package web

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	assets "peso"
	"peso/internal/application"
	"peso/internal/infrastructure/logging"
	"peso/internal/infrastructure/middleware"
	"peso/internal/interfaces"
)

type mergedFileServer struct {
	fs1, fs2 http.FileSystem
}

func (m *mergedFileServer) Open(name string) (http.File, error) {
	if m.fs1 != nil {
		if f, err := m.fs1.Open(name); err == nil {
			return f, nil
		}
	}
	if m.fs2 != nil {
		return m.fs2.Open(name)
	}
	return nil, fmt.Errorf("file not found: %s", name)
}

func (m *mergedFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.FileServer(m).ServeHTTP(w, r)
}

func NewRouter(
	weightTracker *application.WeightTracker,
	goalTracker *application.GoalTracker,
	authService *application.AuthService,
	userRepo interfaces.UserRepository,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	handlers := NewHandlers(weightTracker, goalTracker, userRepo, logger)
	authHandlers := NewAuthHandlers(authService, logger)

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /ready", readyHandler)

	registerStaticRoutes(mux, logger)

	mux.HandleFunc("GET /login", authHandlers.LoginPageHandler)
	mux.HandleFunc("POST /login", authHandlers.LoginHandler)
	mux.HandleFunc("GET /register", authHandlers.RegisterPageHandler)
	mux.HandleFunc("POST /register", authHandlers.RegisterHandler)
	mux.HandleFunc("GET /set-password", authHandlers.SetPasswordPageHandler)
	mux.HandleFunc("POST /set-password", authHandlers.SetPasswordHandler)
	mux.HandleFunc("POST /logout", authHandlers.LogoutHandler)
	mux.HandleFunc("GET /logout", authHandlers.LogoutHandler)

	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("GET /users/{userID}", handlers.UserDashboardHandler)
	mux.HandleFunc("GET /users/{userID}/recent-weights", handlers.RecentWeightsHandler)
	mux.HandleFunc("GET /users/{userID}/weight-form", handlers.WeightFormHandler)
	mux.HandleFunc("GET /users/{userID}/goal-form", handlers.GoalFormHandler)
	mux.HandleFunc("GET /users/{userID}/goal-summary", handlers.GoalSummaryHandler)
	mux.HandleFunc("GET /users/{userID}/goal-badge", handlers.GoalBadgeHandler)
	mux.HandleFunc("GET /users/{userID}/stat-hero", handlers.StatHeroHandler)
	mux.HandleFunc("GET /users/{userID}/stat-pills", handlers.StatPillsHandler)

	mux.HandleFunc("POST /api/weights", handlers.AddWeightHandler)
	mux.HandleFunc("DELETE /api/weights/{userID}/{weightID}", handlers.DeleteWeightHandler)
	mux.HandleFunc("GET /api/weights/{userID}", handlers.WeightHistoryHandler)
	mux.HandleFunc("GET /api/weights/latest/{userID}", handlers.WeightLatestHandler)
	mux.HandleFunc("POST /api/goals", handlers.AddGoalHandler)

	var handler http.Handler = mux
	handler = middleware.SessionMiddleware(authService)(handler)
	handler = logging.RequestLogger(logger)(handler)
	handler = logging.Recoverer(logger)(handler)
	handler = logging.RequestID(handler)

	return handler
}

func registerStaticRoutes(mux *http.ServeMux, logger *slog.Logger) {
	var fs1, fs2 http.FileSystem

	if staticSub, err := fs.Sub(assets.FS, "static"); err == nil {
		fs1 = http.FS(staticSub)
		logger.Info("static_fs_loaded", "source", "static/")
	} else {
		logger.Error("static_fs_failed", "error", err)
	}

	if webStaticSub, err := fs.Sub(assets.FS, "web/static"); err == nil {
		fs2 = http.FS(webStaticSub)
		logger.Info("web_static_fs_loaded", "source", "web/static/")
	} else {
		logger.Error("web_static_fs_failed", "error", err)
	}

	staticHandler := &mergedFileServer{fs1, fs2}
	mux.HandleFunc("GET /static/{filepath...}", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", staticHandler).ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ready"}`))
}
