package blogs

import (
	"context"

	"gorm.io/gorm"
)

type BlogDetailsRepository struct {
	database *gorm.DB
}

func NewBlogDetailsRepository(database *gorm.DB) *BlogDetailsRepository {
	return &BlogDetailsRepository{
		database: database,
	}
}

func (repository BlogDetailsRepository) CreateOne(blogDetails *BlogDetails, context context.Context) error {
	result := gorm.WithResult()
	createError := gorm.G[BlogDetails](repository.database, result).Create(context, blogDetails)

	return createError
}
