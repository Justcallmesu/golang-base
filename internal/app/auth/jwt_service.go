package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType int

const (
	AccessTokenType = iota
	RefreshTokenType
)

type JWTService struct {
	accessTokenSecret  string
	refreshTokenSecret string

	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

func NewJWTService() *JWTService {
	accessTokenExpiration, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_EXPIRATION"), 10, 64)
	if err != nil {
		panic("Invalid JWT_ACCESS_EXPIRATION value")
	}

	refreshTokenExpiration, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_EXPIRATION"), 10, 64)
	if err != nil {
		panic("Invalid JWT_REFRESH_EXPIRATION value")
	}

	accessTokenExpirationDate := time.Duration(accessTokenExpiration) * time.Minute
	refreshTokenExpirationDate := time.Duration(refreshTokenExpiration) * time.Hour

	return &JWTService{
		accessTokenSecret:      os.Getenv("JWT_ACCESS_SECRET"),
		refreshTokenSecret:     os.Getenv("JWT_REFRESH_SECRET"),
		accessTokenExpiration:  accessTokenExpirationDate,
		refreshTokenExpiration: refreshTokenExpirationDate,
	}
}

func (service *JWTService) GenerateToken(userId uint, username string, tokenType TokenType) (string, error) {
	var generateError error
	var token *jwt.Token
	var signedToken, secret string
	var expiresAt *jwt.NumericDate


	switch tokenType {
	case AccessTokenType:
		expiresAt = jwt.NewNumericDate(time.Now().Add(service.accessTokenExpiration))
		secret = service.accessTokenSecret
	case RefreshTokenType:
		expiresAt = jwt.NewNumericDate(time.Now().Add(service.refreshTokenExpiration))
		secret = service.refreshTokenSecret
	default:
		return "", errors.New("invalid token type")
	}

	claims := &JWTClaims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "justcallmesu",
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, generateError = token.SignedString([]byte(secret))

	if generateError != nil {
		return "", fmt.Errorf("failed to sign access token: %w", generateError)
	}
	return signedToken, nil
}

func (service *JWTService) ParseToken(tokenString string, tokenType TokenType) (*JWTClaims, error) {
	claims := &JWTClaims{}

	var token *jwt.Token
	var parseError error

	token, parseError = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		switch tokenType {
		case AccessTokenType:
			return []byte(service.accessTokenSecret), nil
		case RefreshTokenType:
			return []byte(service.refreshTokenSecret), nil
		}
		return nil, fmt.Errorf("invalid token type: %d", tokenType)
	})

	if parseError != nil {
		if errors.Is(parseError, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired, please relogin")
		}
		return nil, parseError
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid, please relogin")
	}

	parsedClaims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return parsedClaims, nil
}
