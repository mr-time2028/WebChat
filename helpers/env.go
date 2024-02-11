package helpers

import (
	"os"
	"strconv"
)

// GetEnvOrDefaultString read string data from env file
func GetEnvOrDefaultString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvOrDefaultBool read bool data from env file
func GetEnvOrDefaultBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}
