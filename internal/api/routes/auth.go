package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/app"
	"justcallmesu.com/rest-api/internal/app/auth"
)

func AuthRoutes(engine *gin.Engine, services *app.Services, middlewares *app.Middlewares) {

	authRouter := engine.Group("/auth")

	//Handler
	authHandler := auth.NewAuthHandler(services.AuthService, services.CookieService)

	// Routes
	authRouter.POST("/login", authHandler.Login)
	authRouter.GET("/logout", middlewares.AuthMiddleware.EnsureSessionIsValid(), authHandler.Logout)
}
