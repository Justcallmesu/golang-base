package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api"
)

func SetupRoutes(ginEngine *gin.Engine, database *sql.DB) {

	AuthRoutes(ginEngine, database)
	ginEngine.Use(api.AuthMiddleware())
}
