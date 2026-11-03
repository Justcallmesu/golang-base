package user_dtos

import (
	"justcallmesu.com/golang-base/internal/common/types"
)

type CreateUserDto struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"email"`
	FullName string `json:"fullName" binding:"required"`

	Password string      `json:"password" binding:"required"`
	RoleId   *types.UUID `json:"roleId" binding:"required"`
}

type UpdateUserDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName" binding:"required"`

	RoleId   *types.UUID `json:"roleId"`
	Password string      `json:"password"`
}

type UserQueryDto struct {
	types.PaginationQueryParams
	RoleId *types.UUID `form:"roleId"`
}
