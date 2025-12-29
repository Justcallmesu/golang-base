package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/routes"
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
)

func main() {
	config.LoadConfig()

	Engine := gin.Default()

	database := database.InitConnection()

	routes.SetupRoutes(Engine, database)

	engineError := Engine.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))

	if engineError != nil {
		log.Fatal(engineError)
	}
	log.Println("Server is running on port:", os.Getenv("APP_PORT"))
}
