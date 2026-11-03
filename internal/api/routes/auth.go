package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/app"
	auth_handler "justcallmesu.com/golang-base/internal/app/handlers/auth"
)

func AuthRoutes(router *gin.RouterGroup, services *app.Services, middlewares *app.Middlewares) {

	authRouter := router.Group("/auth")

	//Handler
	authHandler := auth_handler.NewAuthHandler(services.AuthService, services.CookieService)

	// Routes
	authRouter.POST("/login", authHandler.Login)

	authRouter.Use(middlewares.AuthMiddleware.EnsureSessionIsValid())
	authRouter.GET("/me", authHandler.Me)
	authRouter.PATCH("/change-password", authHandler.ChangePassword)
	authRouter.POST("/logout", authHandler.Logout)
}
