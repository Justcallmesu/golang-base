package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/api"
)

func SetupRoutes(ginEngine *gin.Engine, database *gorm.DB) {

	AuthRoutes(ginEngine, database)
	ginEngine.Use(api.AuthMiddleware())
}
