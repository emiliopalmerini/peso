package main

import (
	"log"
	"os"

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

	if err := application.Run(); err != nil {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
