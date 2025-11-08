package files

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type FilesService struct {
	fileRepository *FilesRepository
}

func NewFilesService(fileRepository *FilesRepository) *FilesService {
	return &FilesService{
		fileRepository: fileRepository,
	}
}

func (service FilesService) FindMany(context *gin.Context, query FilesQuery) (*[]Files, error) {
	return service.fileRepository.FindAll(query, context)
}

func (service *FilesService) FindOne(context *gin.Context, id int) (*Files, error) {
	return service.fileRepository.FindOne(id, context)
}

func (service *FilesService) CreateOne(context *gin.Context, file *Files) (*Files, error) {
	createError := service.fileRepository.CreateOne(file, context)

	if createError != nil {
		return nil, createError
	}

	fmt.Printf("%v", file)

	return file, nil
}
