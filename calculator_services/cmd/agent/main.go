package main

import (
	"go_calculator/internal/agent"
	"log/slog"
	"os"
)

var TotalWorkers = 4

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	for _ = range TotalWorkers {
		go agent.RunWorker()
	}
	waitUntilTermination()
}

func waitUntilTermination() {
	exitSignal := make(chan os.Signal, 1)
	//signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
