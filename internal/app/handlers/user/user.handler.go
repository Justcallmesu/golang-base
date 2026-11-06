package user_handler

import (
	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/request"
	"justcallmesu.com/golang-base/internal/api/response"
	user_dtos "justcallmesu.com/golang-base/internal/app/dtos/user"
	user_service "justcallmesu.com/golang-base/internal/app/services/user"
	"justcallmesu.com/golang-base/internal/common/types"
)

type Handler struct {
	userService *user_service.Service
}

func NewHandler(userService *user_service.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (handler *Handler) FindAllUsers(context *gin.Context) {
	var queryParams user_dtos.UserQueryDto = user_dtos.UserQueryDto{
		PaginationQueryParams: types.NewPaginationQueryParams(),
	}

	if err := context.ShouldBindQuery(&queryParams); err != nil {
		context.Error(request.BindError(err))
		return
	}

	paginatedResponse, err := handler.userService.GetAllUsers(context.Request.Context(), &queryParams)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(200, paginatedResponse)
}

func (handler *Handler) GetUserById(context *gin.Context) {
	userId, parseError := request.ParamUUID(context, "id")

	if parseError != nil {
		context.Error(parseError)
		return
	}

	user, err := handler.userService.GetUserByID(context.Request.Context(), userId)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(200, response.NewResponse("User retrieved successfully", true, user))
}

func (handler *Handler) CreateUser(context *gin.Context) {
	var createUserDto user_dtos.CreateUserDto

	if err := context.ShouldBindJSON(&createUserDto); err != nil {
		context.Error(request.BindError(err))
		return
	}

	user, err := handler.userService.CreateUser(context.Request.Context(), &createUserDto)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(201, response.NewResponse("User created successfully", true, user))
}

func (handler *Handler) UpdateUser(context *gin.Context) {
	userId, parseError := request.ParamUUID(context, "id")

	if parseError != nil {
		context.Error(parseError)
		return
	}
	var updateUserDto user_dtos.UpdateUserDto

	if err := context.ShouldBindJSON(&updateUserDto); err != nil {
		context.Error(request.BindError(err))
		return
	}

	user, err := handler.userService.UpdateUser(context.Request.Context(), &updateUserDto, userId)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(200, response.NewResponse("User updated successfully", true, user))

}

func (handler *Handler) DeleteUser(context *gin.Context) {
	userId, err := request.ParamUUID(context, "id")

	if err != nil {
		context.Error(err)
		return
	}

	err = handler.userService.DeleteUser(context.Request.Context(), userId)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(200, response.NewResponse("User deleted successfully", true, nil))
}
