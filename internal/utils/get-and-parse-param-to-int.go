package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/response"
)

func GetAndParseParamToInt[T int](context *gin.Context, paramName string, base int, bitSize int) (int64, *response.ErrorResponse) {
	blogId := context.Param(paramName)

	if blogId == "" {
		return 0, response.NewErrorResponse("Parameter tidak valid", errors.New("parameter tidak ditemukan"))
	}

	value, parsingError := strconv.ParseInt(blogId, base, bitSize)

	if parsingError != nil {
		return 0, response.NewErrorResponse("Parameter harus berupa angka", nil)
	}

	return value, nil
}
