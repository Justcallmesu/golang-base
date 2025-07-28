package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/cookies"
	"justcallmesu.com/rest-api/internal/api/response"
	"justcallmesu.com/rest-api/internal/app/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString, tokenError := context.Cookie(os.Getenv("COOKIE_NAME"))

		if tokenError != nil {
			context.JSON(http.StatusUnauthorized, response.NewResponse("Not Authorized, Please Login", false, nil))
			context.Abort()
			return
		}

		claims, claimsError := auth.ParseToken(tokenString)

		if claimsError != nil {
			context.JSON(http.StatusUnauthorized, response.NewResponse(claimsError.Error(), false, nil))
			setCookieError := cookies.SetCookie(context, os.Getenv("COOOKIE_NAME"), "", -1)

			if setCookieError != nil {
				context.JSON(http.StatusUnauthorized, response.NewResponse(setCookieError.Error(), false, nil))
			}

			context.Abort()
			return
		}

		context.Set("UserData", claims)

		context.Next()
	}
}
