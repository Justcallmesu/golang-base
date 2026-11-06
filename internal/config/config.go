package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}
