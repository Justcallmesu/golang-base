package role_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/request"
	"justcallmesu.com/golang-base/internal/api/response"
	role_dto "justcallmesu.com/golang-base/internal/app/dtos/role"
	role_service "justcallmesu.com/golang-base/internal/app/services/role"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
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
	var queryParams = role_dto.RoleQueryDto{
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

	context.JSON(http.StatusOK, response.NewResponse("Roles retrieved successfully", true, paginatedResponse))
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

	context.JSON(http.StatusOK, response.NewResponse("Role retrieved successfully", true, role))
}

func (handler *Handler) CreateRole(context *gin.Context) {
	var createRoleDto role_dto.CreateRoleDto

	if err := context.ShouldBindJSON(&createRoleDto); err != nil {
		context.Error(request.BindError(err))
		return
	}

	role, err := handler.roleService.CreateRole(context.Request.Context(), &createRoleDto)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusCreated, response.NewResponse("Role created successfully", true, role))
}

func (handler *Handler) UpdateRole(context *gin.Context) {
	roleId, parseError := request.ParamUUID(context, "id")

	if parseError != nil {
		context.Error(parseError)
		return
	}

	var updateRoleDto role_dto.UpdateRoleDto

	if err := context.ShouldBindJSON(&updateRoleDto); err != nil {
		context.Error(request.BindError(err))
		return
	}

	role, err := handler.roleService.UpdateRole(context.Request.Context(), &updateRoleDto, roleId)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, response.NewResponse("Role updated successfully", true, role))
}

func (handler *Handler) DeleteRole(context *gin.Context) {
	roleId, parseError := request.ParamUUID(context, "id")

	if parseError != nil {
		context.Error(parseError)
		return
	}

	userData, err := utils.GetUserData(context)

	if err != nil {
		context.Error(err)
		return
	}

	if userData.RoleID == roleId {
		context.Error(apperror.BadRequest("you cannot delete your own role"))
		return
	}

	err = handler.roleService.DeleteRole(context.Request.Context(), roleId)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, response.NewResponse("Role deleted successfully", true, nil))
}
