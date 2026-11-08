package user_dtos

import (
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type UpdateUserDto struct {
	Username *string `json:"username" binding:"omitempty,min=3"`
	Email    *string `json:"email" binding:"omitempty,email"`
	FullName *string `json:"fullName" binding:"omitempty,min=1"`

	RoleId   *types.UUID `json:"roleId"`
	Password *string     `json:"password" binding:"omitempty,min=8"`
}

func (dto *UpdateUserDto) ApplyTo(userData *models.User) {

	utils.SetIfPresent(&userData.Username, dto.Username)
	utils.SetIfPresent(&userData.Email, dto.Email)
	utils.SetIfPresent(&userData.FullName, dto.FullName)

	if dto.RoleId != nil {
		userData.RoleId = dto.RoleId
	}

}
