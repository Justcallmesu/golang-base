package main

import (
	"log/slog"
	"os"

	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
)

func main() {

	config.LoadConfig()

	dbConnection, databaseError := database.InitConnection()

	if databaseError != nil {
		slog.Error("Error initializing database connection", "error", databaseError)
		os.Exit(1)
	}

	migrationError := database.Migrations(dbConnection)

	if migrationError != nil {
		slog.Error("Database migration failed", "error", migrationError)
		os.Exit(1)
	}
}
