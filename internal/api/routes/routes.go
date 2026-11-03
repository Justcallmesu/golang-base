package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/app"
)

func SetupRoutes(ginEngine *gin.Engine,
	repositories *app.Repositories,
	services *app.Services,
	middlewares *app.Middlewares,
) {
	var startTime = time.Now()

	// Index Route
	ginEngine.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Server Is Up And Running",
			"status":  "success",
			"data": gin.H{
				"server-runtime": time.Since(startTime).Truncate(time.Second).String(),
			},
		})
	})

	apiRoutes := ginEngine.Group("/api")
	v1Routes := apiRoutes.Group("/v1")

	AuthRoutes(v1Routes, services, middlewares)
	v1Routes.Use(middlewares.AuthMiddleware.EnsureSessionIsValid())
	UserRoutes(v1Routes, services, middlewares)
	RoleRoutes(v1Routes, services, middlewares)
	ginEngine.NoRoute(func(context *gin.Context) {
		context.JSON(404, gin.H{
			"message": "Route Not Found",
			"status":  "error",
		})
	})
}
