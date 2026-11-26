package role_dto

import (
	"justcallmesu.com/golang-base/internal/app/models"
)

type CreateRoleDto struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

func (dto *CreateRoleDto) ToModel() *models.Role {
	return &models.Role{
		Name:        dto.Name,
		Description: dto.Description,
	}
}
