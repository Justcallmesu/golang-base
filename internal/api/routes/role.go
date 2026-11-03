package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/app"
	role_handler "justcallmesu.com/golang-base/internal/app/handlers/role"
)

func RoleRoutes(router *gin.RouterGroup, services *app.Services, middlewares *app.Middlewares) {

	roleRouter := router.Group("/roles")

	//Handler
	roleHandler := role_handler.NewHandler(services.RoleService)

	// Routes
	roleRouter.GET("", roleHandler.FindAllRoles)
	roleRouter.GET(":id", roleHandler.GetRoleById)

}
