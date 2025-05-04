package main

import (
	"log/slog"
	"os"

	"github.com/dzherb/go_calculator/internal/agent"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	app := agent.New()
	app.Run()
}
