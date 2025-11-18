package blogs

import (
	files "justcallmesu.com/rest-api/internal/app/Files"
	"justcallmesu.com/rest-api/internal/types"
)

type BlogDetailsType string

const (
	Image = "Image"
	Code  = "Code"
	Text  = "Text"
)

type BlogDetails struct {
	types.BaseEntityModel
	BlogId  uint            `json:"blogId"`
	Blog    Blog            `json:"blog,omitzero"`
	Content string          `json:"content"`
	FileId  *uint           `json:"fileId,omitempty"`
	File    files.Files     `json:"file,omitzero"`
	Type    BlogDetailsType `json:"type"`
	Order   uint            `json:"order" binding:"required,number,gte=1" gorm:"not null"`
}

type UpdateBlogDetailsOrder struct {
	Id    uint `json:"id" binding:"required,number,gte=1"`
	Order uint `json:"order" binding:"required,number,gte=1"`
}
