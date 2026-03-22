package libs

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the SQLite database and creates tables
func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	// Enable WAL mode for better concurrent read/write
	DB.SetMaxOpenConns(1)

	createTables()
	seedUser()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT NOT NULL UNIQUE,
			user_id INTEGER NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS test_runs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			domain TEXT NOT NULL,
			iterations INTEGER NOT NULL,
			total_requests INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'running',
			created_at TEXT NOT NULL,
			finished_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS test_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id INTEGER NOT NULL,
			url TEXT NOT NULL,
			status_code INTEGER NOT NULL DEFAULT 0,
			error TEXT,
			timestamp TEXT NOT NULL,
			FOREIGN KEY (run_id) REFERENCES test_runs(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_run_id ON test_results(run_id)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Failed to create table: %s", err)
		}
	}
}

// CreateTestRun inserts a new test run and returns its ID
func CreateTestRun(domain string, iterations, totalRequests int) (int64, error) {
	result, err := DB.Exec(
		`INSERT INTO test_runs (domain, iterations, total_requests, status, created_at) VALUES (?, ?, ?, 'running', ?)`,
		domain, iterations, totalRequests, time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertResult inserts a single test result
func InsertResult(runID int64, url string, statusCode int, errMsg string, timestamp string) error {
	_, err := DB.Exec(
		`INSERT INTO test_results (run_id, url, status_code, error, timestamp) VALUES (?, ?, ?, ?, ?)`,
		runID, url, statusCode, errMsg, timestamp,
	)
	return err
}

// FinishTestRun marks a test run as completed
func FinishTestRun(runID int64) error {
	_, err := DB.Exec(
		`UPDATE test_runs SET status = 'completed', finished_at = ? WHERE id = ?`,
		time.Now().UTC().Format(time.RFC3339), runID,
	)
	return err
}

// TestRunResult represents a result row from the database
type TestRunResult struct {
	URL        string `json:"url"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
	Timestamp  string `json:"timestamp"`
}

// GetResultsSince returns results for a run after a given offset
func GetResultsSinceDB(runID int64, since int) ([]TestRunResult, error) {
	rows, err := DB.Query(
		`SELECT url, status_code, error, timestamp FROM test_results WHERE run_id = ? ORDER BY id ASC LIMIT -1 OFFSET ?`,
		runID, since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TestRunResult
	for rows.Next() {
		var r TestRunResult
		var errStr sql.NullString
		if err := rows.Scan(&r.URL, &r.StatusCode, &errStr, &r.Timestamp); err != nil {
			return nil, err
		}
		if errStr.Valid {
			r.Error = errStr.String
		}
		results = append(results, r)
	}
	return results, nil
}

// TestRunInfo holds info about a test run
type TestRunInfo struct {
	ID            int64  `json:"id"`
	Domain        string `json:"domain"`
	Iterations    int    `json:"iterations"`
	TotalRequests int    `json:"totalRequests"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
	FinishedAt    string `json:"finishedAt,omitempty"`
}

// GetTestRun returns info about a specific test run
func GetTestRun(runID int64) (*TestRunInfo, error) {
	var info TestRunInfo
	var finishedAt sql.NullString
	err := DB.QueryRow(
		`SELECT id, domain, iterations, total_requests, status, created_at, finished_at FROM test_runs WHERE id = ?`,
		runID,
	).Scan(&info.ID, &info.Domain, &info.Iterations, &info.TotalRequests, &info.Status, &info.CreatedAt, &finishedAt)
	if err != nil {
		return nil, err
	}
	if finishedAt.Valid {
		info.FinishedAt = finishedAt.String
	}
	return &info, nil
}

// GetAllTestRuns returns all test runs ordered by most recent
func GetAllTestRuns() ([]TestRunInfo, error) {
	rows, err := DB.Query(
		`SELECT id, domain, iterations, total_requests, status, created_at, finished_at FROM test_runs ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []TestRunInfo
	for rows.Next() {
		var r TestRunInfo
		var finishedAt sql.NullString
		if err := rows.Scan(&r.ID, &r.Domain, &r.Iterations, &r.TotalRequests, &r.Status, &r.CreatedAt, &finishedAt); err != nil {
			return nil, err
		}
		if finishedAt.Valid {
			r.FinishedAt = finishedAt.String
		}
		runs = append(runs, r)
	}
	return runs, nil
}

// GetResultCount returns the number of results for a test run
func GetResultCount(runID int64) (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM test_results WHERE run_id = ?`, runID).Scan(&count)
	return count, err
}

// hashPassword creates a SHA-256 hash of the password
func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

// seedUser creates the default user if it doesn't exist
func seedUser() {
	var count int
	DB.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, "harmoniousmoss").Scan(&count)
	if count == 0 {
		hash := hashPassword("harmoniousmoss")
		_, err := DB.Exec(`INSERT INTO users (username, password_hash) VALUES (?, ?)`, "harmoniousmoss", hash)
		if err != nil {
			log.Printf("Failed to seed user: %s", err)
		}
	}
}

// generateToken creates a random session token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Authenticate checks username/password and returns a session token
func Authenticate(username, password string) (string, error) {
	hash := hashPassword(password)
	var userID int64
	err := DB.QueryRow(`SELECT id FROM users WHERE username = ? AND password_hash = ?`, username, hash).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := generateToken()
	if err != nil {
		return "", err
	}

	_, err = DB.Exec(
		`INSERT INTO sessions (token, user_id, created_at) VALUES (?, ?, ?)`,
		token, userID, time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateSession checks if a session token is valid
func ValidateSession(token string) bool {
	if token == "" {
		return false
	}
	var id int64
	err := DB.QueryRow(`SELECT id FROM sessions WHERE token = ?`, token).Scan(&id)
	return err == nil
}

// DeleteSession removes a session token (logout)
func DeleteSession(token string) error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}
