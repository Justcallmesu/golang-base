package app

import (
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/api/middleware"
	files "justcallmesu.com/rest-api/internal/app/Files"
	"justcallmesu.com/rest-api/internal/app/auth"
	"justcallmesu.com/rest-api/internal/app/blogs"
	"justcallmesu.com/rest-api/internal/app/cookies"
	"justcallmesu.com/rest-api/internal/app/images"
	"justcallmesu.com/rest-api/internal/app/system"
	"justcallmesu.com/rest-api/internal/app/users"
)

type Services struct {
	CookieService *cookies.TokenCookieService
	AuthService   *auth.AuthService
	JWTService    *auth.JWTService
	BlogService   *blogs.BlogService
	SystemService *system.SystemService
	FileService   *files.FilesService
}

type Repositories struct {
	UserRepository        *users.UserRepository
	BlogRepository        *blogs.BlogRepository
	blogDetailsRepository *blogs.BlogDetailsRepository
	FileRepository        *files.FilesRepository
}

type Middlewares struct {
	AuthMiddleware *middleware.AuthMiddleware
}

func NewRepositories(database *gorm.DB) *Repositories {
	userRepository := users.NewUserRepository(database)
	blogDetailsRepository := blogs.NewBlogDetailsRepository(database)
	blogRepository := blogs.NewBlogRepository(database)
	fileRepository := files.NewFilesRepository(database)
	return &Repositories{
		UserRepository:        userRepository,
		BlogRepository:        blogRepository,
		FileRepository:        fileRepository,
		blogDetailsRepository: blogDetailsRepository,
	}
}

func NewServices(Repositories *Repositories, defaultWritePath string) *Services {

	/**
	System Service
	*/
	systemService := system.NewSystemService(defaultWritePath)

	/**
	Utils Service
	*/
	fileService := files.NewFilesService(Repositories.FileRepository)
	jwtService := auth.NewJWTService()
	cookieService := cookies.NewTokenCookieService()
	imageService := images.NewImageUploadService(defaultWritePath, systemService, fileService)

	authService := auth.NewAuthService(Repositories.UserRepository, jwtService, cookieService)
	blogService := blogs.NewBlogService(Repositories.BlogRepository, Repositories.blogDetailsRepository, imageService)

	return &Services{
		CookieService: cookieService,
		AuthService:   authService,
		JWTService:    jwtService,
		BlogService:   blogService,
		SystemService: systemService,
		FileService:   fileService,
	}
}

func NewMiddlewares(Services *Services) *Middlewares {
	return &Middlewares{
		AuthMiddleware: middleware.NewAuthMiddleware(Services.AuthService, Services.JWTService, Services.CookieService),
	}
}
