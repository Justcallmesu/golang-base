package routes

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/app"
	role_handler "justcallmesu.com/golang-base/internal/app/handlers/role"
	"justcallmesu.com/golang-base/internal/app/models"
)

func RoleRoutes(router *gin.RouterGroup, services *app.Services, middlewares *app.Middlewares) {

	roleRouter := router.Group("/roles")

	//Handler
	roleHandler := role_handler.NewHandler(services.RoleService)

	// Routes
	roleRouter.GET("", middlewares.AuthMiddleware.RequirePermission(models.PERMISSION_VIEW_ROLE), roleHandler.FindAllRoles)
	roleRouter.GET(":id", middlewares.AuthMiddleware.RequirePermission(models.PERMISSION_VIEW_ROLE), roleHandler.GetRoleById)
}
