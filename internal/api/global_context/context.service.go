package global_context

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/app/auth"
)

func GetUserData(context *gin.Context) (*auth.JWTClaims, error) {
	claims, isExist := context.Get("UserData")

	if !isExist {
		return nil, fmt.Errorf("unauthorized, please login")
	}

	return claims.(*auth.JWTClaims), nil
}
