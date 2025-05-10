package orchestrator

import (
	"os"
	"strconv"
	"time"

	"github.com/dzherb/go_calculator/calculator/internal/pkg"
)

type Config struct {
	Host               string
	Port               string
	GRPCPort           string
	AdditionTime       time.Duration
	SubtractionTime    time.Duration
	MultiplicationTime time.Duration
	DivisionTime       time.Duration
	TaskMaxProcessTime time.Duration
	SecretKey          string
	AccessTokenTTL     time.Duration
}

func ConfigFromEnv() *Config {
	config := new(Config)

	config.Host = common.EnvOrDefault("ORCHESTRATOR_HOST", "0.0.0.0")
	config.Port = common.EnvOrDefault("ORCHESTRATOR_HTTP_PORT", "8080")
	config.GRPCPort = common.EnvOrDefault("ORCHESTRATOR_GRPC_PORT", "8081")

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

	config.SecretKey = common.EnvOrDefault("SECRET_KEY", "insecure")

	if accessTokenTTL, exists := os.LookupEnv("ACCESS_TOKEN_TTL"); exists {
		config.AccessTokenTTL = getDurationInMin(accessTokenTTL)
	}

	return config
}

func getDurationInMs(duration string) time.Duration {
	t, _ := strconv.Atoi(duration)
	return time.Duration(t) * time.Millisecond
}

func getDurationInMin(duration string) time.Duration {
	t, _ := strconv.Atoi(duration)
	return time.Duration(t) * time.Minute
}
