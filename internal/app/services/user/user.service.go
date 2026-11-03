package user_service

import (
	"github.com/gin-gonic/gin"
	user_dtos "justcallmesu.com/golang-base/internal/app/dtos/user"
	"justcallmesu.com/golang-base/internal/app/models"
	user_repository "justcallmesu.com/golang-base/internal/app/repositories/user"
	"justcallmesu.com/golang-base/internal/common/types"
)

type Service struct {
	userRepository *user_repository.UserRepository
}

func NewService(userRepo *user_repository.UserRepository) *Service {
	return &Service{
		userRepository: userRepo,
	}
}
func (service *Service) GetAllUsers(context *gin.Context, queryParams *user_dtos.UserQueryDto) (*types.PaginatedResponse[models.User], error) {
	return service.userRepository.FindAll(queryParams, context.Request.Context())
}

func (service *Service) CreateUser(context *gin.Context, createUserDto *user_dtos.CreateUserDto) (*models.User, error) {
	user := &models.User{
		Username: createUserDto.Username,
		FullName: createUserDto.FullName,
		Email:    createUserDto.Email,
	}

	hashError := user.SetPassword(createUserDto.Password)

	if hashError != nil {
		return nil, hashError
	}

	createError := service.userRepository.CreateUser(user, context.Request.Context())

	if createError != nil {
		return nil, createError
	}

	return user, nil
}

func (service *Service) UpdateUser(context *gin.Context, updateUserDto *user_dtos.UpdateUserDto, userID types.UUID) (*models.User, error) {
	user, updateError := service.userRepository.FindUserByID(userID, context.Request.Context())

	if updateError != nil {
		return nil, updateError
	}

	user.Username = updateUserDto.Username
	user.RoleId = updateUserDto.RoleId
	user.FullName = updateUserDto.FullName
	user.Email = updateUserDto.Email

	if updateUserDto.Password != "" {
		updateError = user.SetPassword(updateUserDto.Password)
		if updateError != nil {
			return nil, updateError
		}
	}

	updateError = service.userRepository.UpdateUser(user, context.Request.Context())

	if updateError != nil {
		return nil, updateError
	}

	return user, nil
}

func (service *Service) DeleteUser(context *gin.Context, userID types.UUID) error {
	return service.userRepository.DeleteUser(userID, context.Request.Context())
}

func (service *Service) GetUserByID(context *gin.Context, userID types.UUID) (*models.User, error) {
	return service.userRepository.FindUserByID(userID, context.Request.Context())

}
