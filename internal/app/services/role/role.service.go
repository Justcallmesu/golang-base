package role_service

import (
	"github.com/gin-gonic/gin"
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

func (service *Service) GetAllRoles(context *gin.Context, queryParams *role_dto.RoleQueryDto) (*types.PaginatedResponse[models.Role], error) {
	return service.roleRepository.FindAll(queryParams, context.Request.Context())
}

func (service *Service) GetRoleByID(context *gin.Context, roleID types.UUID) (*models.Role, error) {
	return service.roleRepository.FindByID(roleID, context.Request.Context())
}
