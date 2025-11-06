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
	"justcallmesu.com/rest-api/internal/app/system"
	"justcallmesu.com/rest-api/internal/utils"
)

type ImageUploadService struct {
	DefaultWritePath string
	SystemService    *system.SystemService
}

func NewImageUploadService(defaultWritePath string, systemService *system.SystemService) *ImageUploadService {
	return &ImageUploadService{
		DefaultWritePath: defaultWritePath,
		SystemService:    systemService,
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

func (service ImageUploadService) ProcessImageUpload(context *gin.Context, fieldName string, resolutions []ImageResolution) error {

	multipartFile, formError := context.FormFile(fieldName)

	if formError != nil {
		return formError
	}

	buffer, bufferError := service.ParseMultipartFileIntoBuffer(multipartFile)

	if bufferError != nil {
		return bufferError
	}

	_, imageProcessingError := service.HandleMultiResolutionWrite(buffer, resolutions, multipartFile.Filename)

	if imageProcessingError != nil {
		return imageProcessingError
	}

	return nil
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

func (service ImageUploadService) HandleMultiResolutionWrite(imageBuffer []byte, resolutions []ImageResolution, originalFileName string) ([]ImageResolutionWriteReturn, error) {
	for _, resolution := range resolutions {
		bimgOptions := service.GetResizeOptions(resolution)

		newImage, imageProcessingError := bimg.NewImage(imageBuffer).Process(bimgOptions)

		if imageProcessingError != nil {
			return nil, imageProcessingError
		}

		originalFileName := strings.Join(strings.Split(originalFileName, ".")[:1], "")

		fileName := service.GetFileName(resolution, originalFileName)

		diskSaveDestination := filepath.Join(service.DefaultWritePath, "webp", strconv.Itoa(int(resolution)))

		isExist, directoryCheckingError := service.SystemService.FileSystemService.CheckIfDirectoryExists(true, diskSaveDestination)

		if directoryCheckingError != nil {
			return nil, directoryCheckingError
		}

		if !isExist {
			return nil, fmt.Errorf("%s is doesn't exist in the machine", diskSaveDestination)
		}

		writeError := bimg.Write(filepath.Join(diskSaveDestination, fileName), newImage)

		if writeError != nil {
			return nil, writeError
		}

	}

	return nil, nil
}
