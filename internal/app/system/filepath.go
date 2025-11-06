package system

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type FileSystemService struct {
}

func NewFileSystemService(defaultUploadWritePath string) *FileSystemService {
	return &FileSystemService{}
}

func (service FileSystemService) CheckIfDirectoryExists(createDirIfDontExist bool, targetPath string) (bool, error) {
	var candidatePath = filepath.Join(".", targetPath)

	_, fileCheckingError := os.Stat(candidatePath)

	if errors.Is(fileCheckingError, fs.ErrNotExist) && !createDirIfDontExist {
		return false, nil
	}

	if errors.Is(fileCheckingError, fs.ErrNotExist) && createDirIfDontExist {

		mkdirError := os.MkdirAll(candidatePath, os.ModePerm)

		if mkdirError != nil {
			return false, mkdirError
		}

		return true, nil
	}

	if fileCheckingError != nil {
		return false, fileCheckingError
	}

	return true, nil
}
