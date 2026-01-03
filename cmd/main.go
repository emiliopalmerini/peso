package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"peso/internal/app"
	"peso/internal/config"
)

func main() {
	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		log.Printf("failed to initialize application: %v", err)
		os.Exit(1)
	}
	defer application.Close()

	errChan := make(chan error, 1)
	go func() {
		errChan <- application.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("server error: %v", err)
			os.Exit(1)
		}
	case sig := <-quit:
		log.Printf("received signal: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
		os.Exit(1)
	}
}
