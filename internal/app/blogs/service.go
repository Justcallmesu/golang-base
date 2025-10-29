package blogs

import (
	"context"
)

type BlogService struct {
	BlogRepository *BlogRepository
}

func NewBlogService(blogRepository *BlogRepository) *BlogService {
	return &BlogService{
		BlogRepository: blogRepository,
	}
}

func (service *BlogService) FindMany(context context.Context, parameter *BlogQuery) (*[]Blog, error) {
	return service.BlogRepository.FindMany(parameter, context)
}

func (service *BlogService) FindOne(context context.Context, id int) (Blog, error) {
	return service.BlogRepository.FindOne(id, context)
}

func (service *BlogService) CreateOne(context context.Context, newBlog *Blog) error {
	return service.BlogRepository.CreateOne(newBlog, context)
}

func (service *BlogService) UpdateOne(context context.Context, updatedBlog *Blog) error {
	updateError := service.BlogRepository.UpdateOne(updatedBlog, context)

	if updateError != nil {
		return updateError
	}

	return nil
}

func (service *BlogService) DeleteOne(context context.Context, targetId int) error {
	deleteError := service.BlogRepository.DeleteOne(targetId, context)

	if deleteError != nil {
		return deleteError
	}

	return nil
}
