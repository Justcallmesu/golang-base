package models

import (
	"justcallmesu.com/golang-base/internal/common/types"
)

type Role struct {
	types.BaseSoftDeleteModel
	Name        string       `gorm:"unique;not null;index" json:"name"`
	Description string       `json:"description,omitempty" `
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}

func (role *Role) HasPermission(permissionCode PermissionCodeEnum) bool {
	for _, permission := range role.Permissions {
		if permission.PermissionCode == permissionCode {
			return true
		}
	}
	return false
}
