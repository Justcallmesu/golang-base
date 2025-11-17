package blogs

import (
	"justcallmesu.com/rest-api/internal/types"
)

type BlogQuery struct {
	types.QueryParams
}

type Blog struct {
	types.BaseEntityModel
	Title       string        `json:"title,omitempty" binding:"required"`
	Description string        `json:"description,omitempty"  binding:"required"`
	ReadingTime int32         `json:"readingTime,omitempty" binding:"required,number,gte=1"`
	ViewCount   uint64        `json:"viewCount"`
	Details     []BlogDetails `json:"details,omitzero" gorm:"foreignKey:BlogId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;references:Id"`
}

func NewBlog(title string, description string, readingTime int32, blogDetails []BlogDetails) *Blog {
	return &Blog{
		Title:       title,
		Description: description,
		ReadingTime: readingTime,
		Details:     blogDetails,
	}
}
