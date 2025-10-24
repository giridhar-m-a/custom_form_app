package utils

import (
	"log"
	"os"
	"strconv"
)

// GetEnv retrieves the environment variable named by key, returning defaultValue if the variable is not set or is empty.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt retrieves the environment variable named by key and returns its value as an int.
// If the variable is not set or is empty, or if its value cannot be parsed as an integer, defaultValue is returned.
// A parse failure is logged as a warning that includes the key, the raw value, and the defaultValue used.
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid integer value for %s: '%s', using default: %d",
			key, value, defaultValue)
	}
	return defaultValue
}

// GetEnvAsBool retrieves the environment variable named by key and parses it as a boolean.
// If the variable is unset or empty, or if parsing fails, it returns defaultValue; parsing failures are logged as a warning.
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
		log.Printf("Warning: Invalid boolean value for %s: '%s', using default: %t",
			key, value, defaultValue)
	}
	return defaultValue
}