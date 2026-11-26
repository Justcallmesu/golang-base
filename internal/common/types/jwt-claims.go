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
	UserID       UUID   `json:"userId"`
	RoleID       UUID   `json:"roleId"`
	Username     string `json:"username"`
	TokenVersion int    `json:"tokenVersion"`
	jwt.RegisteredClaims
}
