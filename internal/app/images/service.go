package images

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
	files "justcallmesu.com/rest-api/internal/app/Files"
	"justcallmesu.com/rest-api/internal/app/system"
	"justcallmesu.com/rest-api/internal/utils"
)

type ImageUploadService struct {
	DefaultWritePath string
	SystemService    *system.SystemService
	FileService      files.FilesService
}

func NewImageUploadService(defaultWritePath string, systemService *system.SystemService, fileService *files.FilesService) *ImageUploadService {
	return &ImageUploadService{
		DefaultWritePath: defaultWritePath,
		SystemService:    systemService,
		FileService:      *fileService,
	}
}

func (service ImageUploadService) GetResizeOptions(imageResolution ImageResolution) bimg.Options {
	return bimg.Options{
		Width:   int(imageResolution),
		Quality: 80,
		Type:    bimg.WEBP,
	}
}

func (service ImageUploadService) ProcessManyImageUpload(context *gin.Context, fieldName string, resolutions []ImageResolution) error {
	form, formError := context.MultipartForm()

	if formError != nil {
		return formError
	}

	files := form.File[fieldName]

	for _, file := range files {

		buffer, bufferError := service.ParseMultipartFileIntoBuffer(file)

		if bufferError != nil {
			return bufferError
		}

		_, imageProcessingError := service.HandleMultiResolutionWrite(buffer, resolutions, file.Filename)

		if imageProcessingError != nil {
			return imageProcessingError
		}

	}

	return nil

}

func (service ImageUploadService) ProcessImageUpload(context *gin.Context, fieldName string, resolutions []ImageResolution) (*files.Files, error) {

	multipartFile, formError := context.FormFile(fieldName)

	if formError != nil {
		return nil, formError
	}

	buffer, bufferError := service.ParseMultipartFileIntoBuffer(multipartFile)

	if bufferError != nil {
		return nil, bufferError
	}

	createdFile, imageProcessingError := service.HandleMultiResolutionWrite(buffer, resolutions, multipartFile.Filename)

	if imageProcessingError != nil {
		return nil, imageProcessingError
	}

	savedFile, fileSaveError := service.FileService.CreateOne(context, createdFile)

	if fileSaveError != nil {
		return nil, fileSaveError
	}

	return savedFile, nil
}

func (service ImageUploadService) ParseMultipartFileIntoBuffer(multipartFile *multipart.FileHeader) ([]byte, error) {

	file, fileOpenError := multipartFile.Open()

	if fileOpenError != nil {
		return nil, fileOpenError
	}

	buffer, bufferError := io.ReadAll(file)

	if bufferError != nil {
		return nil, bufferError
	}

	return buffer, nil

}

func (service ImageUploadService) GetFileName(resolution ImageResolution, originalFileName string) string {
	return filepath.Join(fmt.Sprintf("%s-%s.%s", strconv.Itoa(int(resolution)), utils.SlugGenerator(originalFileName), "webp"))
}

func (service ImageUploadService) HandleMultiResolutionWrite(imageBuffer []byte, resolutions []ImageResolution, originalFileName string) (*files.Files, error) {
	baseDirectorySaveLocation := filepath.Join(service.DefaultWritePath, "webp")
	baseOriginalFileSaveLocation := filepath.Join(service.DefaultWritePath, "original")
	fileNameWithoutExtension := strings.Join(strings.Split(originalFileName, ".")[:1], "")

	var createdFile = &files.Files{
		OriginalName: originalFileName,
		Renditions:   files.Renditions{},
		FileType:     files.IMAGE,
	}

	_, directoryCheckingError := service.SystemService.FileSystemService.CheckIfDirectoryExists(true, baseOriginalFileSaveLocation)

	if directoryCheckingError != nil {
		return nil, directoryCheckingError
	}

	originalFileSaveLocation := filepath.Join(baseOriginalFileSaveLocation, originalFileName)
	writeError := bimg.Write(originalFileSaveLocation, imageBuffer)

	if writeError != nil {
		return nil, writeError
	}

	createdFile.OriginalUrl = originalFileSaveLocation

	for _, resolution := range resolutions {
		bimgOptions := service.GetResizeOptions(resolution)

		newImage, imageProcessingError := bimg.NewImage(imageBuffer).Process(bimgOptions)

		if imageProcessingError != nil {
			return nil, imageProcessingError
		}

		fileName := service.GetFileName(resolution, fileNameWithoutExtension)

		fileSaveLocation := filepath.Join(baseDirectorySaveLocation, strconv.Itoa(int(resolution)))

		_, directoryCheckingError := service.SystemService.FileSystemService.CheckIfDirectoryExists(true, fileSaveLocation)

		if directoryCheckingError != nil {
			return nil, directoryCheckingError
		}

		fileLocation := filepath.Join(fileSaveLocation, fileName)

		writeError := bimg.Write(fileLocation, newImage)

		if writeError != nil {
			return nil, writeError
		}

		switch resolution {
		case SMALL:
			createdFile.Renditions.Small = fileLocation
		case MEDIUM:
			createdFile.Renditions.Medium = fileLocation
		case BIG:
			createdFile.Renditions.Big = fileLocation
		}
	}

	return createdFile, nil
}
