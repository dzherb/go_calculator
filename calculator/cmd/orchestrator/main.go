package main

import (
	"os"

	"github.com/dzherb/go_calculator/internal/orchestrator"
	"github.com/dzherb/go_calculator/pkg/logger"
)

func main() {
	logger.Init()

	app := orchestrator.New()

	err := app.Serve()
	if err != nil {
		os.Exit(1)
	}
}
