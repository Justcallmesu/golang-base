package main

import (
	"justcallmesu.com/rest-api/internal/app/users"
	"justcallmesu.com/rest-api/internal/config"
	"justcallmesu.com/rest-api/internal/database"
)

func main() {

	config.LoadConfig()

	database := database.InitConnection()

	migrateError := database.AutoMigrate(
		&users.User{},
	)

	if migrateError != nil {
		panic("Error migrating database: " + migrateError.Error())
	}

}
