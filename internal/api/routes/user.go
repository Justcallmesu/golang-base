package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/app"
	user_handler "justcallmesu.com/golang-base/internal/app/handlers/user"
)

func UserRoutes(router *gin.RouterGroup, services *app.Services, middlewares *app.Middlewares) {

	userRouter := router.Group("/users")
	//Handler
	userHandler := user_handler.NewHandler(services.UserService)

	// Routes
	userRouter.GET("", userHandler.FindAllUsers)
	userRouter.GET(":id", userHandler.GetUserById)

	userRouter.POST("", userHandler.CreateUser)
	userRouter.PUT(":id", userHandler.UpdateUser)
	userRouter.DELETE(":id", userHandler.DeleteUser)
}
