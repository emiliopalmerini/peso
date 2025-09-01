# ADR-0003: Entità di Dominio (User, Weight e Goal)

## Status

Accepted

## Context

Dobbiamo definire le entità di dominio core per l'applicazione Peso. Il dominio è focalizzato sul tracking del peso di due utenti specifici (Giada ed Emilio) con la possibilità di settare obiettivi di peso. Dobbiamo decidere:

- Come modellare User, Weight e Goal come aggregates
- Quali attributi e comportamenti appartengono a ciascuna entità
- Come gestire le relazioni tra User, Weight e Goal
- Quali validation rules applicare
- Come gestire l'identificazione univoca
- Come calcolare i progressi verso gli obiettivi

Il sistema deve essere semplice ma estensibile per future evoluzioni.

## Decision

Defineremo tre **Aggregates** principali seguendo i principi DDD:

### User Aggregate

```go
type User struct {
    ID       UserID
    Name     string
    Email    string  // opzionale per future estensioni
    Active   bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserID string
```

**Responsabilità:**
- Gestire informazioni anagrafiche utente
- Validare formato nome (non vuoto, lunghezza max)
- Controllare stato attivo/inattivo
- Fornire identificazione univoca per le misurazioni

**Invarianti:**
- Name non può essere vuoto
- UserID deve essere univoco
- CreatedAt non può essere modificato dopo creazione

### Weight Aggregate

```go
type Weight struct {
    ID          WeightID
    UserID      UserID
    Value       WeightValue
    Unit        WeightUnit
    MeasuredAt  time.Time
    Notes       string  // opzionale
    CreatedAt   time.Time
}

type WeightValue float64
type WeightUnit string
type WeightID string

const (
    WeightUnitKg WeightUnit = "kg"
    WeightUnitLb WeightUnit = "lb"
)
```

**Responsabilità:**
- Registrare singola misurazione peso
- Validare valore peso (positivo, range realistico)
- Associare misurazione a utente specifico
- Mantenere timestamp preciso della misurazione
- Supportare note opzionali per contesto

**Invarianti:**
- Value deve essere > 0 e < 1000kg (validazione realistica)
- UserID deve riferirsi a utente esistente e attivo
- MeasuredAt non può essere nel futuro
- Unit deve essere tra quelle supportate

### Goal Aggregate

```go
type Goal struct {
    ID         GoalID
    UserID     UserID
    TargetWeight WeightValue
    Unit       WeightUnit
    TargetDate TargetDate
    Description string  // opzionale
    Active     bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type GoalID string

type TargetDate struct {
    year  int
    month int
    day   int
}

// Methods per TargetDate
func NewTargetDate(year, month, day int) (TargetDate, error)
func (td TargetDate) Year() int
func (td TargetDate) Month() int  
func (td TargetDate) Day() int
func (td TargetDate) IsValid() bool
func (td TargetDate) IsPast() bool
func (td TargetDate) DaysUntil() int
func (td TargetDate) ToTime() time.Time
```

**Responsabilità:**
- Definire obiettivo di peso per un utente
- Validare peso target realistico
- Gestire data target dell'obiettivo
- Calcolare giorni rimanenti per raggiungere goal
- Supportare descrizione opzionale per motivazione

**Invarianti:**
- TargetWeight deve essere > 0 e < 1000kg
- UserID deve riferirsi a utente esistente e attivo
- TargetDate non può essere nel passato al momento della creazione
- Solo un Goal attivo per utente
- Unit deve essere tra quelle supportate

### Value Objects

**UserID, WeightID, GoalID**: Identificatori tipizzati per type safety

**WeightValue**: Wrapper per float64 con validazione integrata

**WeightUnit**: Enum per unità di misura supportate

**TargetDate**: Value object personalizzato per date senza dipendenze esterne, con validazioni e utility methods specifici per il dominio

### Domain Services

```go
type WeightTracker interface {
    RecordWeight(userID UserID, value WeightValue, unit WeightUnit, measuredAt time.Time) (*Weight, error)
    GetWeightHistory(userID UserID, period TimePeriod) ([]*Weight, error)
    CalculateWeightTrend(userID UserID, period TimePeriod) (WeightTrend, error)
}

type GoalTracker interface {
    SetGoal(userID UserID, targetWeight WeightValue, unit WeightUnit, targetDate TargetDate, description string) (*Goal, error)
    GetActiveGoal(userID UserID) (*Goal, error)
    CalculateProgress(userID UserID) (GoalProgress, error)
    DeactivateGoal(goalID GoalID) error
}

type GoalProgress struct {
    Goal           *Goal
    CurrentWeight  WeightValue
    WeightToLose   WeightValue  // può essere negativo se deve guadagnare
    DaysRemaining  int
    WeightPerDay   WeightValue  // peso da perdere/guadagnare per giorno
    ProgressPercent float64
    IsOnTrack      bool
}
```

**Responsabilità:**
- Coordinare operazioni tra aggregates
- Calcolare statistiche e trend
- Applicare business rules complesse
- Gestire obiettivi e calcolare progressi verso goal
- Determinare se l'utente è sulla buona strada per raggiungere l'obiettivo

### Repository Interfaces (Ports)

```go
type UserRepository interface {
    Save(user *User) error
    FindByID(id UserID) (*User, error)
    FindByName(name string) (*User, error)
    FindActive() ([]*User, error)
}

type WeightRepository interface {
    Save(weight *Weight) error
    FindByID(id WeightID) (*Weight, error)
    FindByUserID(userID UserID, limit int) ([]*Weight, error)
    FindByUserIDAndPeriod(userID UserID, from, to time.Time) ([]*Weight, error)
}

type GoalRepository interface {
    Save(goal *Goal) error
    FindByID(id GoalID) (*Goal, error)
    FindActiveByUserID(userID UserID) (*Goal, error)
    FindByUserID(userID UserID) ([]*Goal, error)
    DeactivateByUserID(userID UserID) error
}
```

### Business Rules

1. **User Management**
   - Solo utenti attivi possono registrare pesi
   - Nome utente deve essere unique
   - Non è possibile eliminare utenti con misurazioni

2. **Weight Recording**
   - Massimo 10 misurazioni per giorno per utente
   - Peso deve essere in range 10-500 kg
   - Non è possibile modificare peso storico (immutabilità)
   - Data misurazione non può essere nel futuro

3. **Goal Management**
   - Solo un Goal attivo per utente alla volta
   - TargetDate deve essere almeno 7 giorni nel futuro
   - TargetWeight deve essere diverso dal peso attuale (min 0.1kg differenza)
   - Goal automaticamente disattivato se TargetDate è passata
   - Peso target deve essere realistico (max 2kg/settimana di variazione)

4. **Data Integrity**
   - Soft delete per preservare storico
   - Audit trail per tutte le modificazioni
   - Validazione referential integrity User <-> Weight <-> Goal

## Consequences

### Positive

**Type Safety**
- Identificatori tipizzati prevengono errori di assegnazione
- Value objects incapsulano validazioni
- Compile-time check per business rules

**Domain Clarity**
- Aggregates chiari con responsabilità ben definite
- Business rules esplicite e documentate
- Separazione tra dati e comportamenti

**Extensibility**
- Struttura pronta per nuovi attributi (email, profilo utente)
- Support multi-unità di misura
- Framework per nuove business rules
- Sistema di obiettivi estendibile per nuovi tipi di goal

**Data Quality**
- Validazioni integrate negli aggregates
- Invarianti garantiti a livello di dominio
- Immutabilità per dati storici
- Value object personalizzato TargetDate elimina dipendenze esterne

**Goal Tracking**
- Calcolo automatico del progresso verso obiettivi
- Indicatori di performance e feasibility
- Motivazione attraverso visualizzazione progressi

### Negative

**Initial Complexity**
- Più codice per dominio semplice
- Value objects aggiungono overhead
- Repository interfaces da implementare
- Logic aggiuntiva per gestione goal e calcoli

**Over-engineering Risk**
- Pattern DDD per dominio molto semplice
- Possibile complessità eccessiva per 2 utenti
- Value object personalizzato TargetDate vs time.Time standard

**Goal Management Complexity**
- Business rules complesse per validazione goal realistici
- Calcoli di progresso che richiedono logica sofisticata
- Gestione stati del goal (attivo/scaduto/raggiunto)

### Mitigazioni

**Progressive Enhancement**
- Iniziare con versioni semplificate degli aggregates
- Aggiungere complessità solo quando necessaria
- Mantenere pragmatismo nelle implementazioni
- Implementare Goal feature come MVP prima della full logic

**Documentation**
- Documentare business rules chiaramente
- Esempi di utilizzo per ogni aggregate
- Test case che dimostrano invarianti
- Documentare algoritmi di calcolo progresso goal

**Pragmatic Choices**
- TargetDate personalizzato ma con conversion methods verso time.Time
- Goal calculations semplificati inizialmente
- Possibilità di switch a time.Time se TargetDate non aggiunge valore