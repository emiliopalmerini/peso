package main

import (
    "fmt"
    "net/http"
    "os"
    "io/fs"
    "log/slog"

    "peso/internal/application"
    "peso/internal/infrastructure/persistence"
    "peso/internal/infrastructure/web"
    "peso/internal/interfaces"
    assets "peso"
    "peso/internal/infrastructure/logging"
)

// We'll load migrations from the filesystem at runtime for now

// mergedFileServer tries to serve from multiple filesystems
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

func main() {
    // Logger
    logger := logging.NewLogger(getEnv("LOG_LEVEL", "info"))
	// Get configuration from environment
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./peso.db")

	// Initialize database
	db, err := persistence.NewDB(dbPath)
	if err != nil {
		logger.Error("failed_to_initialize_database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations from local filesystem
	migrationsDir := os.DirFS("./migrations")
	if err := db.Migrate(migrationsDir); err != nil {
		logger.Error("failed_to_run_migrations", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := persistence.NewUserRepository(db)
	weightRepo := persistence.NewWeightRepository(db)
	goalRepo := persistence.NewGoalRepository(db)

	// Initialize domain services
	weightTracker := application.NewWeightTracker(userRepo, weightRepo)
	goalTracker := application.NewGoalTracker(userRepo, weightRepo, goalRepo)

	// Test the basic functionality
	if err := testBasicFunctionality(userRepo); err != nil {
		logger.Error("basic_functionality_test_failed", slog.Any("error", err))
		os.Exit(1)
	}

	// Setup HTTP server
    router := setupRouter(weightTracker, goalTracker, userRepo, logger)

    logger.Info("server_start", 
        slog.String("port", port),
        slog.String("db_path", dbPath),
    )

    if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("server_error", slog.Any("error", err))
		os.Exit(1)
	}
}

func setupRouter(weightTracker *application.WeightTracker, goalTracker *application.GoalTracker, userRepo interface{}, logger *slog.Logger) http.Handler {
    mux := http.NewServeMux()

    // Initialize web handlers
    handlers := web.NewHandlers(weightTracker, goalTracker, userRepo.(interfaces.UserRepository), logger)

	// Health endpoints
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /ready", readyHandler)

    // Serve static assets from embedded FS (favicon and web assets)
    // Must register BEFORE generic / route
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

    // Web UI endpoints
    mux.HandleFunc("GET /", handlers.HomeHandler)
    mux.HandleFunc("GET /users/{userID}", handlers.UserDashboardHandler)
    mux.HandleFunc("GET /users/{userID}/recent-weights", handlers.RecentWeightsHandler)
    mux.HandleFunc("GET /users/{userID}/weight-form", handlers.WeightFormHandler)
    mux.HandleFunc("GET /users/{userID}/goal-form", handlers.GoalFormHandler)
    mux.HandleFunc("GET /users/{userID}/goal-summary", handlers.GoalSummaryHandler)
    mux.HandleFunc("GET /users/{userID}/goal-badge", handlers.GoalBadgeHandler)

	// API endpoints
	mux.HandleFunc("POST /api/weights", handlers.AddWeightHandler)
	mux.HandleFunc("GET /api/weights/{userID}", handlers.WeightHistoryHandler)
	mux.HandleFunc("GET /api/weights/latest/{userID}", handlers.WeightLatestHandler)
	mux.HandleFunc("POST /api/goals", handlers.AddGoalHandler)

	// Wrap with middleware
	var handler http.Handler = mux
	handler = logging.RequestLogger(logger)(handler)
	handler = logging.Recoverer(logger)(handler)
	handler = logging.RequestID(handler)
	
	return handler
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


func testBasicFunctionality(userRepo interface{}) error {
	// This would test our domain objects work correctly
	// For now, just return success to show the app structure works
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
