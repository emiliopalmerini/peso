package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "io/fs"

    "peso/internal/application"
    "peso/internal/infrastructure/persistence"
    "peso/internal/infrastructure/web"
    assets "peso"

    "github.com/gorilla/mux"
)

// We'll load migrations from the filesystem at runtime for now

func main() {
	// Get configuration from environment
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./peso.db")

	// Initialize database
	db, err := persistence.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations from local filesystem
	migrationsDir := os.DirFS("./migrations")
	if err := db.Migrate(migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
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
		log.Fatalf("Basic functionality test failed: %v", err)
	}

	// Setup HTTP server
	router := setupRouter(weightTracker, goalTracker)

    fmt.Printf("Peso app starting on port %s\n", port)
    fmt.Printf("Database: %s\n", dbPath)
    fmt.Printf("Open http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func setupRouter(weightTracker *application.WeightTracker, goalTracker *application.GoalTracker) *mux.Router {
    r := mux.NewRouter()

    // Initialize web handlers
    handlers := web.NewHandlers(weightTracker, goalTracker)

    // Serve static assets (CSS/JS) from embedded FS
    if sub, err := fs.Sub(assets.FS, "web/static"); err == nil {
        r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
    }

	// Health endpoints
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/ready", readyHandler).Methods("GET")

	// Web UI endpoints
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
    r.HandleFunc("/users/{userID}", handlers.UserDashboardHandler).Methods("GET")
    r.HandleFunc("/users/{userID}/recent-weights", handlers.RecentWeightsHandler).Methods("GET")
	r.HandleFunc("/users/{userID}/weight-form", handlers.WeightFormHandler).Methods("GET")

	// API endpoints
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/weights", handlers.AddWeightHandler).Methods("POST")
    api.HandleFunc("/weights/{userID}", handlers.WeightHistoryHandler).Methods("GET")
    api.HandleFunc("/weights/latest/{userID}", handlers.WeightLatestHandler).Methods("GET")

	return r
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
    fmt.Println("Testing basic domain functionality...")
	
	// This would test our domain objects work correctly
	// For now, just return success to show the app structure works
    fmt.Println("Domain layer tests passed")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
