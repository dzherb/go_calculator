package agent

import (
	"strconv"

	"github.com/dzherb/go_calculator/internal/pkg"
)

type Config struct {
	orchestratorHost string
	orchestratorPort string
	TotalWorkers     int
}

func ConfigFromEnv() *Config {
	config := new(Config)

	config.orchestratorHost = common.EnvOrDefault(
		"ORCHESTRATOR_HOST",
		"localhost",
	)
	config.orchestratorPort = common.EnvOrDefault("ORCHESTRATOR_PORT", "8080")

	workers := common.EnvOrDefault("COMPUTING_POWER", "4")
	config.TotalWorkers, _ = strconv.Atoi(workers)

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
