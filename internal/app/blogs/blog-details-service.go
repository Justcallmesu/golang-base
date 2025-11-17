package blogs

import "github.com/gin-gonic/gin"

type BlogDetailService struct {
	blogDetailRepository *BlogDetailsRepository
}

func NewBlogDetailsService(blogDetailRepository BlogDetailsRepository) *BlogDetailService {
	return &BlogDetailService{
		blogDetailRepository: &blogDetailRepository,
	}
}

func (service *BlogDetailService) CreateOne(context *gin.Context, newBlog *BlogDetails) error {
	return service.blogDetailRepository.CreateOne(newBlog, context)
}

func (service *BlogDetailService) UpdateOne(context *gin.Context, updatedBlog *BlogDetails) error {
	updateError := service.blogDetailRepository.UpdateOne(updatedBlog, context)

	if updateError != nil {
		return updateError
	}

	return nil
}

func (service *BlogDetailService) UpdateDetailsOrder(context *gin.Context, updateBlogDetailOrder *[]UpdateBlogDetailsOrder) error {
	return service.blogDetailRepository.UpdateDetailOrder(updateBlogDetailOrder)
}
