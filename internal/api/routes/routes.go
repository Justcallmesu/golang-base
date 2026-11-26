package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/response"
	"justcallmesu.com/golang-base/internal/app"
	"justcallmesu.com/golang-base/internal/common/apperror"
)

func SetupRoutes(ginEngine *gin.Engine,
	repositories *app.Repositories,
	services *app.Services,
	middlewares *app.Middlewares,
) {
	// Index Route
	ginEngine.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, response.NewResponse("🚀 Server up and running, for server health checking use /health", true, nil))
	})

	apiRoutes := ginEngine.Group("/api")
	v1Routes := apiRoutes.Group("/v1")

	AuthRoutes(v1Routes, services, middlewares)
	v1Routes.Use(middlewares.AuthMiddleware.EnsureSessionIsValid())
	UserRoutes(v1Routes, services, middlewares)
	RoleRoutes(v1Routes, services, middlewares)

	ginEngine.NoRoute(func(context *gin.Context) {
		context.Error(apperror.NotFound("Route not found"))
	})
}
