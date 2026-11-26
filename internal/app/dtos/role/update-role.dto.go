package role_dto

import (
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type UpdateRoleDto struct {
	CreateRoleDto
}

func (dto *UpdateRoleDto) ApplyTo(roleData *models.Role) {
	utils.SetIfPresent(&roleData.Name, &dto.Name)
	utils.SetIfPresent(&roleData.Description, &dto.Description)
}
