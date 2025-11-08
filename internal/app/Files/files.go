package files

import (
	"justcallmesu.com/rest-api/internal/types"
)

type FileType string

const (
	IMAGE    FileType = "Images"
	DOCUMENT FileType = "Document"
)

type Renditions struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Small     string `json:"small,omitempty"`
	Medium    string `json:"medium,omitempty"`
	Big       string `json:"big,omitempty"`
}

type Files struct {
	types.BaseEntityModel
	OriginalName string     `json:"originalName,omitempty"`
	OriginalUrl  string     `json:"originalUrl,omitempty"`
	Renditions   Renditions `json:"renditions,omitempty" gorm:"type:json;serializer:json"`
	FileType     FileType   `json:"fileType,omitempty" gorm:"type:ENUM('Images','Document')"`
}

type FilesQuery struct {
	types.QueryParams
}
