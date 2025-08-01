package auth

import (
	"context"
	"fmt"

	"justcallmesu.com/rest-api/internal/app/users"
)

type AuthService struct {
	UserRepository *users.UserRepository
}

func NewAuthService(repository *users.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: repository,
	}
}

func (service *AuthService) Login(context context.Context, credentials *LoginUser) (string, error) {
	foundUser, loginError := service.UserRepository.FindUserByUsername(credentials.Username, context)

	if loginError != nil {
		return "", loginError
	}

	loginError = foundUser.ComparePassword(credentials.Password)

	if loginError != nil {
		return "", fmt.Errorf("email atau password tidak sesuai")
	}

	token, loginError := GenerateToken(foundUser.ID, foundUser.Username)

	if loginError != nil {
		return "", loginError
	}

	return token, nil
}
