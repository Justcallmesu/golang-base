package files

import (
	"gorm.io/datatypes"
	"justcallmesu.com/rest-api/internal/types"
)

type FileType string

const (
	IMAGE    FileType = "Images"
	DOCUMENT FileType = "Document"
)

type Files struct {
	types.BaseEntityModel
	OriginalName string            `json:"originalName"`
	OriginalUrl  string            `json:"originalUrl"`
	Renditions   datatypes.JSONMap `json:"renditions"`
	FileType     FileType          `json:"fileType" gorm:"type:ENUM('Images','Document')"`
}

type FilesQuery struct {
	types.QueryParams
}
