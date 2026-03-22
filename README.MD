![Cargo](https://icbotbz40ngrmv6r.public.blob.vercel-storage.com/cargo.png)

# Cargo - Server Performance Testing Tool
Cargo tests how well your server handles heavy traffic by firing concurrent HTTP requests across your endpoints. Spot failures before your users do.

## Features
- **Concurrent Load Testing**: Fires hundreds of simultaneous HTTP requests using Go's goroutines
- **Configurable Iterations**: Multiply load tests with iteration parameters (99 URLs × N iterations)
- **Real-Time Results**: Live polling of test results with per-request status codes
- **SQLite Persistence**: All test runs and results stored in SQLite for historical review
- **Authentication**: Session-based login to protect access to the tool
- **RESTful API**: Full API for starting tests, polling results, and reviewing past runs

## How It Works

1. Configure a target domain and number of iterations
2. Cargo reads 99 URL paths from `urls.txt`
3. All requests fire concurrently via goroutines (99 × iterations)
4. Each response (status code, URL, timestamp) is stored in SQLite
5. Frontend polls for results every 500ms during active tests
6. Review results: 200 OK, 403 Forbidden, 404 Not Found counts

## Tech Stack
- **Go 1.22** with **Fiber v2** (HTTP framework)
- **SQLite** via go-sqlite3 (requires `CGO_ENABLED=1`)
- **Air** for hot reload during development

## Getting Started

### Prerequisites
- Go 1.22+
- Air (optional, for hot reload)
- GCC (required for go-sqlite3)

### Setup
```bash
# Clone the repo
git clone git@github.com:harmoniousmoss/cargo.git
cd cargo

# Create .env file
echo "DOMAIN=https://example.com/" > .env

# Run with Air (hot reload)
CGO_ENABLED=1 air

# Or run directly
CGO_ENABLED=1 go run .
```

The server starts on port 8080.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Health check |
| `POST` | `/api/login` | Authenticate (`{ username, password }`) |
| `GET` | `/api/validate` | Validate session token |
| `POST` | `/api/logout` | End session |
| `POST` | `/api/start-test` | Start load test (`{ domain, iterations }`) |
| `GET` | `/api/results?since=N` | Poll current test results (incremental) |
| `GET` | `/api/runs` | List all past test runs |
| `GET` | `/api/runs/:id/results` | Get results for a specific past run |

### Start Test Example
```json
POST /api/start-test
{
  "domain": "https://example.com",
  "iterations": 5
}
```

## Frontend
The Svelte frontend lives in a separate private repo. It connects to this backend via server-side proxy routes (no CORS needed).

## Project Structure
```
├── main.go                    # Fiber app, API routes
├── handlers/
│   └── main_handler.go        # RPS test logic, concurrent requests
├── libs/
│   ├── database.go            # SQLite schema, queries, auth
│   └── env_utils.go           # .env file loading
├── urls.txt                   # 99 URL paths to test against
├── .air.toml                  # Air hot reload config
└── .env                       # Domain configuration
```

# Terms and Conditions
Cargo is designed strictly for testing and evaluating the performance of servers that you own or have explicit permission to test. By using this tool, you agree to the following terms:

## Authorized Use Only
You must only use Cargo on servers you own or have written consent to test. Unauthorized use on third-party servers is strictly prohibited.

## No Malicious Intent
Cargo must not be used for malicious activities such as Distributed Denial-of-Service (DDoS) attacks or any other actions intended to disrupt, harm, or degrade systems.

## Compliance with Laws
Users must comply with all applicable local, national, and international laws and regulations when using Cargo.

## Liability Disclaimer
The developers of Cargo are not responsible for any misuse, illegal activity, or damage caused by the tool. Use it at your own risk.

By using this tool, you acknowledge and agree to these terms. Violations may lead to legal consequences.
