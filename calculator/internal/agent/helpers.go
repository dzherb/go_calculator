package agent

import (
	"os"
	"os/signal"
	"syscall"
)

func waitUntilTermination() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
