package utils

import (
	"os"
	"strconv"
)

// GetEnv returns a string based on the OS environment variable, and returns a default value if not found
func GetEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

// GetEnvBool returns a boolean based on the OS environment variable, and returns false if not found
func GetEnvBool(key string) (envValBool bool) {
	if envVal, ok := os.LookupEnv(key); ok {
		envValBool, _ = strconv.ParseBool(envVal)
	}
	return
}

// GetEnvInt retuns an integer based on the OS environment variable, and returns a default value if not found
func GetEnvInt(key string, defaultVal int) int {
	if envVal, ok := os.LookupEnv(key); ok {
		if val, ok := strconv.ParseInt(envVal, 0, 0); ok == nil {
			return int(val)
		}
	}
	return defaultVal
}
