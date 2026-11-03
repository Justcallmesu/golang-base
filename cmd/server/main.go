package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"justcallmesu.com/golang-base/internal/api/middleware"
	"justcallmesu.com/golang-base/internal/api/routes"
	"justcallmesu.com/golang-base/internal/app"
	"justcallmesu.com/golang-base/internal/common/utils"
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Get environment variables
	appPort := utils.GetEnv("APP_PORT", "3000")
	syncOnstart, parseErr := strconv.ParseBool(utils.GetEnv("DB_SYNC_ON_START", "false"))

	if parseErr != nil {
		log.Fatal("Error parsing DB_SYNC_ON_START:", parseErr)
	}

	// Initialize Gin engine and database connection
	Engine := gin.Default()
	databaseConnection := database.InitConnection()

	// Perform database migration if enabled
	if syncOnstart {
		database.Migrations(databaseConnection)
		log.Println("Database migration completed successfully.")
	}

	// Initialize Modules
	var repositories = app.NewRepositories(databaseConnection)
	var services = app.NewServices(repositories)
	var middlewares = app.NewMiddlewares(services)

	// Cors
	Engine.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length", "Content-Type"},
			AllowCredentials: true,
			AllowFiles:       true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			MaxAge:           12 * time.Hour,
		},
	))

	// Validator

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("form"), ",", 2)[0]

			if name == "" {
				name = strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			}

			if name == "-" {
				return ""
			}
			return name
		})

	}

	// Setup Routes
	Engine.Use(middleware.ErrorHandler())
	routes.SetupRoutes(Engine, repositories, services, middlewares)

	// Start the server
	engineError := Engine.Run(fmt.Sprintf(":%s", appPort))

	if engineError != nil {
		log.Fatal(engineError)
	}

	log.Println("Server is running on port:", appPort)
}
