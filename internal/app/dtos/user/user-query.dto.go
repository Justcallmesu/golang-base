package user_dtos

import (
	"justcallmesu.com/golang-base/internal/common/types"
)

type UserQueryDto struct {
	types.PaginationQueryParams
	RoleId *types.UUID `form:"roleId"`
}
