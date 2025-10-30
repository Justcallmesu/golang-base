package blogs

import (
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
	Type    BlogDetailsType `json:"type"`
}
