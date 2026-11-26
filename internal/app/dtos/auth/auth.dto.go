package auth_dtos

import "justcallmesu.com/golang-base/internal/app/models"

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=72"`
}

type SessionResponse struct {
	User           *models.User             `json:"user"`
	NavigationMenu *[]models.NavigationMenu `json:"navigationMenu"`
}
