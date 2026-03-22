package handlers

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"gorps/libs"
)

// RunState tracks the currently active test run
type RunState struct {
	mu      sync.Mutex
	RunID   int64
	Running bool
	Total   int
}

var currentRun = &RunState{}

// GetCurrentRunState returns the current run ID and status
func GetCurrentRunState() (int64, bool, int) {
	currentRun.mu.Lock()
	defer currentRun.mu.Unlock()
	return currentRun.RunID, currentRun.Running, currentRun.Total
}

// TestRPSWithIterations performs the RPS test with specified iterations
func TestRPSWithIterations(domain string, pathsFile string, iterations int) {
	paths := loadPaths(pathsFile)
	if len(paths) == 0 {
		log.Fatal("No paths found in the file")
	}

	totalRequests := len(paths) * iterations

	// Create a test run in the database
	runID, err := libs.CreateTestRun(domain, iterations, totalRequests)
	if err != nil {
		log.Printf("Failed to create test run: %s", err)
		return
	}

	// Update in-memory state
	currentRun.mu.Lock()
	currentRun.RunID = runID
	currentRun.Running = true
	currentRun.Total = totalRequests
	currentRun.mu.Unlock()

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	fmt.Printf("Starting RPS test run #%d with %d iterations (%d total requests)\n", runID, iterations, totalRequests)

	var wg sync.WaitGroup
	for i := 0; i < totalRequests; i++ {
		randomPath := paths[rand.Intn(len(paths))]
		fullURL := fmt.Sprintf("%s%s", domain, randomPath)

		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sendRequest(client, url, runID)
		}(fullURL)
	}

	wg.Wait()

	// Mark run as completed
	if err := libs.FinishTestRun(runID); err != nil {
		log.Printf("Failed to finish test run: %s", err)
	}

	currentRun.mu.Lock()
	currentRun.Running = false
	currentRun.mu.Unlock()

	fmt.Printf("✅ RPS test run #%d completed! (%d requests processed)\n", runID, totalRequests)
}

// sendRequest sends a single HTTP GET request and stores the result in the database
func sendRequest(client *http.Client, url string, runID int64) {
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to send request to %s: %s\n", url, err)
		if dbErr := libs.InsertResult(runID, url, 0, err.Error(), timestamp); dbErr != nil {
			log.Printf("Failed to insert result: %s", dbErr)
		}
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response: %s -> Status Code: %d\n", url, resp.StatusCode)
	if dbErr := libs.InsertResult(runID, url, resp.StatusCode, "", timestamp); dbErr != nil {
		log.Printf("Failed to insert result: %s", dbErr)
	}
}

// loadPaths loads paths from a file
func loadPaths(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file %s: %s", filename, err)
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file %s: %s", filename, err)
	}

	return paths
}
