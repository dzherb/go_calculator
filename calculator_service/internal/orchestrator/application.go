package orchestrator

import "os"

type Config struct {
	Addr string
	Port string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("GO_CALC_ADDR")
	if config.Addr == "" {
		config.Addr = "127.0.0.1"
	}
	config.Port = os.Getenv("GO_CALC_PORT")
	if config.Port == "" {
		config.Port = "8080"
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
