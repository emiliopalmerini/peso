# ADR-0001: Weight Tracking Application

## Status

Accepted

## Context

Nel nostro homelab abbiamo la necessità di tracciare il peso di due utenti specifici (Giada ed Emilio) in modo semplice e efficace. L'obiettivo è creare un'applicazione web che permetta di:

- Registrare facilmente le misurazioni del peso
- Visualizzare l'andamento del peso nel tempo tramite grafici
- Mantenere uno storico completo delle misurazioni
- Essere facilmente deployabile nell'infrastruttura homelab esistente

## Decision

Implementeremo un'applicazione web chiamata "Peso" con le seguenti caratteristiche:

### Stack Tecnologico
- **Backend**: Go (Golang) per le performance, semplicità e facilità di deployment
- **Frontend**: HTMX per interattività senza JavaScript complesso
- **Database**: SQLite per semplicità e zero-configuration
- **Visualizzazione**: Chart.js per i grafici interattivi
- **Deployment**: Docker + Docker Compose per facilità di gestione nel homelab

### Funzionalità Core
1. **Inserimento Peso**: Form semplice per registrare peso per utente specifico
2. **Visualizzazione Storico**: Lista delle misurazioni recenti
3. **Grafici**: Andamento temporale del peso per entrambi gli utenti
4. **API Salute**: Endpoints per health check e monitoring

### Architettura
- Applicazione monolitica per semplicità
- Separazione chiara tra presentation, business logic e persistence
- RESTful API per future estensioni (mobile app, integazioni)

## Consequences

### Positive
- **Semplicità**: Stack tecnologico minimale e ben conosciuto
- **Performance**: Go offre ottime performance con footprint ridotto
- **Maintainability**: Codebase piccola e focalizzata su un dominio specifico
- **Deployment**: Docker semplifica deployment e gestione nel homelab
- **Zero Dependencies**: SQLite elimina la necessità di database server separato

### Negative
- **Scalabilità Limitata**: SQLite non adatto per high-concurrency (accettabile per 2 utenti)
- **Single Point of Failure**: Applicazione monolitica (mitigato da backup regolari)
- **Vendor Lock-in**: Alcune scelte tecnologiche specifiche (mitigato dalla semplicità)

### Neutral
- **Learning Curve**: HTMX può richiedere apprendimento iniziale
- **Future Evolution**: L'architettura monolitica può richiedere refactoring per crescita futura