# Peso - Weight Tracking App

A simple and effective web application built with Go and HTMX for monitoring and visualizing weight progress over time, with support for personalized goals.

## Features

- **Weight Recording**: Quick daily measurement input
- **Interactive Charts**: Temporal trend visualization with Chart.js
- **Personal Goals**: Weight goal setting with automatic progress calculation
- **Trend Analysis**: Automatic calculation of variations and statistics
- **Multi-User**: Separate tracking for multiple users
- **Homelab Ready**: Optimized for home deployment with Docker

**Tech Stack**: Go, HTMX, Chart.js, SQLite

## Prerequisites

- Go 1.21+
- SQLite3
- Make (optional)

## Quick Start

```bash
# Clone the repository
git clone <repo-url>
cd peso

# Install dependencies
go mod tidy

# Run in development mode
go run cmd/main.go

# Or using Make
make run
```

The application will be available at http://localhost:8082


## Configuration

Environment variables (see `.env.example` for defaults):

- `PORT`: Server port (default: 8082)
- `DB_PATH`: SQLite database path (default: ./peso.db)
- `LOG_LEVEL`: Log level (default: info)

## Development

### Available Make Commands

- `make run`: Run application in development mode
- `make build`: Compile binary
- `make test`: Run all tests
- `make clean`: Clean compiled files
- `make docker-build`: Build Docker image

## Docker Deployment

### Building and Running

```bash
# Build the image
docker build -t peso .

# Run the container
docker run -p 8082:8082 -v $(pwd)/data:/app/data peso
```

### Docker Compose

```yaml
version: '3.8'
services:
  peso:
    build: .
    ports:
      - "8082:8082"
    volumes:
      - ./data:/app/data
    environment:
      - DB_PATH=/app/data/peso.db
      - PORT=8082
    restart: unless-stopped
```

### Health Checks

- Health check: `GET /health`
- Readiness check: `GET /ready`

## Deployment Guide

For homelab deployment:

1. Use docker-compose for container orchestration
2. Configure reverse proxy (nginx/caddy) if needed
3. Setup periodic SQLite database backups via volumes
4. Configure SSL certificates for HTTPS access

## Contributing

### Commit Message Convention

This project follows [Conventional Commits](https://conventionalcommits.org/) for clear and automated commit messages.

**Format**:
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `build`, `ci`

**Examples**:
```bash
feat(ui): add weight goal progress indicator
fix(api): resolve database connection timeout
docs(readme): update deployment instructions
```

**Guidelines**:
- Use imperative present tense in subject
- Keep subject under 72 characters
- One commit should do one thing well
- Include scope when relevant (ui, api, db, etc.)

### Development Workflow

1. Create feature branch from `main`
2. Make changes following conventional commits
3. Run tests: `make test`
4. Submit pull request

