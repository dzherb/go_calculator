package main

import (
	"log/slog"
	"os"

	"github.com/dzherb/go_calculator/internal/orchestrator"
	"github.com/dzherb/go_calculator/internal/storage"
	"github.com/dzherb/go_calculator/pkg/logger"
)

func main() {
	logger.Init()

	closeSt, err := storage.InitFromEnv()
	if err != nil {
		slog.Error(
			"Failed to initialize storage",
			"error", err,
		)

		return
	}
	defer closeSt()

	app := orchestrator.New()

	err = app.Serve()
	if err != nil {
		defer os.Exit(1)
	}
}
