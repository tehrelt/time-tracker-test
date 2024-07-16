package main

import (
	"em-test/internal/app"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", slog.String("err", err.Error()))
		os.Exit(-1)
	}

	a, rollback, err := app.Init()
	if err != nil {
		slog.Error("Error initializing app", slog.String("err", err.Error()))
		os.Exit(-1)
	}
	defer rollback()
	if err := a.Run(); err != nil {
		slog.Error("Error running app", slog.String("err", err.Error()))
		os.Exit(-1)
	}
}
