package auth

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/cookies"
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

func (service *AuthService) SignUp(context *gin.Context) (*users.User, error) {

	var userData users.User

	createError := context.ShouldBindJSON(&userData)

	if createError != nil {
		return nil, createError
	}

	encryptError := userData.HashPassword()

	if encryptError != nil {
		return nil, encryptError
	}

	createdData, createError := service.UserRepository.Create(&userData)

	if createError != nil {
		return nil, createError
	}

	return createdData, nil
}

func (service *AuthService) Login(context *gin.Context) (string, error) {

	var credentials LoginUser

	loginError := context.ShouldBindJSON(&credentials)

	if loginError != nil {
		return "", loginError
	}

	foundUser, loginError := service.UserRepository.FindUserByEmail(credentials.Email)

	if loginError != nil {
		return "", loginError
	}

	loginError = foundUser.ComparePassword(credentials.Password)

	if loginError != nil {
		return "", fmt.Errorf("email atau password tidak sesuai")
	}

	token, loginError := GenerateToken(foundUser.ID, foundUser.Email)

	if loginError != nil {
		return "", loginError
	}

	return token, nil
}

func (service *AuthService) Logout(context *gin.Context) {
	cookies.SetCookie(context, os.Getenv("COOKIE_NAME"), "", -1)
}
