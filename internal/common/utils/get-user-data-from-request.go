package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/common/types"
)

func GetUserData(context *gin.Context) (*types.JWTClaims, error) {
	claims, isExist := context.Get(UserDataKey)

	if !isExist || claims == nil {
		return nil, fmt.Errorf("unauthorized, please login")
	}

	parsedClaims, ok := claims.(*types.JWTClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized, please login")
	}

	return parsedClaims, nil
}
