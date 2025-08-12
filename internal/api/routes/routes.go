package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/app"
)

func SetupRoutes(ginEngine *gin.Engine, database *gorm.DB) {

	var repositories = app.NewRepositories(database)
	var services = app.NewServices(repositories)

	var middlewares = app.NewMiddlewares(services)

	AuthRoutes(ginEngine, services,middlewares)
	ginEngine.Use(middlewares.AuthMiddleware.EnsureSessionIsValid())
}
