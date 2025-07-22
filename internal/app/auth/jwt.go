package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJwtClaims(userId int64, email string) (*JWTClaims, error) {

	jwtExpirationDuration, jwtExpirationError := strconv.ParseInt(os.Getenv("JWT_EXPIRATION"), 10, 64)

	if jwtExpirationError != nil {
		return nil, fmt.Errorf("error parsing: Error Jwt Expiration Duration")
	}

	expiresAt := time.Now().Add(time.Duration(jwtExpirationDuration) * time.Hour)

	return &JWTClaims{
		UserID: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "justcallmesu",
		},
	}, nil
}
