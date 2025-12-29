package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/response"
	"justcallmesu.com/golang-base/internal/app/auth"
	"justcallmesu.com/golang-base/internal/app/tokens"
)

type AuthMiddleware struct {
	AuthService *auth.Service
	JWTService  *tokens.JWTService
}

func NewAuthMiddleware(authService *auth.Service, jwtService *tokens.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
		JWTService:  jwtService,
	}
}

func (middleware *AuthMiddleware) EnsureSessionIsValid() gin.HandlerFunc {
	return func(context *gin.Context) {
		var claims *tokens.JWTClaims
		var claimsError, regenerateError error

		tokenString, tokenError := context.Cookie(os.Getenv("COOKIE_ACCESS_TOKEN"))

		if tokenError != nil {
			regenerateError = middleware.AuthService.RegenerateAccessToken(context)

		} else {
			claims, claimsError = middleware.JWTService.ParseToken(tokenString, tokens.AccessTokenType)
			fmt.Println(tokenString)

			if claimsError != nil {
				regenerateError = middleware.AuthService.RegenerateAccessToken(context)
			}
		}

		if regenerateError != nil || claimsError != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse("unauthorized", false, nil))
			return
		}

		context.Set("UserData", claims)

		context.Next()
	}
}
