package web

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log/slog"
    "net/http"
    "strconv"
    "strings"
    "time"

    "peso/internal/application"
    "peso/internal/domain/user"
    "peso/internal/domain/weight"
    "peso/internal/domain/goal"
    assets "peso"

    "github.com/gorilla/mux"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	weightTracker *application.WeightTracker
	goalTracker   *application.GoalTracker
	templates     *template.Template
	logger        *slog.Logger
}

// NewHandlers creates new web handlers
func NewHandlers(weightTracker *application.WeightTracker, goalTracker *application.GoalTracker, logger *slog.Logger) *Handlers {
	return &Handlers{
		weightTracker: weightTracker,
		goalTracker:   goalTracker,
		templates:     loadTemplates(),
		logger:        logger,
	}
}

// HomeHandler serves the main dashboard
func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		Users []string
	}{
		Title: "Peso - Weight Tracking",
		Users: []string{"giada", "emilio"},
	}
	
    if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
        writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
        return
    }
}

// AddWeightHandler handles weight recording
func (h *Handlers) AddWeightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse form data
	userIDStr := r.FormValue("user_id")
	weightStr := r.FormValue("weight")
	unitStr := r.FormValue("unit")
	dateStr := r.FormValue("date")
	notes := r.FormValue("notes")
	
	// Validate inputs
	if userIDStr == "" || weightStr == "" || unitStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	
	// Parse weight value
	weightFloat, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		http.Error(w, "Invalid weight value", http.StatusBadRequest)
		return
	}
	
    // Parse date and attach current time (hour/min/sec) automatically
    var measuredAt time.Time
    if dateStr == "" {
        measuredAt = time.Now()
    } else {
        // Parse provided date in local timezone
        d, err2 := time.ParseInLocation("2006-01-02", dateStr, time.Local)
        if err2 != nil {
            http.Error(w, "Invalid date format", http.StatusBadRequest)
            return
        }
        now := time.Now()
        measuredAt = time.Date(d.Year(), d.Month(), d.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
    }
	
	// Create domain objects
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	weightValue, err := weight.NewWeightValue(weightFloat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid weight: %v", err), http.StatusBadRequest)
		return
	}
	
	unit, err := weight.NewWeightUnit(unitStr)
	if err != nil {
		http.Error(w, "Invalid unit", http.StatusBadRequest)
		return
	}
	
	// Record weight using domain service  
	recordedWeight, err := h.weightTracker.RecordWeight(userID, weightValue, unit, measuredAt, notes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to record weight: %v", err), http.StatusBadRequest)
		return
	}
	
	// Return success response (HTMX will handle this)
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Weight  struct {
			ID     string  `json:"id"`
			Value  float64 `json:"value"`
			Unit   string  `json:"unit"`
			Date   string  `json:"date"`
		} `json:"weight"`
	}{
		Success: true,
		Message: "Weight recorded successfully",
		Weight: struct {
			ID     string  `json:"id"`
			Value  float64 `json:"value"`
			Unit   string  `json:"unit"`
			Date   string  `json:"date"`
		}{
			ID:    recordedWeight.ID().String(),
			Value: recordedWeight.Value().Float64(),
			Unit:  recordedWeight.Unit().String(),
			Date:  recordedWeight.MeasuredAt().Format("2006-01-02"),
		},
	}
	
	json.NewEncoder(w).Encode(response)
}

// WeightHistoryHandler returns weight history for a user
func (h *Handlers) WeightHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	// Get period from query params (default to last month)
	period := application.TimePeriodLastMonth
	if periodStr := r.URL.Query().Get("period"); periodStr != "" {
		switch periodStr {
		case "week":
			period = application.TimePeriodLastWeek
		case "month":
			period = application.TimePeriodLastMonth
		case "3months":
			period = application.TimePeriodLast3Months
		case "6months":
			period = application.TimePeriodLast6Months
		case "year":
			period = application.TimePeriodLastYear
		case "all":
			period = application.TimePeriodAll
		}
	}
	
	weights, err := h.weightTracker.GetWeightHistory(userID, period)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get weight history: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Convert to JSON-friendly format
    type WeightResponse struct {
        ID     string  `json:"id"`
        Value  float64 `json:"value"`
        Unit   string  `json:"unit"`
        Date   string  `json:"date"`
        Time   string  `json:"time"`
        Notes  string  `json:"notes"`
    }
	
	var response []WeightResponse
	for _, w := range weights {
        response = append(response, WeightResponse{
            ID:    w.ID().String(),
            Value: w.Value().Float64(),
            Unit:  w.Unit().String(),
            Date:  w.MeasuredAt().Format("2006-01-02"),
            Time:  w.MeasuredAt().Format("15:04"),
            Notes: w.Notes(),
        })
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// WeightLatestHandler returns the latest weight for a user
func (h *Handlers) WeightLatestHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["userID"]

    userID, err := user.NewUserID(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    latest, err := h.weightTracker.GetLatestWeight(userID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to get latest weight: %v", err), http.StatusInternalServerError)
        return
    }

    resp := struct {
        ID    string  `json:"id"`
        Value float64 `json:"value"`
        Unit  string  `json:"unit"`
        Date  string  `json:"date"`
        Time  string  `json:"time"`
        Notes string  `json:"notes"`
    }{
        ID:    latest.ID().String(),
        Value: latest.Value().Float64(),
        Unit:  latest.Unit().String(),
        Date:  latest.MeasuredAt().Format("2006-01-02"),
        Time:  latest.MeasuredAt().Format("15:04"),
        Notes: latest.Notes(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// UserDashboardHandler serves individual user dashboard
func (h *Handlers) UserDashboardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
    // Get active goal if exists
    activeGoal, _ := h.goalTracker.GetActiveGoal(userID)
	
	// Calculate goal progress if goal exists
	var progress *application.GoalProgress
	if activeGoal != nil {
		p, err := h.goalTracker.CalculateProgress(userID)
		if err == nil {
			progress = &p
		}
	}
	
    data := struct {
        UserID     string
        ActiveGoal interface{}
        Progress   *application.GoalProgress
    }{
        UserID:     userIDStr,
        ActiveGoal: activeGoal,
        Progress:   progress,
    }
	
	if err := h.templates.ExecuteTemplate(w, "user_dashboard.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

// GoalFormHandler serves the goal entry form
func (h *Handlers) GoalFormHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["userID"]

    userID, err := user.NewUserID(userIDStr)
    if err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    // Try to get current weight for helper text
    latest, _ := h.weightTracker.GetLatestWeight(userID)
    var current struct{
        Value float64
        Unit  string
    }
    if latest != nil {
        current.Value = latest.Value().Float64()
        current.Unit = latest.Unit().String()
    }

    data := struct {
        UserID        string
        Today         string
        CurrentWeight *struct{ Value float64; Unit string }
    }{
        UserID:        userIDStr,
        Today:         time.Now().Format("2006-01-02"),
        CurrentWeight: func() *struct{ Value float64; Unit string } { if latest != nil { return &current }; return nil }(),
    }

    if err := h.templates.ExecuteTemplate(w, "goal_form.html", data); err != nil {
        writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
        return
    }
}

// AddGoalHandler handles goal creation
func (h *Handlers) AddGoalHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeError(h.logger, w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    userIDStr := r.FormValue("user_id")
    goalType := r.FormValue("goal_type") // optional for now
    targetWeightStr := r.FormValue("target_weight")
    unitStr := r.FormValue("unit")
    targetDateStr := r.FormValue("target_date")
    notes := r.FormValue("notes")

    if userIDStr == "" || targetWeightStr == "" || unitStr == "" || targetDateStr == "" {
        writeError(h.logger, w, r, http.StatusBadRequest, "Missing required fields", nil)
        return
    }

    userID, err := user.NewUserID(userIDStr)
    if err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    tw, err := strconv.ParseFloat(targetWeightStr, 64)
    if err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Invalid target weight", err)
        return
    }
    targetWeight, err := weight.NewWeightValue(tw)
    if err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Invalid target weight", err)
        return
    }

    unit, err := weight.NewWeightUnit(unitStr)
    if err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Invalid unit", err)
        return
    }

    // Parse target date (YYYY-MM-DD)
    t, err := time.Parse("2006-01-02", targetDateStr)
    if err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Invalid target date", err)
        return
    }
    td, err := goal.NewTargetDate(t.Year(), int(t.Month()), t.Day())
    if err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Invalid target date", err)
        return
    }

    // Optionally, enforce direction for goalType (not used by domain yet)
    _ = goalType

    if _, err := h.goalTracker.SetGoal(userID, targetWeight, unit, td, notes); err != nil {
        writeError(h.logger, w, r, http.StatusBadRequest, "Failed to set goal", err)
        return
    }

    // HTMX redirect to refresh dashboard with new goal
    w.Header().Set("HX-Redirect", "/users/"+userIDStr)
    w.WriteHeader(http.StatusNoContent)
}
// RecentWeightsHandler returns the HTML partial with recent weights list
func (h *Handlers) RecentWeightsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["userID"]

    userID, err := user.NewUserID(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    weights, err := h.weightTracker.GetRecentWeights(userID, 10)
    if err != nil {
        http.Error(w, "Failed to load recent weights", http.StatusInternalServerError)
        return
    }

    // Build view models
    type Row struct {
        Date  string
        Value string
        Unit  string
        Notes string
    }
    var rows []Row
    for _, wgt := range weights {
        rows = append(rows, Row{
            Date:  wgt.MeasuredAt().Format("02/01/2006"),
            Value: fmt.Sprintf("%.1f", wgt.Value().Float64()),
            Unit:  wgt.Unit().String(),
            Notes: wgt.Notes(),
        })
    }

    data := struct {
        Rows []Row
    }{Rows: rows}

    if err := h.templates.ExecuteTemplate(w, "partials_recent_weights.html", data); err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }
}

// WeightFormHandler serves the weight entry form
func (h *Handlers) WeightFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	
	data := struct {
		UserID string
		Today  string
	}{
		UserID: userIDStr,
		Today:  time.Now().Format("2006-01-02"),
	}
	
	if err := h.templates.ExecuteTemplate(w, "weight_form.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

// loadTemplates loads HTML templates from files
func loadTemplates() *template.Template {
    tmpl := template.New("").Funcs(template.FuncMap{
        "title": func(s string) string {
            if len(s) == 0 {
                return s
            }
            return string(s[0]-32) + s[1:]
        },
        // toJson marshals a value to JSON for safe inline JS usage in templates
        "toJson": func(v interface{}) template.JS {
            b, err := json.Marshal(v)
            if err != nil {
                return template.JS("null")
            }
            return template.JS(string(b))
        },
    })

    // Load templates from embedded filesystem
    template.Must(tmpl.ParseFS(assets.FS, "templates/*.html"))

    return tmpl
}

// Helper function to convert typed slice to interface slice for templates
func interfaceSlice(slice interface{}) []interface{} {
    // This is a simplified implementation
    // In a real app, you'd handle this more robustly
    return []interface{}{}
}

// writeError writes a uniform error structure and logs it
func writeError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, status int, message string, err error) {
    reqID := r.Header.Get("X-Request-ID")
    if err != nil {
        logger.Error("http_error",
            slog.Int("status", status),
            slog.String("message", message),
            slog.Any("error", err),
            slog.String("path", r.URL.Path),
            slog.String("request_id", reqID),
        )
    } else {
        logger.Warn("http_error",
            slog.Int("status", status),
            slog.String("message", message),
            slog.String("path", r.URL.Path),
            slog.String("request_id", reqID),
        )
    }

    // JSON for /api, HTML/plain otherwise
    isAPI := strings.HasPrefix(r.URL.Path, "/api/")
    if isAPI {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        _ = json.NewEncoder(w).Encode(map[string]any{
            "success":    false,
            "error":      http.StatusText(status),
            "message":    message,
            "request_id": reqID,
        })
        return
    }
    http.Error(w, message, status)
}
