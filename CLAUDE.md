# Peso - Weight Tracking App

## Project Overview

Weight tracking web application built with Go, HTMX, and Chart.js. Uses SQLite for data persistence with multi-user support.

## Tech Stack

- **Language**: Go 1.24.0
- **Module**: peso
- **Frontend**: HTMX, Chart.js (no templ templates)
- **Database**: SQLite (modernc.org/sqlite)
- **Architecture**: Multi-user support

## Development Commands

### Build

```bash
make build
# or
go build -o bin/peso ./cmd/main.go
```

### Run

```bash
make run
# or
./bin/peso
```

Server runs on port 8082 by default.

### Test

```bash
make test
# or
go test ./...
```

### Docker

```bash
make docker-build
docker-compose up
```

## Environment Variables

| Variable    | Default     | Description              |
|-------------|-------------|--------------------------|
| `PORT`      | 8082        | Server port              |
| `DB_PATH`   | ./peso.db   | SQLite database path     |
| `LOG_LEVEL` | info        | Log level                |

## Health Endpoints

- `/health` - Health check
- `/ready` - Readiness check

## Project Structure

```
peso/
├── cmd/main.go       # Application entry point
├── internal/         # Internal packages
├── migrations/       # Database migrations
├── static/           # Static assets
├── templates/        # HTML templates
├── bin/              # Compiled binaries
└── Makefile          # Build automation
```

## Code Style

- Follow standard Go conventions (go fmt, go vet)
- Use conventional commits for git messages
- Run `make lint` before committing
