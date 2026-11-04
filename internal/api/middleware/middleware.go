package middleware

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/response"
	"justcallmesu.com/golang-base/internal/app/models"
	auth_service "justcallmesu.com/golang-base/internal/app/services/auth"
	token_service "justcallmesu.com/golang-base/internal/app/services/token"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
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

func (middleware *AuthMiddleware) RequirePermission(code models.PermissionCodeEnum) gin.HandlerFunc {
	return func(context *gin.Context) {
		parsedClaims, err := utils.GetUserData(context)

		if err != nil {
			context.Error(apperror.Unauthorized("Unauthorized access"))
			context.Abort()
			return
		}

		hasPermission, err := middleware.AuthService.HasPermission(context.Request.Context(), parsedClaims.UserID, code)

		if err != nil {
			context.Error(err)
			context.Abort()
			return
		}

		if !hasPermission {
			slog.Warn("Permission Denied",
				"user_id", parsedClaims.UserID,
				"username", parsedClaims.Username,
				"permission_code", code,
				"request_path", context.Request.URL.Path,
				"request_method", context.Request.Method,
				"ip_address", context.ClientIP(),
			)

			context.Error(apperror.Forbidden("You do not have permission to perform this action"))
			context.Abort()

			return
		}

		context.Next()
	}
}
