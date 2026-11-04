package token_service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type Service struct {
	accessTokenSecret  string
	refreshTokenSecret string

	/*
		Jwt Configuration
	*/
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
	issuer                 string
}

func NewService() *Service {
	accessTokenExpiration, err := strconv.ParseInt(utils.GetEnv("JWT_ACCESS_EXPIRATION", "15"), 10, 64)
	if err != nil {
		panic("Invalid JWT_ACCESS_EXPIRATION value")
	}

	refreshTokenExpiration, err := strconv.ParseInt(utils.GetEnv("JWT_REFRESH_EXPIRATION", "24"), 10, 64)
	if err != nil {
		panic("Invalid JWT_REFRESH_EXPIRATION value")
	}

	issuer := utils.GetEnv("APP_NAME", "MyApp")

	accessTokenExpirationDate := time.Duration(accessTokenExpiration) * time.Minute
	refreshTokenExpirationDate := time.Duration(refreshTokenExpiration) * time.Hour

	return &Service{
		accessTokenSecret:      os.Getenv("JWT_ACCESS_SECRET"),
		refreshTokenSecret:     os.Getenv("JWT_REFRESH_SECRET"),
		accessTokenExpiration:  accessTokenExpirationDate,
		refreshTokenExpiration: refreshTokenExpirationDate,
		issuer:                 issuer,
	}
}

func (service *Service) GenerateToken(userId types.UUID, username string, tokenType types.TokenType) (string, error) {
	var generateError error
	var token *jwt.Token
	var signedToken, secret string
	var expiresAt *jwt.NumericDate

	switch tokenType {
	case types.AccessTokenType:
		expiresAt = jwt.NewNumericDate(time.Now().Add(service.accessTokenExpiration))
		secret = service.accessTokenSecret
	case types.RefreshTokenType:
		expiresAt = jwt.NewNumericDate(time.Now().Add(service.refreshTokenExpiration))
		secret = service.refreshTokenSecret
	default:
		return "", apperror.Unauthorized(fmt.Sprintf("invalid token type: %s", tokenType))
	}

	claims := service.NewJwtClaims(userId, username, service.issuer, expiresAt.Time)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, generateError = token.SignedString([]byte(secret))

	if generateError != nil {
		return "", apperror.Unauthorized("failed to sign the token")
	}
	return signedToken, nil
}

func (service *Service) ParseToken(tokenString string, tokenType types.TokenType) (*types.JWTClaims, error) {
	claims := &types.JWTClaims{}

	var token *jwt.Token
	var parseError error

	token, parseError = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.Unauthorized(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		switch tokenType {
		case types.AccessTokenType:
			return []byte(service.accessTokenSecret), nil
		case types.RefreshTokenType:
			return []byte(service.refreshTokenSecret), nil
		}
		return nil, apperror.Unauthorized(fmt.Sprintf("invalid token type: %s", tokenType))
	})

	if parseError != nil {
		if errors.Is(parseError, jwt.ErrTokenExpired) {
			return nil, apperror.Unauthorized("token has expired, please relogin")
		}
		return nil, parseError
	}

	if !token.Valid {
		return nil, apperror.Unauthorized("invalid token")
	}

	parsedClaims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, apperror.Unauthorized("invalid token")
	}

	return parsedClaims, nil
}

func (service *Service) NewJwtClaims(userId types.UUID, username, issuer string, expiresAt time.Time) *types.JWTClaims {
	return &types.JWTClaims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}
}
