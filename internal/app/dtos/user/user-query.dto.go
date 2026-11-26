package user_dtos

import (
	"justcallmesu.com/golang-base/internal/common/types"
)

type UserQueryDto struct {
	types.PaginationQueryParams
	RoleId string `form:"roleId" binding:"omitempty,uuid"`
}
