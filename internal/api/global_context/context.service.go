package global_context

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/app/tokens"
)

func GetUserData(context *gin.Context) (*tokens.JWTClaims, error) {
	claims, isExist := context.Get("UserData")

	if !isExist {
		return nil, fmt.Errorf("unauthorized, please login")
	}

	return claims.(*tokens.JWTClaims), nil
}
