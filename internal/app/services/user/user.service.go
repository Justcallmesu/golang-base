package user_service

import (
	"context"

	user_dtos "justcallmesu.com/golang-base/internal/app/dtos/user"
	"justcallmesu.com/golang-base/internal/app/models"
	role_repository "justcallmesu.com/golang-base/internal/app/repositories/role"
	user_repository "justcallmesu.com/golang-base/internal/app/repositories/user"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
)

type Service struct {
	userRepository *user_repository.Repository
	roleRepository *role_repository.Repository
}

func NewService(userRepo *user_repository.Repository, roleRepo *role_repository.Repository) *Service {
	return &Service{
		userRepository: userRepo,
		roleRepository: roleRepo,
	}
}
func (service *Service) GetAllUsers(context context.Context, queryParams *user_dtos.UserQueryDto) (*types.PaginatedResponse[models.User], error) {
	return service.userRepository.FindAll(queryParams, context)
}

func (service *Service) CreateUser(context context.Context, createUserDto *user_dtos.CreateUserDto) (*models.User, error) {
	user := &models.User{
		Username: createUserDto.Username,
		FullName: createUserDto.FullName,
		Email:    createUserDto.Email,
		RoleId:   createUserDto.RoleId,
	}

	if err := service.validate(user, context); err != nil {
		return nil, err
	}

	hashError := user.SetPassword(createUserDto.Password)

	if hashError != nil {
		return nil, hashError
	}

	createError := service.userRepository.CreateUser(user, context)

	if createError != nil {
		return nil, createError
	}

	return user, nil
}

func (service *Service) UpdateUser(context context.Context, updateUserDto *user_dtos.UpdateUserDto, userID types.UUID) (*models.User, error) {
	user, updateError := service.userRepository.FindUserByID(userID, context)

	if updateError != nil {
		return nil, updateError
	}

	user.ApplyTo(updateUserDto)

	if updateUserDto.RoleId != nil {
		user.RoleId = updateUserDto.RoleId
		if err := service.validate(user, context); err != nil {
			return nil, err
		}
	}

	if updateUserDto.Password != "" {
		updateError = user.SetPassword(updateUserDto.Password)
		if updateError != nil {
			return nil, updateError
		}
	}

	updateError = service.userRepository.UpdateUser(user, context)

	if updateError != nil {
		return nil, updateError
	}

	return user, nil
}

func (service *Service) DeleteUser(context context.Context, userID types.UUID) error {
	if _, err := service.userRepository.FindUserByID(userID, context); err != nil {
		return err
	}

	return service.userRepository.DeleteUser(userID, context)
}

func (service *Service) GetUserByID(context context.Context, userID types.UUID) (*models.User, error) {
	return service.userRepository.FindUserByID(userID, context)

}

func (service *Service) validate(user *models.User, context context.Context) error {
	if user.RoleId == nil {
		return apperror.BadRequest("role is required")
	}

	// Check if the role exists
	if _, err := service.roleRepository.FindByID(*user.RoleId, context); err != nil {
		if apperror.IsKind(err, apperror.KindNotFound) {
			return apperror.BadRequest("role does not exist")
		}
		return err
	}
	return nil
}
