package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/cookies"
	"justcallmesu.com/rest-api/internal/app/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString, tokenError := context.Cookie(os.Getenv("COOKIE_NAME"))

		if tokenError != nil {
			context.JSON(http.StatusUnauthorized, NewResponse("Not Authorized, Please Login", false, nil))
			context.Abort()
			return
		}

		claims, claimsError := auth.ParseToken(tokenString)

		if claimsError != nil {
			context.JSON(http.StatusUnauthorized, NewResponse(claimsError.Error(), false, nil))
			cookies.SetCookie(context, os.Getenv("COOOKIE_NAME"), "", -1)
			context.Abort()
			return
		}

		context.Set("UserData", claims)

		context.Next()
	}
}
