package blogs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/response"
	"justcallmesu.com/rest-api/internal/utils"
)

type BlogDetailsHandler struct {
	blogDetailService *BlogDetailService
}

func NewBlogDetailsHandler(blogDetailService *BlogDetailService) *BlogDetailsHandler {
	return &BlogDetailsHandler{
		blogDetailService: blogDetailService,
	}
}

func (handler *BlogDetailsHandler) CreateOne(context *gin.Context) {
	var newBlogDetail BlogDetails

	if bindError := context.ShouldBindJSON(&newBlogDetail); bindError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Validasi gagal", utils.FormatValidationError(bindError)))
		return
	}

	createError := handler.blogDetailService.CreateOne(context, &newBlogDetail)

	if createError != nil {
		context.JSON(http.StatusInternalServerError, response.NewErrorResponse("Gagal membuat blog detail", createError.Error()))
		return
	}

	context.JSON(http.StatusCreated, response.NewResponse("Blog detail berhasil dibuat", true, newBlogDetail))
}

func (handler *BlogDetailsHandler) UpdateDetailsOrder(context *gin.Context) {
	var updatedBlogDetail []UpdateBlogDetailsOrder

	if bindError := context.ShouldBindJSON(&updatedBlogDetail); bindError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Validasi gagal", utils.FormatValidationError(bindError)))
		return
	}

	updateError := handler.blogDetailService.UpdateDetailsOrder(context, &updatedBlogDetail)

	if updateError != nil {
		context.JSON(http.StatusInternalServerError, response.NewErrorResponse("Gagal memperbarui urutan blog detail", updateError.Error()))
		return
	}

	context.JSON(http.StatusOK, response.NewResponse("Urutan blog detail berhasil diperbarui", true, nil))

}
