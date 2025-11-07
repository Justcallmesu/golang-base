package files

import "github.com/gin-gonic/gin"

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
