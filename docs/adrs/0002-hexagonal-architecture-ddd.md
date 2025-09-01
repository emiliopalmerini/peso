# ADR-0002: Hexagonal Architecture e Domain Driven Design

## Status

Accepted

## Context

Nonostante la semplicità dell'applicazione Peso, vogliamo mantenere un design pulito e facilmente testabile. Dobbiamo decidere come strutturare il codice per garantire:

- Separazione delle responsabilità
- Testabilità del business logic
- Facilità di manutenzione e evoluzione
- Indipendenza dalle tecnologie infrastrutturali (database, web framework)

Anche se il dominio è semplice (peso tracking), vogliamo applicare buone pratiche architetturali per facilitare future estensioni.

## Decision

Adotteremo una **Hexagonal Architecture** (Ports & Adapters) combinata con principi **Domain Driven Design (DDD)**, adattati alla semplicità del nostro contesto.

### Struttura dei Package

```
internal/
├── domain/           # Core business logic
│   ├── user/        # Aggregate User
│   └── weight/      # Aggregate Weight
├── application/     # Use cases e servizi applicativi
├── interfaces/      # Ports (interfacce)
└── infrastructure/  # Adapters (implementazioni concrete)
    ├── persistence/ # Database adapters
    └── web/        # HTTP/HTMX handlers
```

### Principi Applicati

**Hexagonal Architecture:**
- **Domain** al centro, indipendente da infrastruttura
- **Ports** (interfacce) definiscono contratti
- **Adapters** implementano i ports per tecnologie specifiche
- **Application Layer** orchestra i use cases

**Domain Driven Design:**
- **Aggregates**: User e Weight come radici di aggregati
- **Value Objects**: per concetti immutabili (es. peso, timestamp)
- **Domain Services**: per logica che non appartiene a un singolo aggregate
- **Repository Pattern**: per astrazione persistenza

### Regole di Dipendenza

1. **Domain** non dipende da niente (zero imports esterni)
2. **Application** dipende solo dal Domain
3. **Infrastructure** dipende da Application e Domain
4. **Interfaces** definiscono contratti per Infrastructure

## Consequences

### Positive

**Testabilità**
- Business logic completamente isolata e testabile
- Mock facili per i ports/interfaces
- Test unitari veloci senza dipendenze esterne

**Maintainability**
- Codice organizzato per responsabilità
- Modifiche tecnologiche limitate agli adapters
- Evoluzione del dominio indipendente dall'infrastruttura

**Flessibilità**
- Possibilità di cambiare database senza impatto sul domain
- Possibilità di aggiungere interfacce (API, CLI) facilmente
- Extension point chiari per nuove funzionalità

**Code Quality**
- Separazione netta delle responsabilità
- Dependency inversion principle applicato
- Codice più leggibile e comprensibile

### Negative

**Complessità Iniziale**
- Più file e directory da gestire
- Curva di apprendimento per pattern DDD
- Overhead per applicazione semplice

**Over-engineering Risk**
- Potenziale over-engineering per dominio molto semplice
- Più codice boilerplate iniziale

**Development Overhead**
- Più tempo per setup iniziale
- Necessità di definire interfacce anche per operazioni semplici

### Mitigazioni

**Semplicità Adattiva**
- Implementazione graduale dei pattern
- Start semplice, evoluzione iterativa
- Focus su value objects e aggregates essenziali

**Pragmatismo**
- Non forzare pattern dove non aggiungono valore
- Permettere eccezioni giustificate per semplicità
- Documentare decisioni di design

## Implementation Guidelines

1. **Start Small**: Iniziare con aggregates essenziali (User, Weight)
2. **Evolutionary Design**: Aggiungere complessità solo quando necessaria  
3. **Clear Boundaries**: Mantenere boundaries ben definiti tra layers
4. **Interface Segregation**: Ports piccoli e focalizzati
5. **Domain First**: Iniziare sempre dal modeling del dominio