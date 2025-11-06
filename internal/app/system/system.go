package system

type SystemService struct {
	FileSystemService *FileSystemService
}

func NewSystemService(defaultUploadWritePath string) *SystemService {
	return &SystemService{
		FileSystemService: NewFileSystemService(defaultUploadWritePath),
	}
}
