package routes

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/app"
	"justcallmesu.com/rest-api/internal/utils"
)

func SetupRoutes(ginEngine *gin.Engine, database *gorm.DB) {

	var defaultWritePath = filepath.Clean(utils.GetDefaultValue(os.Getenv("UPLOAD_PATH"), "public/"))

	var repositories = app.NewRepositories(database)
	var services = app.NewServices(repositories, defaultWritePath)

	var middlewares = app.NewMiddlewares(services)

	ginEngine.Static("/public", filepath.Clean(defaultWritePath))

	AuthRoutes(ginEngine, services, middlewares)
	ginEngine.Use(middlewares.AuthMiddleware.EnsureSessionIsValid())
	BlogRoutes(ginEngine, services, middlewares)
}
