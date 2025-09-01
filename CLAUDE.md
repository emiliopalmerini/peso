# Repository Guidelines

## Project Structure & Module Organization

```
peso/
├── cmd/
│   └── main.go                 # Entry point dell'applicazione
├── internal/
│   ├── domain/                 # Business logic e entità di dominio
│   │   ├── user/              # Entità User e logica correlata
│   │   └── weight/            # Entità Weight e logica correlata
│   ├── infrastructure/        # Implementazioni concrete (database, web)
│   │   ├── persistence/       # Repository SQLite
│   │   └── web/              # Handler HTMX e routing
│   ├── application/           # Use cases e servizi applicativi
│   └── interfaces/            # Port interfaces per hexagonal architecture
├── web/
│   ├── static/               # CSS, JS, immagini
│   └── templates/           # Template HTML/HTMX
├── migrations/              # Script SQL per database
└── docs/
    └── adrs/               # Architectural Decision Records
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

La documentazione delle decisioni architetturali è disponibile in `docs/adr`.
Indice ADR: [docs/adrs/README.md](./docs/adrs/README.md)
