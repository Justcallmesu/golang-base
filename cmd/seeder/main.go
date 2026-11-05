package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
	"justcallmesu.com/golang-base/internal/database/seeds"
)

func main() {

	config.LoadConfig()

	database, databaseError := database.InitConnection()

	if databaseError != nil {
		slog.Error("Error initializing database connection", "error", databaseError)
		os.Exit(1)
	}

	defer func() {

		sqlDB, sqlDbError := database.DB()

		if sqlDbError != nil {
			slog.Error("Error getting database instance", "error", sqlDbError)
			os.Exit(1)
		}

		sqlDbError = sqlDB.Close()

		if sqlDbError != nil {
			slog.Error("Error closing database connection", "error", sqlDbError)
			os.Exit(1)
		}

		slog.Info("Database connection closed successfully")
	}()

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.StampMilli,
	})

	seeds.PermissionSeeder(database, logger)
	seeds.RoleSeeder(database, logger)
	seeds.RolePermissionSeeder(database, logger)
	seeds.UserSeeder(database, logger)

}
