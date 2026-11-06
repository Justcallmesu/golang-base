package role_handler

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/request"
	role_dto "justcallmesu.com/golang-base/internal/app/dtos/role"
	role_service "justcallmesu.com/golang-base/internal/app/services/role"
	"justcallmesu.com/golang-base/internal/common/types"
)

type Handler struct {
	roleService *role_service.Service
}

func NewHandler(roleService *role_service.Service) *Handler {
	return &Handler{
		roleService: roleService,
	}
}

func (handler *Handler) FindAllRoles(context *gin.Context) {
	var queryParams role_dto.RoleQueryDto = role_dto.RoleQueryDto{
		PaginationQueryParams: types.NewPaginationQueryParams(),
	}

	if err := context.ShouldBindQuery(&queryParams); err != nil {
		context.Error(request.BindError(err))
		return
	}

	paginatedResponse, err := handler.roleService.GetAllRoles(context.Request.Context(), &queryParams)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(200, paginatedResponse)
}

func (handler *Handler) GetRoleById(context *gin.Context) {
	roleId, parseError := request.ParamUUID(context, "id")

	if parseError != nil {
		context.Error(parseError)
		return
	}

	role, err := handler.roleService.GetRoleByID(context.Request.Context(), roleId)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(200, role)
}
