package common

import "os"

// EnvOrDefault reads an environment variable or returns a default value.
func EnvOrDefault(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
