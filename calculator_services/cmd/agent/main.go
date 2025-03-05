package main

import (
	"go_calculator/internal/agent"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	app := agent.New()
	app.Run()
}
