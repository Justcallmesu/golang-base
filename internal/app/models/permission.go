package models

import "justcallmesu.com/golang-base/internal/common/types"

type Permission struct {
	types.BaseModel
	PermissionCode PermissionCodeEnum `gorm:"unique;not null;index" json:"permission_code"`
	Group          string             `json:"group,omitempty"`
	Description    string             `json:"description,omitempty"`
}

type PermissionCodeEnum string

const (
	PERMISSION_CREATE_USER PermissionCodeEnum = "CREATE_USER"
	PERMISSION_UPDATE_USER PermissionCodeEnum = "EDIT_USER"
	PERMISSION_DELETE_USER PermissionCodeEnum = "DELETE_USER"
	PERMISSION_VIEW_USER   PermissionCodeEnum = "VIEW_USER"

	PERMISSION_CREATE_ROLE PermissionCodeEnum = "CREATE_ROLE"
	PERMISSION_UPDATE_ROLE PermissionCodeEnum = "UPDATE_ROLE"
	PERMISSION_DELETE_ROLE PermissionCodeEnum = "DELETE_ROLE"
	PERMISSION_VIEW_ROLE   PermissionCodeEnum = "VIEW_ROLE"
)
