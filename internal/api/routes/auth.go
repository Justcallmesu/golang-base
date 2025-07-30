package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/api/handler"
)

func AuthRoutes(engine *gin.Engine, database *gorm.DB) {

	authRouter := engine.Group("/auth")

	//Handler
	authHandler := handler.NewAuthHandler(database)

	// Routes
	authRouter.POST("/sign-up", authHandler.SignUp)
	authRouter.POST("/login", authHandler.Login)
	authRouter.GET("/logout", authHandler.Logout)
}
