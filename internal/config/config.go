package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := godotenv.Load()

	if err != nil {
		return fmt.Errorf("Error loading .env file: %w", err)
	}
	return nil
}
