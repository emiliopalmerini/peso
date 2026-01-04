package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"peso/internal/application"
	"peso/internal/domain/goal"
	"peso/internal/domain/user"
	"peso/internal/domain/weight"
	"peso/internal/infrastructure/middleware"
	"peso/internal/interfaces"
	"strconv"
	"strings"
	"time"

	assets "peso"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	weightTracker *application.WeightTracker
	goalTracker   *application.GoalTracker
	userRepo      interfaces.UserRepository
	templates     *template.Template
	logger        *slog.Logger
}

// NewHandlers creates new web handlers
func NewHandlers(weightTracker *application.WeightTracker, goalTracker *application.GoalTracker, userRepo interfaces.UserRepository, logger *slog.Logger) *Handlers {
	return &Handlers{
		weightTracker: weightTracker,
		goalTracker:   goalTracker,
		userRepo:      userRepo,
		templates:     loadTemplates(),
		logger:        logger,
	}
}

// HomeHandler serves the landing page
func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	u := middleware.UserFromContext(r.Context())
	if u != nil {
		http.Redirect(w, r, "/users/"+u.ID().String(), http.StatusSeeOther)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Peso - Weight Tracking",
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
		return
	}
}

// AddWeightHandler handles weight recording
func (h *Handlers) AddWeightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(h.logger, w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Parse form data
	userIDStr := r.FormValue("user_id")
	weightStr := r.FormValue("weight")

	// Validate inputs
	if userIDStr == "" || weightStr == "" {
		writeError(h.logger, w, r, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	// Parse weight value
	weightFloat, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid weight value", err)
		return
	}

	// Use current server time
	measuredAt := time.Now()

	// Create domain objects
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	weightValue, err := weight.NewWeightValue(weightFloat)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid weight", err)
		return
	}

	unit, err := weight.NewWeightUnit("kg")
	if err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Invalid unit", err)
		return
	}

	// Record weight using domain service
	recordedWeight, err := h.weightTracker.RecordWeight(userID, weightValue, unit, measuredAt, "")
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Failed to record weight", err)
		return
	}

	// Return success response (HTMX will handle this)
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Weight  struct {
			ID    string  `json:"id"`
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
			Date  string  `json:"date"`
		} `json:"weight"`
	}{
		Success: true,
		Message: "Weight recorded successfully",
		Weight: struct {
			ID    string  `json:"id"`
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
			Date  string  `json:"date"`
		}{
			ID:    recordedWeight.ID().String(),
			Value: recordedWeight.Value().Float64(),
			Unit:  recordedWeight.Unit().String(),
			Date:  recordedWeight.MeasuredAt().Format("02/01/2006"),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// WeightHistoryHandler returns weight history for a user
func (h *Handlers) WeightHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")

	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
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
		writeError(h.logger, w, r, http.StatusInternalServerError, "Failed to get weight history", err)
		return
	}

	// Convert to JSON-friendly format
	type WeightResponse struct {
		ID    string  `json:"id"`
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
		Date  string  `json:"date"`
		Time  string  `json:"time"`
		Notes string  `json:"notes"`
	}

	var response []WeightResponse
	for _, w := range weights {
		response = append(response, WeightResponse{
			ID:    w.ID().String(),
			Value: w.Value().Float64(),
			Unit:  w.Unit().String(),
			Date:  w.MeasuredAt().Format("02/01/2006"),
			Time:  w.MeasuredAt().Format("15:04"),
			Notes: w.Notes(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// WeightLatestHandler returns the latest weight for a user
func (h *Handlers) WeightLatestHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")

	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	latest, err := h.weightTracker.GetLatestWeight(userID)
	if err != nil {
		writeError(h.logger, w, r, http.StatusNotFound, "No weight found", err)
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
		Date:  latest.MeasuredAt().Format("02/01/2006"), Time: latest.MeasuredAt().Format("15:04"),
		Notes: latest.Notes(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UserDashboardHandler serves individual user dashboard
func (h *Handlers) UserDashboardHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.UserFromContext(r.Context())
	if currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userIDStr := r.PathValue("userID")
	if userIDStr != currentUser.ID().String() {
		http.Redirect(w, r, "/users/"+currentUser.ID().String(), http.StatusSeeOther)
		return
	}

	userID := currentUser.ID()

	// Get active goal if exists
	activeGoal, _ := h.goalTracker.GetActiveGoal(userID)

	// Calculate goal progress if goal exists
	var progress *application.GoalProgress
	var startWeight interface{}
	var createdAt interface{}
	if activeGoal != nil {
		p, err := h.goalTracker.CalculateProgress(userID)
		if err == nil {
			progress = &p
		}
		
		// Get starting weight for trajectory calculation
		if startWeightRecord, err := h.goalTracker.GetStartingWeightForGoal(userID, activeGoal.CreatedAt()); err == nil {
			startWeight = startWeightRecord.Value().Float64()
		}
		createdAt = activeGoal.CreatedAt().Format("02/01/2006")
	}

	data := struct {
		UserID      string
		UserName    string
		ActiveGoal  interface{}
		Progress    *application.GoalProgress
		StartWeight interface{}
		CreatedAt   interface{}
	}{
		UserID:      userID.String(),
		UserName:    currentUser.Name(),
		ActiveGoal:  activeGoal,
		Progress:    progress,
		StartWeight: startWeight,
		CreatedAt:   createdAt,
	}

	if err := h.templates.ExecuteTemplate(w, "user_dashboard.html", data); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
		return
	}
}

// GoalFormHandler serves the goal entry form
func (h *Handlers) GoalFormHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")

	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Try to get current weight for helper text
	latest, _ := h.weightTracker.GetLatestWeight(userID)
	var current struct {
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
		CurrentWeight *struct {
			Value float64
			Unit  string
		}
	}{
		UserID: userIDStr,
		Today:  time.Now().Format("02/01/2006"),
		CurrentWeight: func() *struct {
			Value float64
			Unit  string
		} {
			if latest != nil {
				return &current
			}
			return nil
		}(),
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
	targetDateStr := r.FormValue("target_date")
	notes := r.FormValue("notes")

	if userIDStr == "" || targetWeightStr == "" || targetDateStr == "" {
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

	unit, err := weight.NewWeightUnit("kg")
	if err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Invalid unit", err)
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

	// Return simple success response for HTMX
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// RecentWeightsHandler returns the HTML partial with recent weights list
func (h *Handlers) RecentWeightsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")

	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	weights, err := h.weightTracker.GetRecentWeights(userID, 10)
	if err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Failed to load recent weights", err)
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
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
		return
	}
}

// WeightFormHandler serves the weight entry form
func (h *Handlers) WeightFormHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")

	data := struct {
		UserID string
	}{
		UserID: userIDStr,
	}

	if err := h.templates.ExecuteTemplate(w, "weight_form.html", data); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
		return
	}
}

// GoalSummaryHandler returns the goal summary partial
func (h *Handlers) GoalSummaryHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")

	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Check if user has any weight records
	hasWeights := false
	if _, err := h.weightTracker.GetLatestWeight(userID); err == nil {
		hasWeights = true
	}

	type vm struct {
		UserID        string
		Active        bool
		TargetWeight  string
		Unit          string
		TargetDate    string
		HasProgress   bool
		WeightToLose  string
		DaysRemaining int
		HasWeights    bool
	}

	out := vm{UserID: userIDStr, HasWeights: hasWeights}
	if g, _ := h.goalTracker.GetActiveGoal(userID); g != nil {
		out.Active = true
		out.TargetWeight = fmt.Sprintf("%.1f", g.TargetWeight().Float64())
		out.Unit = g.Unit().String()
		out.TargetDate = g.TargetDate().String()
		if p, err := h.goalTracker.CalculateProgress(userID); err == nil {
			out.HasProgress = true
			out.WeightToLose = fmt.Sprintf("%.1f", p.WeightToLose.Float64())
			out.DaysRemaining = p.DaysRemaining
		}
	}

	if err := h.templates.ExecuteTemplate(w, "partials_goal_summary.html", out); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
		return
	}
}

// GoalBadgeHandler returns just the small badge for the chart header
func (h *Handlers) GoalBadgeHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	type vm struct {
		Active        bool
		HasProgress   bool
		WeightToLose  string
		DaysRemaining int
	}
	out := vm{}
	if g, _ := h.goalTracker.GetActiveGoal(userID); g != nil {
		out.Active = true
		if p, err := h.goalTracker.CalculateProgress(userID); err == nil {
			out.HasProgress = true
			out.WeightToLose = fmt.Sprintf("%.1f", p.WeightToLose.Float64())
			out.DaysRemaining = p.DaysRemaining
		}
	}
}

// StatHeroHandler returns the hero stat card with current weight and trend
func (h *Handlers) StatHeroHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	type vm struct {
		HasData       bool
		CurrentWeight string
		Unit          string
		LastDate      string
		LastTime      string
		TrendValue    string
		TrendClass    string
	}

	out := vm{}

	latest, err := h.weightTracker.GetLatestWeight(userID)
	if err == nil && latest != nil {
		out.HasData = true
		out.CurrentWeight = fmt.Sprintf("%.1f", latest.Value().Float64())
		out.Unit = latest.Unit().String()
		out.LastDate = latest.MeasuredAt().Format("02/01")
		out.LastTime = latest.MeasuredAt().Format("15:04")

		// Calculate 7-day trend
		weights, _ := h.weightTracker.GetWeightHistory(userID, application.TimePeriodLastWeek)
		if len(weights) >= 2 {
			oldest := weights[len(weights)-1].Value().Float64()
			newest := weights[0].Value().Float64()
			diff := newest - oldest
			if diff < -0.1 {
				out.TrendValue = fmt.Sprintf("%.1f", diff)
				out.TrendClass = "stat-hero__trend--down"
			} else if diff > 0.1 {
				out.TrendValue = fmt.Sprintf("+%.1f", diff)
				out.TrendClass = "stat-hero__trend--up"
			} else {
				out.TrendValue = "0.0"
				out.TrendClass = "stat-hero__trend--neutral"
			}
		}
	}

	if err := h.templates.ExecuteTemplate(w, "partials_stat_hero.html", out); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
		return
	}
}

// StatPillsHandler returns the stat pills with goal and average
func (h *Handlers) StatPillsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		writeError(h.logger, w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	type vm struct {
		GoalWeight   string
		GoalUnit     string
		HasGoal      bool
		WeekAvg      string
		WeekAvgUnit  string
		HasWeekAvg   bool
	}

	out := vm{}

	// Get goal info
	if g, _ := h.goalTracker.GetActiveGoal(userID); g != nil {
		out.HasGoal = true
		out.GoalWeight = fmt.Sprintf("%.1f", g.TargetWeight().Float64())
		out.GoalUnit = g.Unit().String()
	}

	// Calculate 7-day average
	weights, _ := h.weightTracker.GetWeightHistory(userID, application.TimePeriodLastWeek)
	if len(weights) > 0 {
		var sum float64
		for _, wgt := range weights {
			sum += wgt.Value().Float64()
		}
		avg := sum / float64(len(weights))
		out.HasWeekAvg = true
		out.WeekAvg = fmt.Sprintf("%.1f", avg)
		out.WeekAvgUnit = weights[0].Unit().String()
	}

	if err := h.templates.ExecuteTemplate(w, "partials_stat_pills.html", out); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
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
