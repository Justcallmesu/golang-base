package utils

import "os"

func GetEnv(lookupKey string, defaultValue string) string {
	envValue := os.Getenv(lookupKey)

	if envValue == "" {
		return defaultValue
	}
	return envValue
}
