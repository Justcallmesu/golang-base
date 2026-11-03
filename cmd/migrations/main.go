package main

import (
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
)

func main() {

	config.LoadConfig()

	dbConnection := database.InitConnection()

	database.Migrations(dbConnection)
}
