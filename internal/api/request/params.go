package request

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
)

func ParamUUID(context *gin.Context, name string) (types.UUID, error) {
	id, err := types.ParseUUID(context.Param(name))
	if err != nil || id.UUID == uuid.Nil {
		return types.UUID{}, apperror.BadRequest(fmt.Sprintf("%s must be a valid UUID", name))
	}
	return id, nil
}
