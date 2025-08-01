package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/app/auth"
)

func AuthRoutes(engine *gin.Engine, database *gorm.DB) {

	authRouter := engine.Group("/auth")

	//Handler
	authHandler := auth.NewAuthHandler(database)

	// Routes
	authRouter.POST("/login", authHandler.Login)
	authRouter.GET("/logout", authHandler.Logout)
}
