# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Repository Guidelines

## Project Structure & Module Organization

```
peso/
├── cmd/
│   └── main.go                 # Entry point dell'applicazione
├── internal/
│   ├── domain/                 # Business logic e entità di dominio
│   │   ├── user/              # Entità User e logica correlata
│   │   ├── weight/            # Entità Weight e logica correlata
│   │   └── goal/              # Entità Goal e logica correlata
│   ├── infrastructure/        # Implementazioni concrete (database, web)
│   │   ├── persistence/       # Repository SQLite
│   │   └── web/              # Handler HTMX e routing
│   ├── application/           # Use cases e servizi applicativi
│   └── interfaces/            # Port interfaces per hexagonal architecture
├── templates/                  # Template HTML/HTMX (root level)
├── web/static/                 # CSS, JS, immagini (embedded in assets.go)
└── migrations/              # Script SQL per database
```

## Build, Test, and Development Commands

```bash
# Sviluppo locale
go run cmd/main.go

# Build di produzione
go build -o bin/peso cmd/main.go

# Test con coverage
go test ./... -cover

# Test di integrazione
go test ./... -tags=integration

# Linting
golangci-lint run

# Format del codice
go fmt ./...

# Tidy delle dipendenze
go mod tidy

# Makefile commands (alternative)
make run          # Esegue in modalità sviluppo
make build        # Compila il binario
make test         # Esegue tutti i test
make test-coverage # Test con coverage
make test-domain   # Test solo del domain layer
make lint          # Linting (golangci-lint o go fmt)
make fmt           # Format del codice
make clean         # Pulisce i file compilati
make reset-db      # Resetta il database locale
```

## Coding Style & Naming Conventions

- Seguire le convenzioni standard di Go (gofmt, golint)
- Nomi di package in minuscolo, senza underscore
- Interfacce terminano con -er quando possibile (es. `Weigher`, `UserRepository`)
- Costanti in PascalCase
- Variabili locali in camelCase
- Funzioni pubbliche in PascalCase, private in camelCase
- Commenti su funzioni pubbliche obbligatori
- Errori personalizzati con il suffisso Error (es. `ValidationError`)

## Testing Guidelines

- Test unitari per la business logic nel dominio
- Test di integrazione per i repository
- Test end-to-end per i handler web
- Mocking delle dipendenze esterne (database, servizi esterni)
- Coverage minimo del 80% per il domain layer
- File di test nella stessa directory del codice sorgente
- Test fixtures in `testdata/` directory
- Utilizzare table-driven tests quando appropriato

## Commit & Pull Request Guidelines

- Seguire Conventional Commits (vedi README.md)
- Un commit per logica/feature
- Squash di commit prima del merge quando necessario
- PR description deve includere:
  - Descrizione del problema risolto
  - Approccio utilizzato
  - Test effettuati
  - Screenshot per modifiche UI
- Review obbligatoria prima del merge

## Security & Configuration Tips

- Non mai committare credenziali nel codice
- Utilizzare variabili d'ambiente per configurazioni sensibili
- Validazione input rigorosa sui dati in ingresso
- Sanitizzazione dei dati prima del salvataggio in database
- HTTPS in produzione (configurato nel reverse proxy)
- Rate limiting su endpoints critici
- Backup regolari del database SQLite

## Architecture Overview

Il progetto segue una architettura esagonale (Ports & Adapters) con Domain-Driven Design:

- **Domain Layer** (`internal/domain/`): Entità di business (User, Weight, Goal) con value objects e business logic
- **Application Layer** (`internal/application/`): Use cases e servizi applicativi (WeightTracker, GoalTracker) 
- **Infrastructure Layer** (`internal/infrastructure/`): Implementazioni concrete per database (SQLite) e web (HTMX handlers)
- **Interfaces** (`internal/interfaces/`): Port interfaces per hexagonal architecture

Stack tecnologico:
- Backend: Go 1.23+ con Gorilla Mux per routing
- Frontend: HTMX per interazioni dinamiche, Chart.js per grafici
- Database: SQLite con migrazioni in `/migrations`
- Templates: HTML templates in `/templates` (embedded via assets.go)

