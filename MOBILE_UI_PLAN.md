# Modern Mobile-First UI Improvements for Peso

## Overview

Transform the Peso weight tracking app into a native-feeling mobile experience while maintaining the current clean aesthetic and HTMX architecture.

## Design Principles

1. **Thumb-zone optimized** - Primary actions in bottom 1/3 of screen
2. **One-handed usability** - All critical actions reachable with thumb
3. **Progressive disclosure** - Show summary first, details on demand
4. **Native-feeling** - Match iOS/Android interaction patterns
5. **Friction-free daily use** - Minimize taps to log weight

---

## Implementation Plan

### Phase 1: Floating Action Button (FAB) for Weight Entry

**Goal:** One-tap access to weight logging from anywhere

**Files to modify:**
- `web/static/min.css` - Add FAB styles
- `templates/user_dashboard.html` - Add FAB component

**Implementation:**
```
- Fixed position FAB (bottom-right, 24px from edges)
- 56px diameter, primary color
- Plus icon or scale icon
- Triggers bottom sheet for weight entry
- Hide on scroll down, show on scroll up
```

**CSS additions:**
- `.fab` - Base FAB styles with shadow, positioning
- `.fab--hidden` - Transform for scroll hide
- Pulse animation on successful save

---

### Phase 2: Bottom Sheet for Weight Entry

**Goal:** Thumb-friendly form that maintains context

**Files to modify:**
- `web/static/min.css` - Bottom sheet styles
- `templates/user_dashboard.html` - Bottom sheet container
- `templates/weight_form.html` - Adapt form for bottom sheet

**Implementation:**
```
- Slides up from bottom (transform: translateY)
- Backdrop overlay with tap-to-dismiss
- Drag handle at top for swipe-to-dismiss
- Large numeric display for weight value
- Pre-fill with last recorded weight
- Full-width save button at bottom
```

**CSS additions:**
- `.bottom-sheet` - Container with transform animation
- `.bottom-sheet__backdrop` - Semi-transparent overlay
- `.bottom-sheet__handle` - Drag indicator
- `.bottom-sheet--open` - Active state

---

### Phase 3: Dashboard Redesign - Summary Cards

**Goal:** Glanceable information with progressive disclosure

**Files to modify:**
- `web/static/min.css` - Summary card styles
- `templates/user_dashboard.html` - New layout structure
- `templates/partials_goal_summary.html` - Compact card format

**New dashboard layout:**
```
┌─────────────────────────────┐
│ Topbar                      │
├─────────────────────────────┤
│ ┌─────────────────────────┐ │
│ │ Current: 75.2 kg   ↓0.3 │ │  <- Hero stat card
│ │ Last: ieri 14:30        │ │
│ └─────────────────────────┘ │
├─────────────────────────────┤
│ ┌───────────┐ ┌───────────┐ │
│ │ Obiettivo │ │ Media 7g  │ │  <- Stat pills (2-col grid)
│ │ 70.0 kg   │ │ 75.5 kg   │ │
│ └───────────┘ └───────────┘ │
├─────────────────────────────┤
│ [Sparkline - 30 day trend]  │  <- Mini chart, tap for full
├─────────────────────────────┤
│ Storico recente        Vedi>│
│ ├─ Oggi      75.2 kg        │
│ ├─ Ieri      75.5 kg        │
│ └─ 2 gen     75.8 kg        │
└─────────────────────────────┘
│        [FAB +]              │
```

**CSS additions:**
- `.stat-hero` - Large current weight display
- `.stat-pills` - 2-column grid of mini stats
- `.sparkline` - Mini chart container
- `.trend-up`, `.trend-down`, `.trend-neutral` - Trend indicators

---

### Phase 4: Chart Improvements

**Goal:** Touch-friendly period selection, better mobile sizing

**Files to modify:**
- `web/static/min.css` - Period chips, chart sizing
- `templates/user_dashboard.html` - Replace dropdown with chips

**Implementation:**
```
- Replace dropdown with horizontal scrolling chips
- Chip buttons: 1S, 1M, 3M, 6M, 1A, Tutto
- Active chip has filled background
- Chart height: 240px mobile, 320px tablet
- Tap-to-show tooltips (not hover)
```

**CSS additions:**
- `.period-chips` - Horizontal scroll container
- `.period-chip` - Individual chip button
- `.period-chip--active` - Selected state

---

### Phase 5: Swipe Actions on History List

**Goal:** Quick edit/delete without buttons

**Files to modify:**
- `web/static/min.css` - Swipe row styles
- `templates/partials_recent_weights.html` - Swipeable row structure
- `templates/user_dashboard.html` - Add swipe JS handler

**Implementation:**
```
- Swipe left reveals delete action (red)
- Swipe right reveals edit action (blue)
- Touch feedback on swipe
- Auto-close on action or tap elsewhere
- Undo snackbar after delete
```

**CSS additions:**
- `.swipe-row` - Overflow hidden container
- `.swipe-row__content` - Main row content
- `.swipe-row__action--delete` - Red delete background
- `.swipe-row__action--edit` - Blue edit background

---

### Phase 6: Toast Notifications

**Goal:** Non-blocking feedback for actions

**Files to modify:**
- `web/static/min.css` - Toast styles
- `templates/user_dashboard.html` - Toast container + JS

**Implementation:**
```
- Position: bottom center, above FAB
- Auto-dismiss after 3 seconds
- Slide-up animation on appear
- Types: success, error, info
- Undo action for destructive operations
```

**CSS additions:**
- `.toast-container` - Fixed positioning
- `.toast` - Base toast styles
- `.toast--success`, `.toast--error` - Variants
- `.toast__action` - Undo button

---

### Phase 7: Pull-to-Refresh

**Goal:** Native-feeling data refresh

**Files to modify:**
- `web/static/min.css` - Pull indicator styles
- `templates/user_dashboard.html` - Pull-to-refresh handler

**Implementation:**
```
- Pull down on main content area
- Spinner indicator at top
- Triggers HTMX refresh of dynamic sections
- Haptic feedback on trigger threshold
```

---

### Phase 8: PWA Enhancements

**Goal:** Install-to-homescreen, offline support

**Files to create/modify:**
- `web/static/manifest.json` - PWA manifest
- `web/static/sw.js` - Service worker
- `templates/*.html` - Add manifest link, service worker registration
- `web/static/icons/icon-192.png` - App icon
- `web/static/icons/icon-512.png` - App icon
- `web/static/icons/apple-touch-icon.png` - iOS icon

**manifest.json:**
```json
{
  "name": "Peso - Weight Tracker",
  "short_name": "Peso",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#F9FAFB",
  "theme_color": "#111111",
  "icons": [...]
}
```

**Service Worker Strategy:**
```
- Cache static assets (CSS, JS, icons) on install
- Network-first for API calls with cache fallback
- IndexedDB for offline weight entries
- Background sync when connection restored
- Update prompt for new versions
```

**Offline UX:**
```
- Subtle offline indicator in topbar
- Weight entries queued with pending icon
- "Saved offline - will sync" toast
- Sync status in recent weights list
```

---

## File Change Summary

| File | Changes |
|------|---------|
| `web/static/min.css` | FAB, bottom sheet, stat cards, chips, swipe, toast, pull-refresh styles |
| `templates/user_dashboard.html` | New layout, FAB, bottom sheet, chips, swipe JS, toast JS |
| `templates/weight_form.html` | Adapt for bottom sheet context |
| `templates/partials_goal_summary.html` | Compact stat pill format |
| `templates/partials_recent_weights.html` | Swipeable row structure |
| `web/static/manifest.json` | NEW: PWA manifest |
| `web/static/sw.js` | NEW: Service worker |
| `web/static/icons/` | NEW: App icons |

---

## Implementation Order

All 8 phases will be implemented in sequence:

1. **FAB + Bottom Sheet** (Phase 1-2) - Core mobile entry pattern
2. **Dashboard Redesign** (Phase 3) - New summary cards layout
3. **Chart Chips** (Phase 4) - Touch-friendly period selection
4. **Swipe Actions** (Phase 5) - Quick edit/delete
5. **Toast Notifications** (Phase 6) - Action feedback
6. **Pull-to-Refresh** (Phase 7) - Native refresh pattern
7. **PWA Setup** (Phase 8) - Installable app with offline support

---

## Technical Considerations

- All animations use `transform` and `opacity` (GPU accelerated)
- Maintain HTMX architecture - no heavy JS frameworks
- Touch events with passive listeners for scroll performance
- CSS-only where possible, minimal vanilla JS
- Test on iOS Safari and Android Chrome
