package role_service

import (
	"context"

	role_dto "justcallmesu.com/golang-base/internal/app/dtos/role"
	"justcallmesu.com/golang-base/internal/app/models"
	role_repository "justcallmesu.com/golang-base/internal/app/repositories/role"
	"justcallmesu.com/golang-base/internal/common/types"
)

type Service struct {
	roleRepository *role_repository.Repository
}

func NewService(roleRepo *role_repository.Repository) *Service {
	return &Service{
		roleRepository: roleRepo,
	}
}

func (service *Service) GetAllRoles(context context.Context, queryParams *role_dto.RoleQueryDto) (*types.PaginatedResponse[models.Role], error) {
	return service.roleRepository.FindAll(queryParams, context)
}

func (service *Service) GetRoleByID(context context.Context, roleID types.UUID) (*models.Role, error) {
	return service.roleRepository.FindByID(roleID, context)
}

func (service *Service) CreateRole(context context.Context, roleDto *role_dto.CreateRoleDto) (*models.Role, error) {
	role := roleDto.ToModel()

	return role, service.roleRepository.Create(role, context)
}

func (service *Service) UpdateRole(context context.Context, roleDto *role_dto.UpdateRoleDto, roleID types.UUID) (*models.Role, error) {
	role, err := service.roleRepository.FindByID(roleID, context)

	if err != nil {
		return nil, err
	}

	roleDto.ApplyTo(role)
	return role, service.roleRepository.Update(role, context)
}

func (service *Service) DeleteRole(context context.Context, roleID types.UUID) error {
	return service.roleRepository.Delete(roleID, context)
}
