package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJwtClaims(userId uint, username string) (*JWTClaims, error) {

	jwtExpirationDuration, jwtExpirationError := strconv.ParseInt(os.Getenv("JWT_EXPIRATION"), 10, 64)

	if jwtExpirationError != nil {
		return nil, fmt.Errorf("error parsing: Error Jwt Expiration Duration")
	}

	expiresAt := time.Now().Add(time.Duration(jwtExpirationDuration) * time.Hour)

	return &JWTClaims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "justcallmesu",
		},
	}, nil
}
