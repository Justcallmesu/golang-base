package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/response"
	auth_service "justcallmesu.com/golang-base/internal/app/services/auth"
	token_service "justcallmesu.com/golang-base/internal/app/services/token"
	"justcallmesu.com/golang-base/internal/common/types"
)

type AuthMiddleware struct {
	AuthService *auth_service.Service
	JWTService  *token_service.Service
}

func NewAuthMiddleware(authService *auth_service.Service, jwtService *token_service.Service) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
		JWTService:  jwtService,
	}
}

func (middleware *AuthMiddleware) EnsureSessionIsValid() gin.HandlerFunc {
	return func(context *gin.Context) {
		var claims *types.JWTClaims
		var regenerateError, claimsParseError error

		accessTokenString, tokenError := context.Cookie(os.Getenv("COOKIE_ACCESS_TOKEN"))

		if tokenError != nil {
			accessTokenString, regenerateError = middleware.AuthService.RegenerateAccessToken(context)
		}

		if regenerateError != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse("unauthorized", false, nil))
			return
		}

		claims, claimsParseError = middleware.JWTService.ParseToken(accessTokenString, types.AccessTokenType)

		if claimsParseError != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse("unauthorized", false, nil))
			return
		}

		context.Set("UserData", claims)

		context.Next()
	}
}
