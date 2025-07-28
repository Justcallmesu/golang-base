package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId int64, email string) (string, error) {

	claims, claimsError := NewJwtClaims(userId, email)

	if claimsError != nil {
		return "", claimsError
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, signError := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if signError != nil {
		return "", signError
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, parseError := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
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
