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
)

// TestRPS performs the RPS test by sending randomized requests
func TestRPS(domain string, pathsFile string) {
	paths := loadPaths(pathsFile)
	if len(paths) == 0 {
		log.Fatal("No paths found in the file")
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	var wg sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		randomPath := paths[rand.Intn(len(paths))]
		fullURL := fmt.Sprintf("%s%s", domain, randomPath)

		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sendRequest(client, url)
		}(fullURL)
	}

	wg.Wait()
	fmt.Println("✅ RPS test completed successfully!")
}

// sendRequest sends a single HTTP GET request and logs the response
func sendRequest(client *http.Client, url string) {
	fmt.Printf("Sending request to: %s\n", url)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to send request to %s: %s\n", url, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Response: %s -> Status Code: %d\n", url, resp.StatusCode)
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
