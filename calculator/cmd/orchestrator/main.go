package main

import (
	"log/slog"
	"os"

	"github.com/dzherb/go_calculator/internal/orchestrator"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	app := orchestrator.New()

	err := app.RunServer()
	if err != nil {
		slog.Error("server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
