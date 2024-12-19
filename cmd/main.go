package main

import (
	"go_calculator/internal/application"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	app := application.New()
	err := app.RunServer()
	if err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
	}
}
