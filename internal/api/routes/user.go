package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/app"
	user_handler "justcallmesu.com/golang-base/internal/app/handlers/user"
	"justcallmesu.com/golang-base/internal/app/models"
)

func UserRoutes(router *gin.RouterGroup, services *app.Services, middlewares *app.Middlewares) {

	userRouter := router.Group("/users")
	//Handler
	userHandler := user_handler.NewHandler(services.UserService)

	// Routes
	userRouter.GET("", middlewares.AuthMiddleware.RequirePermission(models.PERMISSION_VIEW_USER), userHandler.FindAllUsers)
	userRouter.GET(":id", middlewares.AuthMiddleware.RequirePermission(models.PERMISSION_VIEW_USER), userHandler.GetUserById)

	userRouter.POST("", middlewares.AuthMiddleware.RequirePermission(models.PERMISSION_CREATE_USER), userHandler.CreateUser)
	userRouter.PUT(":id", middlewares.AuthMiddleware.RequirePermission(models.PERMISSION_UPDATE_USER), userHandler.UpdateUser)
	userRouter.DELETE(":id", middlewares.AuthMiddleware.RequirePermission(models.PERMISSION_DELETE_USER), userHandler.DeleteUser)
}
