package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/routes"
	"justcallmesu.com/rest-api/internal/config"
	"justcallmesu.com/rest-api/internal/database"
)

func main() {
	config.LoadConfig()

	Engine := gin.Default()

	database, databaseConnectionError := database.InitConnection()

	if databaseConnectionError != nil {
		log.Fatal(databaseConnectionError)
	}

	routes.SetupRoutes(Engine, database)

	engineError := Engine.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))

	if engineError != nil {
		log.Fatal(engineError)
	}
	log.Println("Server is running on port:", os.Getenv("APP_PORT"))
}
