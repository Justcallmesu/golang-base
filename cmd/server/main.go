package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
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
	configError := config.LoadConfig()

	if configError != nil {
		slog.Error("Error loading configuration", "error", configError)
		os.Exit(1)
	}

	// Get environment variables
	appPort := utils.GetEnv("APP_PORT", "3000")
	syncOnstart, parseErr := strconv.ParseBool(utils.GetEnv("DB_SYNC_ON_START", "false"))

	if parseErr != nil {
		slog.Error("Error parsing DB_SYNC_ON_START", "error", parseErr)
		os.Exit(1)
	}

	// Initialize Gin engine and database connection
	Engine := gin.Default()
	databaseConnection, databaseError := database.InitConnection()

	if databaseError != nil {
		slog.Error("Error initializing database connection", "error", databaseError)
		os.Exit(1)
	}

	// Perform database migration if enabled
	if syncOnstart {
		migrationError := database.Migrations(databaseConnection)
		if migrationError != nil {
			slog.Error("Database migration failed", "error", migrationError)
			os.Exit(1)
		}
		slog.Info("Database migration completed successfully.")
	}

	// Initialize Modules
	var repositories = app.NewRepositories(databaseConnection)
	var services, serviceError = app.NewServices(repositories)

	if serviceError != nil {
		slog.Error("Invalid Configuration", "error", serviceError)
		os.Exit(1)
	}

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

	Engine.Use(middleware.ErrorHandler())
	routes.SetupRoutes(Engine, repositories, services, middlewares)

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
	serverInstance := &http.Server{
		Addr:              fmt.Sprintf(":%s", appPort),
		Handler:           Engine, // Gin Engine As Handler
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stopSignal()

	go func() {
		slog.Info("Server listening", "addr", serverInstance.Addr)

		if err := serverInstance.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	slog.Info("Server shutting down...")

	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), 8*time.Second)

	defer shutdownCancel()

	if err := serverInstance.Shutdown(shutdownContext); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	if err := database.CloseConnection(databaseConnection); err != nil {
		slog.Error("Error closing database connection", "error", err)
	}

	slog.Info("Server exited gracefully")

}
