package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/response"
	"justcallmesu.com/rest-api/internal/app/auth"
)

type AuthMiddleware struct {
	AuthService *auth.AuthService
	JWTService  *auth.JWTService
}

func NewAuthMiddleware(authService *auth.AuthService, jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
		JWTService:  jwtService,
	}
}

func (middleware *AuthMiddleware) EnsureSessionIsValid() gin.HandlerFunc {
	return func(context *gin.Context) {
		var claims *auth.JWTClaims
		var claimsError, regenerateError error

		tokenString, tokenError := context.Cookie(os.Getenv("COOKIE_ACCESS_TOKEN"))

		if tokenError != nil {
			regenerateError = middleware.AuthService.RegenerateAccessToken(context)

		} else {
			claims, claimsError = middleware.JWTService.ParseToken(tokenString, auth.AccessTokenType)
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
