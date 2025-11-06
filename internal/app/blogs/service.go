package blogs

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/app/images"
)

type BlogService struct {
	BlogRepository *BlogRepository
	imageService   *images.ImageUploadService
}

func NewBlogService(blogRepository *BlogRepository, imageService *images.ImageUploadService) *BlogService {
	return &BlogService{
		BlogRepository: blogRepository,
		imageService:   imageService,
	}
}

func (service *BlogService) FindMany(context *gin.Context, parameter *BlogQuery) (*[]Blog, error) {
	return service.BlogRepository.FindMany(parameter, context)
}

func (service *BlogService) FindOne(context *gin.Context, id int) (Blog, error) {
	return service.BlogRepository.FindOne(id, context)
}

func (service *BlogService) CreateOne(context *gin.Context, newBlog *Blog) error {
	return service.BlogRepository.CreateOne(newBlog, context)
}

func (service *BlogService) UpdateOne(context *gin.Context, updatedBlog *Blog) error {
	updateError := service.BlogRepository.UpdateOne(updatedBlog, context)

	if updateError != nil {
		return updateError
	}

	return nil
}

func (service *BlogService) HandleImageUpload(context *gin.Context, blogId int) error {

	processError := service.imageService.ProcessImageUpload(context, "image", []images.ImageResolution{
		images.BIG,
		images.MEDIUM,
		images.SMALL,
	})

	fmt.Println(processError)

	if processError != nil {
		return processError
	}

	return nil

}

func (service *BlogService) DeleteOne(context *gin.Context, targetId int) error {
	deleteError := service.BlogRepository.DeleteOne(targetId, context)

	if deleteError != nil {
		return deleteError
	}

	return nil
}
