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
		"0.0.0.0",
	)
	config.orchestratorPort = common.EnvOrDefault(
		"ORCHESTRATOR_GRPC_PORT",
		"8081",
	)

	workers := common.EnvOrDefault("COMPUTING_POWER", "4")
	config.TotalWorkers, _ = strconv.Atoi(workers)

	return config
}
