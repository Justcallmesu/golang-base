package user_dtos

import (
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/types"
)

type CreateUserDto struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"email"`
	FullName string `json:"fullName" binding:"required"`

	Password string      `json:"password" binding:"required,min=8"`
	RoleId   *types.UUID `json:"roleId" binding:"required"`
}

func (dto *CreateUserDto) ToModel() *models.User {
	return &models.User{
		Username: dto.Username,
		Email:    dto.Email,
		FullName: dto.FullName,
		RoleId:   dto.RoleId,
	}
}
