package main

import (
	"justcallmesu.com/golang-base/internal/app/users"
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
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
