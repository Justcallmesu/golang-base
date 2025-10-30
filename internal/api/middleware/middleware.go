package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/response"
	"justcallmesu.com/rest-api/internal/app/auth"
	"justcallmesu.com/rest-api/internal/app/cookies"
)

type AuthMiddleware struct {
	AuthService   *auth.AuthService
	JWTService    *auth.JWTService
	cookieService *cookies.TokenCookieService
}

func NewAuthMiddleware(authService *auth.AuthService, jwtService *auth.JWTService, cookieService *cookies.TokenCookieService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService:   authService,
		JWTService:    jwtService,
		cookieService: cookieService,
	}
}

func (middleware *AuthMiddleware) EnsureSessionIsValid() gin.HandlerFunc {
	return func(context *gin.Context) {
		var claims *auth.JWTClaims
		var claimsError, regenerateError error

		tokenString, tokenError := middleware.cookieService.GetAccessTokenCookie(context)

		if tokenError != nil {
			regenerateError = middleware.AuthService.RegenerateAccessToken(context)

		} else {
			claims, claimsError = middleware.JWTService.ParseToken(tokenString, auth.AccessTokenType)

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
