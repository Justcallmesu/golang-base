package auth

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/app/cookies"
	"justcallmesu.com/rest-api/internal/app/users"
)

type AuthService struct {
	UserRepository *users.UserRepository
	JwtService     *JWTService
	CookieService  *cookies.TokenCookieService
}

func NewAuthService(repository *users.UserRepository, jwtService *JWTService, cookieService *cookies.TokenCookieService) *AuthService {
	return &AuthService{
		UserRepository: repository,
		JwtService:     jwtService,
		CookieService:  cookieService,
	}
}

func (service *AuthService) Login(context context.Context, credentials *LoginUser) (*cookies.TokenCookie, error) {
	foundUser, loginError := service.UserRepository.FindUserByUsername(credentials.Username, context)

	if loginError != nil {
		return nil, loginError
	}

	loginError = foundUser.ComparePassword(credentials.Password)

	if loginError != nil {
		return nil, fmt.Errorf("email atau password tidak sesuai")
	}

	// Generate JWT token

	accessToken, loginError := service.JwtService.GenerateToken(foundUser.ID, foundUser.Username, AccessTokenType)
	if loginError != nil {
		return nil, loginError
	}
	refreshToken, loginError := service.JwtService.GenerateToken(foundUser.ID, foundUser.Username, RefreshTokenType)

	if loginError != nil {
		return nil, loginError
	}

	return cookies.NewTokenCookie(refreshToken, accessToken), nil
}

func (service *AuthService) RegenerateAccessToken(context *gin.Context) error {
	var claims *JWTClaims
	var claimsError error

	tokenString, tokenError := service.CookieService.GetRefreshTokenCookie(context)

	if tokenError != nil {
		return fmt.Errorf("unauthorized")
	}

	claims, claimsError = service.JwtService.ParseToken(tokenString, RefreshTokenType)

	if claimsError != nil {
		return claimsError
	}

	accessToken, accessTokenGenerateError := service.JwtService.GenerateToken(claims.UserID, claims.Username, AccessTokenType)

	if accessTokenGenerateError != nil {
		return accessTokenGenerateError
	}

	service.CookieService.GenerateAccessCookies(context, accessToken)

	regenerateRefreshTokenError := service.RegenerateRefreshToken(context)

	if regenerateRefreshTokenError != nil {
		return fmt.Errorf("unauthorized")
	}

	return nil
}

func (service *AuthService) RegenerateRefreshToken(context *gin.Context) error {
	var claims *JWTClaims

	tokenString, tokenError := service.CookieService.GetRefreshTokenCookie(context)

	if tokenError != nil {
		return fmt.Errorf("unauthorized")
	}

	claims, claimsError := service.JwtService.ParseToken(tokenString, RefreshTokenType)

	if claimsError != nil {
		return claimsError
	}

	refreshToken, refreshTokenGenerateError := service.JwtService.GenerateToken(claims.UserID, claims.Username, RefreshTokenType)

	if refreshTokenGenerateError != nil {
		return refreshTokenGenerateError
	}

	service.CookieService.GenerateRefreshCookies(context, refreshToken)

	return nil
}
