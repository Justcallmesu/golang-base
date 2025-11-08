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
	Url     string          `json:"url,omitempty"`
	FileId  uint            `json:"fileId,omitempty"`
	File    files.Files     `json:"file,omitzero"`
	Type    BlogDetailsType `json:"type"`
}
