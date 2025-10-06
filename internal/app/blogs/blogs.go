package blogs

import (
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/pkg/types"
)

type BlogQuery struct {
	types.QueryParams
}

type Blog struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"  binding:"required"`
	ReadingTime int32  `json:"readingTime" binding:"required,number,min=1"`
}

func NewBlog(title string, description string, readingTime int32) *Blog {
	return &Blog{
		Title:       title,
		Description: description,
		ReadingTime: readingTime,
	}
}
