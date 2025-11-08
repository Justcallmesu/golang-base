package system

import (
	"os"
	"path/filepath"
)

type FileSystemService struct {
}

func NewFileSystemService(defaultUploadWritePath string) *FileSystemService {
	return &FileSystemService{}
}

func (service FileSystemService) CheckIfDirectoryExists(createDirIfDontExist bool, targetPath string) (bool, error) {
	var candidatePath = filepath.Join("./", targetPath)

	_, fileCheckingError := os.Stat(candidatePath)

	if createDirIfDontExist {

		mkdirError := os.MkdirAll(candidatePath, os.ModePerm)

		if mkdirError != nil {
			return false, mkdirError
		}

		return true, nil
	}

	if os.IsNotExist(fileCheckingError) {
		return false, nil
	}

	if fileCheckingError != nil {
		return false, fileCheckingError
	}

	return true, nil
}
