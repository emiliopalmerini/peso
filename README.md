# Peso - Weight Tracking App

Applicazione web per il tracking del peso di Giada ed Emilio nel homelab.
Un'app semplice ed efficace costruita con Go e HTMX per monitorare e visualizzare l'andamento del peso nel tempo, con supporto per obiettivi personalizzati.

## Funzionalità

- ⚖️ **Registrazione Peso**: Inserimento rapido delle misurazioni giornaliere
- 📊 **Grafici Interattivi**: Visualizzazione dell'andamento temporale con Chart.js
- 🎯 **Obiettivi Personali**: Impostazione di goal di peso con calcolo automatico dei progressi
- 📈 **Analisi Trend**: Calcolo automatico di variazioni e statistiche
- 👥 **Multi-Utente**: Tracking separato per Giada ed Emilio
- 🏠 **Homelab Ready**: Ottimizzato per deployment domestico con Docker

Stack: Go, HTMX, Chart.js, SQLite

## Requisiti

- Go 1.21+
- SQLite3
- Make (opzionale)

## Esecuzione locale

```bash
# Clone del repository
git clone <repo-url>
cd peso

# Installazione dipendenze
go mod tidy

# Esecuzione in modalità sviluppo
go run cmd/main.go

# Oppure con Make
make run
```

L'applicazione sarà disponibile su http://localhost:8080


## Variabili d'ambiente supportate

Vedi `.env.example` per i default. Principali:

- `PORT`: Porta del server (default: 8080)
- `DB_PATH`: Path del database SQLite (default: ./peso.db)
- `LOG_LEVEL`: Livello di log (default: info)

## Comandi Makefile utili

- `make run`: Esegue l'applicazione in modalità sviluppo
- `make build`: Compila il binario
- `make test`: Esegue tutti i test
- `make clean`: Pulisce i file compilati
- `make docker-build`: Costruisce l'immagine Docker

## Docker

```bash
# Build dell'immagine
docker build -t peso .

# Run del container
docker run -p 8080:8080 -v $(pwd)/data:/app/data peso
```

## Health & Readiness

- Health check: `GET /health`
- Readiness check: `GET /ready`

## Deploy

Per il deployment nel homelab:

1. Utilizzare docker-compose per il deployment
2. Configurare reverse proxy (nginx/caddy) se necessario
3. Backup periodico del database SQLite tramite volume

Esempio docker-compose.yml:

```yaml
version: '3.8'
services:
  peso:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - DB_PATH=/app/data/peso.db
    restart: unless-stopped
```

## Commit message template (Conventional Commits)

Usiamo lo standard Conventional Commits per messaggi chiari e automatizzabili.

Formato base:

```
<type>(<scope>)<!>: <subject>

<body>

<footer>
```

- type: tipo di cambiamento. Esempi comuni: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `build`, `ci`.
- scope: (opzionale) area toccata, es. `templates`, `encounters`, `router`.
- !: (opzionale) indica breaking change.
- subject: (obbligatorio) riassunto al presente, in minuscolo, senza punto finale.
- body: (opzionale) contesto/motivazione, dettagli tecnici se utili.
- footer: (opzionale) riferimenti a issue/PR o `BREAKING CHANGE:` con spiegazione.

Esempi:

```
fix(templates): allinea route HTMX a /encounters/*

docs(readme): aggiungi template per Conventional Commits
```

Note:
- Imperativo presente nel subject (es. "aggiungi", "correggi").
- Mantieni il subject entro ~72 caratteri quando possibile.
- Un commit dovrebbe fare una cosa sola e bene.

## Pre-commit hook

## ADR (Architectural Decision Records)

La documentazione delle decisioni architetturali è disponibile in `docs/adr`.
Indice ADR: [docs/adrs/README.md](./docs/adrs/README.md)

