package agent

import (
	"os"
	"strconv"
)

type Config struct {
	orchestratorHost string
	orchestratorPort string
	TotalWorkers     int
}

func ConfigFromEnv() *Config {
	config := new(Config)

	orchestratorHost, exists := os.LookupEnv("ORCHESTRATOR_HOST")
	if !exists {
		orchestratorHost = "localhost"
	}
	config.orchestratorHost = orchestratorHost

	orchestratorPort, exists := os.LookupEnv("ORCHESTRATOR_PORT")
	if !exists {
		orchestratorPort = "8080"
	}
	config.orchestratorPort = orchestratorPort

	workers, exists := os.LookupEnv("COMPUTING_POWER")
	if !exists {
		config.TotalWorkers = 4
	} else {
		config.TotalWorkers, _ = strconv.Atoi(workers)
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}
