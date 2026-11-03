package role_dto

import "justcallmesu.com/golang-base/internal/common/types"

type CreateRoleDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateRoleDto struct {
	CreateRoleDto
}

type RoleQueryDto struct {
	types.PaginationQueryParams
}
