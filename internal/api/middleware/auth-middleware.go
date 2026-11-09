package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	cookies_service "justcallmesu.com/golang-base/internal/api/cookies"
	"justcallmesu.com/golang-base/internal/app/models"
	auth_service "justcallmesu.com/golang-base/internal/app/services/auth"
	token_service "justcallmesu.com/golang-base/internal/app/services/token"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type AuthMiddleware struct {
	AuthService   *auth_service.Service
	JWTService    *token_service.Service
	CookieService *cookies_service.Service
}

func NewAuthMiddleware(authService *auth_service.Service, jwtService *token_service.Service, cookieService *cookies_service.Service) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService:   authService,
		JWTService:    jwtService,
		CookieService: cookieService,
	}
}

func (middleware *AuthMiddleware) EnsureSessionIsValid() gin.HandlerFunc {
	return func(context *gin.Context) {

		if access, err := context.Cookie(middleware.CookieService.AccessTokenName); err == nil {
			if claims, err := middleware.JWTService.ParseToken(access, types.AccessTokenType); err == nil {
				context.Set("UserData", claims)
				context.Next()
				return
			}
		}

		refreshToken, err := context.Cookie(middleware.CookieService.RefreshTokenName)

		if err != nil || refreshToken == "" {
			if err = middleware.CookieService.ResetTokenCookies(context); err != nil {
				context.Error(err)
				context.Abort()
				return
			}
			middleware.unauthorized(context, "")
			return
		}

		rotatedTokens, err := middleware.AuthService.RefreshTokens(context.Request.Context(), refreshToken)

		if err != nil {
			if apperror.IsKind(err, apperror.KindUnauthorized) {
				if err = middleware.CookieService.ResetTokenCookies(context); err != nil {
					context.Error(err)
					context.Abort()
					return
				}
				middleware.unauthorized(context, "")
				return
			}
			context.Error(err)
			context.Abort()
			return
		}

		middleware.CookieService.GenerateTokenCookies(context, rotatedTokens)
		claims, err := middleware.JWTService.ParseToken(rotatedTokens.AccessToken, types.AccessTokenType)

		if err != nil {
			if err = middleware.CookieService.ResetTokenCookies(context); err != nil {
				context.Error(err)
				context.Abort()
				return
			}
			middleware.unauthorized(context, "")
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
			middleware.unauthorized(context, "")
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

func (middleware *AuthMiddleware) unauthorized(context *gin.Context, message string) {
	context.Error(apperror.Unauthorized(utils.GetDefaultValue(message, "unauthorized access")))
	context.Abort()
}
