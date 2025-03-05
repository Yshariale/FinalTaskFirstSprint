package main

import (
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/application"
)

func main() {
	slog.Info("Starting application")
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
	app := application.New()
	if err := app.RunServer(); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
