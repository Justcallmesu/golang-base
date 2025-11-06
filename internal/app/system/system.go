package system

type SystemService struct {
	FileSystemService *FileSystemService
}

func NewSystemService(fileSystemService *FileSystemService) *SystemService {
	return &SystemService{
		FileSystemService: fileSystemService,
	}
}
