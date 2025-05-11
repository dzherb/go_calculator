package main

import (
	"github.com/dzherb/go_calculator/calculator/internal/agent"
	"github.com/dzherb/go_calculator/calculator/pkg/logger"
)

func main() {
	logger.Init()

	app := agent.New()
	app.Run()
}
