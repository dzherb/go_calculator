package orchestrator

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr               string
	Port               string
	AdditionTime       time.Duration
	SubtractionTime    time.Duration
	MultiplicationTime time.Duration
	DivisionTime       time.Duration
	TaskMaxProcessTime time.Duration
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("ORCHESTRATOR_HOST")
	if config.Addr == "" {
		config.Addr = "127.0.0.1"
	}
	config.Port = os.Getenv("ORCHESTRATOR_PORT")
	if config.Port == "" {
		config.Port = "8080"
	}

	if addTime, exists := os.LookupEnv("TIME_ADDITION_MS"); exists {
		config.AdditionTime = getDurationInMs(addTime)
	}
	if subTime, exists := os.LookupEnv("TIME_SUBTRACTION_MS"); exists {
		config.SubtractionTime = getDurationInMs(subTime)
	}
	if mulTime, exists := os.LookupEnv("TIME_MULTIPLICATIONS_MS"); exists {
		config.MultiplicationTime = getDurationInMs(mulTime)
	}
	if divTime, exists := os.LookupEnv("TIME_DIVISIONS_MS"); exists {
		config.DivisionTime = getDurationInMs(divTime)
	}

	if maxTime, exists := os.LookupEnv("TASK_MAX_PROCESS_TIME_IN_MS"); exists {
		config.TaskMaxProcessTime = getDurationInMs(maxTime)
	} else {
		config.TaskMaxProcessTime = 1 * time.Minute
	}

	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	app := &Application{
		config: ConfigFromEnv(),
	}
	orchestrator.app = app
	return app
}

func getDurationInMs(duration string) time.Duration {
	t, _ := strconv.Atoi(duration)
	return time.Duration(t) * time.Millisecond
}
