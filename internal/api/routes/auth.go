package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/handler"
)

func AuthRoutes(engine *gin.Engine, database *sql.DB) {

	authRouter := engine.Group("/auth")

	//Handler
	authHandler := handler.NewAuthHandler(database)

	// Routes
	authRouter.POST("/sign-up", authHandler.SignUp)
	authRouter.POST("/login", authHandler.Login)
	authRouter.GET("/logout", authHandler.Logout)
}
