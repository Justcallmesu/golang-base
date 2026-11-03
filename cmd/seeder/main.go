package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
	"justcallmesu.com/golang-base/internal/database/seeds"
)

func main() {

	config.LoadConfig()

	database := database.InitConnection()
	defer func() {

		sqlDB, sqlDbError := database.DB()

		if sqlDbError != nil {
			panic(sqlDbError)
		}

		sqlDbError = sqlDB.Close()

		if sqlDbError != nil {
			panic(sqlDbError)
		}
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
