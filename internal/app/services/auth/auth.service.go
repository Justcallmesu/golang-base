package auth_service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	dtos "justcallmesu.com/golang-base/internal/app/dtos/auth"
	cookie_dtos "justcallmesu.com/golang-base/internal/app/dtos/cookie"
	"justcallmesu.com/golang-base/internal/app/models"
	user_repository "justcallmesu.com/golang-base/internal/app/repositories/user"
	cookies_service "justcallmesu.com/golang-base/internal/app/services/cookie"
	navigation_menu_service "justcallmesu.com/golang-base/internal/app/services/navigation-menu"
	token_service "justcallmesu.com/golang-base/internal/app/services/token"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
)

type Service struct {
	userRepository        *user_repository.UserRepository
	jwtService            *token_service.Service
	cookieService         *cookies_service.Service
	navigationMenuService *navigation_menu_service.Service
}

func NewService(repository *user_repository.UserRepository, jwtService *token_service.Service, cookieService *cookies_service.Service, navigationMenuService *navigation_menu_service.Service) *Service {
	return &Service{
		userRepository:        repository,
		jwtService:            jwtService,
		cookieService:         cookieService,
		navigationMenuService: navigationMenuService,
	}
}

func (service *Service) Login(context context.Context, credentials *dtos.LoginUser) (*cookie_dtos.TokenCookie, *models.User, *[]models.NavigationMenu, error) {
	foundUser, loginError := service.userRepository.FindUserByUsername(credentials.Username, context)

	if loginError != nil {
		return nil, nil, nil, fmt.Errorf("user or password doesnt match")
	}

	loginError = foundUser.ComparePassword(credentials.Password)

	if loginError != nil {
		return nil, nil, nil, fmt.Errorf("user or password doesnt match")
	}

	// Generate JWT token

	accessToken, loginError := service.jwtService.GenerateToken(foundUser.ID, foundUser.Username, types.AccessTokenType)
	if loginError != nil {
		return nil, nil, nil, loginError
	}
	refreshToken, loginError := service.jwtService.GenerateToken(foundUser.ID, foundUser.Username, types.RefreshTokenType)

	if loginError != nil {
		return nil, nil, nil, loginError
	}

	navigationMenu := service.navigationMenuService.GetNavigationMenuBasedOnUserPermission(foundUser)

	return cookies_service.NewTokenCookie(refreshToken, accessToken), foundUser, navigationMenu, nil
}

func (service *Service) ChangePassword(context context.Context, passwords *dtos.ChangePassword, userData *types.JWTClaims) error {

	foundUser, findUserError := service.userRepository.FindUserByID(userData.UserID, context)

	if findUserError != nil {
		return findUserError
	}

	compareError := foundUser.ComparePassword(passwords.OldPassword)

	if compareError != nil {
		return fmt.Errorf("old password is incorrect")
	}

	setError := foundUser.SetPassword(passwords.NewPassword)

	if setError != nil {
		return setError
	}

	updateError := service.userRepository.UpdateUserPassword(foundUser, context)

	if updateError != nil {
		return updateError
	}

	return nil
}

func (service *Service) GetMe(context context.Context, userData *types.JWTClaims) (*models.User, *[]models.NavigationMenu, error) {

	user, findUserError := service.userRepository.FindUserByID(userData.UserID, context)

	if findUserError != nil {
		return nil, nil, findUserError
	}

	navigationMenu := service.navigationMenuService.GetNavigationMenuBasedOnUserPermission(user)

	return user, navigationMenu, nil
}

func (service *Service) IsRefreshTokenExist(context *gin.Context) (bool, string) {
	tokenString, tokenError := context.Cookie(os.Getenv("COOKIE_REFRESH_TOKEN"))

	return tokenError == nil, tokenString
}

func (service *Service) RegenerateAccessToken(context *gin.Context) (string, error) {
	var claims *types.JWTClaims
	var claimsError error

	isRefreshExist, tokenString := service.IsRefreshTokenExist(context)

	if !isRefreshExist {
		return "", fmt.Errorf("refresh token not found")
	}

	claims, claimsError = service.jwtService.ParseToken(tokenString, types.RefreshTokenType)

	if claimsError != nil {
		return "", claimsError
	}

	accessToken, accessTokenGenerateError := service.jwtService.GenerateToken(claims.UserID, claims.Username, types.AccessTokenType)
	refrehToken, refreshTokenGenerateError := service.jwtService.GenerateToken(claims.UserID, claims.Username, types.RefreshTokenType)

	if refreshTokenGenerateError != nil {
		return "", refreshTokenGenerateError
	}

	if accessTokenGenerateError != nil {
		return "", accessTokenGenerateError
	}

	service.cookieService.GenerateAccessCookies(context, accessToken)
	service.cookieService.GenerateRefreshCookies(context, refrehToken)

	return accessToken, nil

}

func (s *Service) HasPermission(ctx context.Context, userID types.UUID, code models.PermissionCodeEnum) (bool, error) {
	foundUser, err := s.userRepository.FindUserByID(userID, ctx)

	if err != nil {
		var appError *apperror.Error
		if errors.As(err, &appError) && appError.Kind == apperror.KindNotFound {
			return false, apperror.Unauthorized("User no longer exists")
		}
		return false, err
	}

	return foundUser.Role.HasPermission(code), nil
}
