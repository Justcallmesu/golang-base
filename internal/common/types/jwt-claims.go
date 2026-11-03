package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessTokenType  TokenType = "Access"
	RefreshTokenType TokenType = "Refresh"
)

type JWTClaims struct {
	UserID   UUID   `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
