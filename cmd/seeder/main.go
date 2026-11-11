package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"justcallmesu.com/golang-base/internal/common/logger"
	"justcallmesu.com/golang-base/internal/common/utils"
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
	"justcallmesu.com/golang-base/internal/database/seeds"
)

func main() {

	configError := config.LoadConfig()

	if configError != nil {
		slog.Error("Error loading configuration", "error", configError)
		os.Exit(1)
	}

	appEnv := utils.GetEnv("APP_ENV", "development")
	slog.SetDefault(logger.New(appEnv))

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

	seedLogger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.StampMilli,
	})

	seeds.PermissionSeeder(database, seedLogger)
	seeds.RoleSeeder(database, seedLogger)
	seeds.RolePermissionSeeder(database, seedLogger)
	seeds.UserSeeder(database, seedLogger)

}
