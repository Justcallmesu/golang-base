package blogs

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/response"
	"justcallmesu.com/rest-api/internal/utils"
)

type BlogHandler struct {
	BlogService *BlogService
}

func NewBlogHandler(blogService *BlogService) *BlogHandler {
	return &BlogHandler{
		BlogService: blogService,
	}
}

func (handler *BlogHandler) FindMany(context *gin.Context) {
	var blogQuery BlogQuery

	queryParseError := context.ShouldBindQuery(&blogQuery)

	if queryParseError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Terjadi kesalahan", utils.FormatValidationError(queryParseError)))
		return
	}

	blogs, findError := handler.BlogService.FindMany(context, &blogQuery)

	if findError != nil {
		context.JSON(http.StatusInternalServerError, response.NewErrorResponse("Terjadi kesalahan pada server", findError.Error()))
		return
	}

	context.JSON(http.StatusOK, response.NewResponse("Success", true, response.NewResponse("Blog berhasil ditemukan", true, blogs)))
}

func (handler *BlogHandler) FindOne(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Id diperlukan", nil))
		return
	}

	parsedId, parseError := strconv.Atoi(id)

	if parseError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Id harus berupa angka", nil))
		return
	}

	blog, findError := handler.BlogService.FindOne(context, parsedId)

	if findError != nil {
		context.JSON(http.StatusInternalServerError, response.NewErrorResponse(findError.Error(), nil))
		return
	}

	context.JSON(http.StatusOK, response.NewResponse("Success", true, response.NewResponse("Blog berhasil ditemukan", true, blog)))
}

func (handler *BlogHandler) Create(context *gin.Context) {
	var newBlog Blog

	if err := context.ShouldBindJSON(&newBlog); err != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Data tidak sesuai", utils.FormatValidationError(err)))
		return
	}

	createError := handler.BlogService.CreateOne(context, &newBlog)

	if createError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Terjadi kesalahan saat menyimpan", utils.FormatValidationError(createError)))
		return
	}

	context.JSON(http.StatusCreated, response.NewResponse("Data berhasil dibuat", true, nil))
}

func (handler *BlogHandler) Update(context *gin.Context) {
	var updatedBlog Blog

	if err := context.ShouldBindJSON(&updatedBlog); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Data tidak sesuai", utils.FormatValidationError(err)))
		return
	}

	updateError := handler.BlogService.UpdateOne(context, &updatedBlog)

	if updateError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Gagal saat mengubah", utils.FormatValidationError(updateError)))
		return
	}

	context.JSON(http.StatusAccepted, response.NewResponse("Blog berhasil diedit", true, nil))
}

func (handler *BlogHandler) Upload(context *gin.Context) {
	blogId, paramError := utils.GetAndParseParamToInt(context, "id", 10, 32)

	if paramError != nil {
		context.JSON(http.StatusBadRequest, paramError)
		return
	}

	imageUploadError := handler.BlogService.HandleImageUpload(context, int(blogId))

	if imageUploadError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Gagal Memproses Upload", errors.New("silahkan upload kembali")))
		return
	}

}

func (handler *BlogHandler) Delete(context *gin.Context) {
	blogId := context.Param("id")

	if blogId == "" {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Parameter tidak valid", errors.New("parameter tidak ditemukan")))
		return
	}

	value, parsingError := strconv.ParseInt(blogId, 10, 64)

	if parsingError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Parameter harus berupa angka", nil))
	}

	deleteError := handler.BlogService.DeleteOne(context, int(value))

	if deleteError != nil {

		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Gagal menghapus Blog", utils.FormatValidationError(deleteError)))
		return
	}

	context.JSON(http.StatusNoContent, response.NewResponse("Blog berhasil dihapus", true, nil))
}
