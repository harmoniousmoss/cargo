# Cargo Backend

## Overview
Go backend for Cargo, an HTTP load testing tool. Fires concurrent requests against target domains and records results.

## Tech Stack
- Go 1.22 with Fiber v2 (HTTP framework)
- SQLite via go-sqlite3 (requires CGO_ENABLED=1)
- Air for hot reload during development

## Project Structure
- `main.go` — Fiber app setup, API routes
- `handlers/main_handler.go` — RPS test logic, concurrent HTTP requests
- `libs/database.go` — SQLite schema, queries, CRUD operations, auth
- `libs/env_utils.go` — .env file loading
- `urls.txt` — URL paths to test against
- `.air.toml` — Air hot reload config

## Running
```bash
# Requires .env file with DOMAIN=https://example.com/
CGO_ENABLED=1 air
```

## API Endpoints
- `GET /` — Health check
- `POST /api/login` — Authenticate with username/password, returns session token
- `GET /api/validate` — Validate session token (via Authorization header)
- `POST /api/logout` — Delete session token
- `POST /api/start-test` — Start a test (`{ domain, iterations }`)
- `GET /api/results?since=N` — Poll current test results (incremental)
- `GET /api/runs` — List all past test runs
- `GET /api/runs/:id/results` — Get results for a specific past run

## Auth
- Users table with SHA-256 hashed passwords
- Sessions table with random tokens
- Default user seeded on startup (see `seedUser()` in database.go)

## Notes
- go-sqlite3 requires CGo — always build with `CGO_ENABLED=1`
- Database file is `cargo.db` (gitignored)
- The frontend connects via server-side proxy (SvelteKit API routes), not direct browser requests
