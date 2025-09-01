# ADR-0004: Sistema di Grafici e Visualizzazione

## Status

Accepted

## Context

L'applicazione Peso deve fornire visualizzazioni grafiche per permettere agli utenti di monitorare l'andamento del peso nel tempo. Dobbiamo decidere:

- Quale libreria utilizzare per i grafici
- Che tipi di grafici implementare
- Come gestire i dati per la visualizzazione
- Come integrare i grafici con HTMX
- Come rendere i grafici responsive e accessibili

I grafici devono mostrare l'andamento del peso di entrambi gli utenti (Giada ed Emilio) in modo chiaro e intuitivo.

## Decision

Implementeremo un sistema di grafici utilizzando **Chart.js** con le seguenti caratteristiche:

### Scelta della Libreria: Chart.js

**Motivazioni:**
- Libreria matura e stabile
- Ottima documentazione e community
- Grafici responsive out-of-the-box  
- Integrazione semplice con HTML/HTMX
- Configurabile e personalizzabile
- Performance adeguate per il nostro use case
- Zero dipendenze da framework frontend complessi

### Tipi di Grafici

**1. Line Chart - Andamento Temporale**
```javascript
// Grafico principale con andamento peso nel tempo
{
  type: 'line',
  data: {
    datasets: [
      {
        label: 'Giada',
        data: [...], // {x: date, y: weight}
        borderColor: '#e91e63',
        backgroundColor: 'rgba(233, 30, 99, 0.1)'
      },
      {
        label: 'Emilio', 
        data: [...],
        borderColor: '#2196f3',
        backgroundColor: 'rgba(33, 150, 243, 0.1)'
      }
    ]
  }
}
```

**2. Bar Chart - Confronto Mensile**
```javascript
// Confronto peso medio mensile
{
  type: 'bar',
  data: {
    labels: ['Gen', 'Feb', 'Mar', ...],
    datasets: [
      {
        label: 'Giada - Media Mensile',
        data: [...],
        backgroundColor: '#e91e63'
      },
      {
        label: 'Emilio - Media Mensile', 
        data: [...],
        backgroundColor: '#2196f3'
      }
    ]
  }
}
```

**3. Progress Indicators**
- Indicatori di trend (in crescita/calo/stabile)
- Variazione percentuale ultimo mese
- Obiettivi peso (se configurati)

### Integrazione con HTMX

**Aggiornamento Dinamico:**
```html
<!-- Container grafico che viene aggiornato via HTMX -->
<div id="weight-chart-container" 
     hx-get="/api/charts/weight-trend?period=6m"
     hx-trigger="load, weight-updated from:body"
     hx-swap="innerHTML">
  <canvas id="weightChart"></canvas>
</div>

<!-- Form che triggera aggiornamento -->
<form hx-post="/weights" 
      hx-trigger="submit"
      hx-after-request="htmx.trigger('body', 'weight-updated')">
  ...
</form>
```

**Endpoint API:**
```go
// GET /api/charts/weight-trend?period=6m&users=giada,emilio
func (h *ChartHandler) WeightTrend(w http.ResponseWriter, r *http.Request) {
    // 1. Parse parametri (period, users)
    // 2. Fetch dati dal repository
    // 3. Transform data per Chart.js
    // 4. Render template con script Chart.js
}
```

### Configurazione Responsiva

```javascript
const chartConfig = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
      labels: {
        usePointStyle: true
      }
    },
    tooltip: {
      mode: 'index',
      intersect: false,
      callbacks: {
        title: function(tooltipItems) {
          return new Date(tooltipItems[0].parsed.x).toLocaleDateString('it-IT');
        },
        label: function(context) {
          return `${context.dataset.label}: ${context.parsed.y} kg`;
        }
      }
    }
  },
  scales: {
    x: {
      type: 'time',
      time: {
        unit: 'day',
        displayFormats: {
          day: 'DD/MM'
        }
      }
    },
    y: {
      beginAtZero: false,
      title: {
        display: true,
        text: 'Peso (kg)'
      }
    }
  }
};
```

### Periodi Temporali Supportati

- **1 settimana**: Dettaglio giornaliero
- **1 mese**: Dettaglio giornaliero con media mobile 3 giorni  
- **3 mesi**: Dettaglio settimanale
- **6 mesi**: Dettaglio settimanale
- **1 anno**: Dettaglio mensile
- **Tutto**: Dettaglio mensile con media

### Features Aggiuntive

**1. Filtri Interattivi**
- Selezione periodo temporale
- Toggle visibilità utenti
- Zoom su sezioni specifiche

**2. Export Capabilities**
- Download grafico come PNG
- Export dati CSV
- Condivisione link con filtri

**3. Accessibility**
- Tabella dati alternativa per screen readers
- Keyboard navigation per controlli
- Color scheme accessibile (contrasto WCAG AA)

## Consequences

### Positive

**User Experience**
- Visualizzazione immediata e intuitiva dei trend
- Confronto facile tra utenti
- Grafici responsive su tutti i dispositivi
- Aggiornamenti real-time con HTMX

**Tecnico**
- Integrazione semplice con stack esistente
- Performance ottimali per dataset piccoli/medi
- Configurazione flessibile e estendibile
- Bundle size ragionevole (~60kb minified)

**Maintenance**
- Libreria stabile e ben mantenuta
- Documentazione eccellente
- Community attiva per supporto
- API consistenti tra versioni

### Negative

**Limitazioni**
- Dipendenza da JavaScript (mitigata da fallback tabellare)
- Chart.js non ottimale per dataset molto grandi
- Personalizzazione avanzata può essere complessa

**Performance**
- Bundle aggiuntivo da caricare
- Rendering client-side (vs server-side)
- Memory usage per dataset storici grandi

### Mitigazioni

**Progressive Enhancement**
- Tabelle HTML come fallback per JS disabilitato
- Lazy loading dei grafici sotto la fold
- Paginazione/filtering per dataset grandi

**Performance Optimization**
- CDN per Chart.js
- Caching dei dati aggregati
- Debouncing degli aggiornamenti real-time

**Accessibility Compliance**
- Alt text per grafici
- Dati tabellari sempre disponibili
- Focus management per controlli interattivi

## Implementation Plan

1. **MVP Charts**
   - Line chart base con dati ultimi 3 mesi
   - Integrazione HTMX per aggiornamenti
   - Responsive design base

2. **Enhanced Features**
   - Selezione periodi temporali
   - Bar chart comparativo mensile
   - Export capabilities

3. **Advanced Features**
   - Zoom e pan interattivi
   - Annotazioni su grafici
   - Configurazione obiettivi peso