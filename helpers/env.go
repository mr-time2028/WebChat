package helpers

import "os"

// GetEnvOrDefaultString read string data from env file
func GetEnvOrDefaultString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
