package main

import (
	"database/sql"
	"fmt"
	"log"

	"justcallmesu.com/rest-api/internal/config"
	"justcallmesu.com/rest-api/internal/database"
	"justcallmesu.com/rest-api/internal/migrations"
)

type MigrationFunctionType func(database *sql.DB, channel chan bool)

func main() {
	config.LoadConfig()

	database, databaseConnectionError := database.InitConnection()

	if databaseConnectionError != nil {
		log.Fatal(databaseConnectionError)
	}

	fmt.Println("RUNNING TABLE MIGRATIONS")

	var migrationFunctions []MigrationFunctionType = []MigrationFunctionType{
		migrations.CreateUsersTable,
	}

	var channels []chan bool = make([]chan bool, len(migrationFunctions))

	for index, migrationFunction := range migrationFunctions {
		channels[index] = make(chan bool)

		fmt.Println(channels[index])

		go migrationFunction(database, channels[index])
	}

	for _, channel := range channels {
		<-channel
	}
}
