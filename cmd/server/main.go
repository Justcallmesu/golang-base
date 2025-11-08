package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"
	"justcallmesu.com/rest-api/internal/api/routes"
	"justcallmesu.com/rest-api/internal/config"
	"justcallmesu.com/rest-api/internal/database"
	"justcallmesu.com/rest-api/internal/serializer"
)

func main() {
	config.LoadConfig()

	Engine := gin.Default()

	schema.RegisterSerializer("json", serializer.JSONSerializer{})

	database := database.InitConnection()

	routes.SetupRoutes(Engine, database)

	engineError := Engine.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))

	if engineError != nil {
		log.Fatal(engineError)
	}
	log.Println("Server is running on port:", os.Getenv("APP_PORT"))
}
