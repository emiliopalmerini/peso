package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"peso/internal/application"
	"peso/internal/config"
	"peso/internal/infrastructure/logging"
	"peso/internal/infrastructure/persistence"
	"peso/internal/infrastructure/web"
)

type App struct {
	config *config.Config
	logger *slog.Logger
	db     *persistence.DB
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	logger := logging.NewLogger(cfg.LogLevel)

	db, err := persistence.NewDB(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	migrationsDir := os.DirFS("./migrations")
	if err := db.Migrate(migrationsDir); err != nil {
		db.Close()
		return nil, err
	}

	userRepo := persistence.NewUserRepository(db)
	weightRepo := persistence.NewWeightRepository(db)
	goalRepo := persistence.NewGoalRepository(db)
	sessionRepo := persistence.NewSessionRepository(db)

	weightTracker := application.NewWeightTracker(userRepo, weightRepo)
	goalTracker := application.NewGoalTracker(userRepo, weightRepo, goalRepo)
	authService := application.NewAuthService(userRepo, sessionRepo)

	if err := authService.CleanupExpiredSessions(); err != nil {
		logger.Warn("failed_to_cleanup_sessions", slog.Any("error", err))
	}

	router := web.NewRouter(weightTracker, goalTracker, authService, userRepo, logger)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &App{
		config: cfg,
		logger: logger,
		db:     db,
		server: server,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("server_start",
		slog.String("port", a.config.Port),
		slog.String("db_path", a.config.DBPath),
	)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("server_shutdown")
	return a.server.Shutdown(ctx)
}

func (a *App) Close() error {
	return a.db.Close()
}
