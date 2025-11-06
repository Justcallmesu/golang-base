package app

import (
	"os"
	"path/filepath"

	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/api/middleware"
	"justcallmesu.com/rest-api/internal/app/auth"
	"justcallmesu.com/rest-api/internal/app/blogs"
	"justcallmesu.com/rest-api/internal/app/cookies"
	"justcallmesu.com/rest-api/internal/app/images"
	"justcallmesu.com/rest-api/internal/app/system"
	"justcallmesu.com/rest-api/internal/app/users"
	"justcallmesu.com/rest-api/internal/utils"
)

type Services struct {
	CookieService *cookies.TokenCookieService
	AuthService   *auth.AuthService
	JWTService    *auth.JWTService
	BlogService   *blogs.BlogService
	SystemService *system.SystemService
}

type Repositories struct {
	UserRepository *users.UserRepository
	BlogRepository *blogs.BlogRepository
}

type Middlewares struct {
	AuthMiddleware *middleware.AuthMiddleware
}

func NewRepositories(database *gorm.DB) *Repositories {
	userRepository := users.NewUserRepository(database)
	blogRepository := blogs.NewBlogRepository(database)
	return &Repositories{
		UserRepository: userRepository,
		BlogRepository: blogRepository,
	}
}

func NewServices(Repositories *Repositories) *Services {
	DEFAULT_WRITE_PATH := filepath.Clean(utils.GetDefaultValue(os.Getenv("UPLOAD_PATH"), "public/uploads"))

	/**
	System Service
	*/
	systemService := system.NewSystemService(DEFAULT_WRITE_PATH)

	cookieService := cookies.NewTokenCookieService()
	jwtService := auth.NewJWTService()
	imageService := images.NewImageUploadService(DEFAULT_WRITE_PATH, systemService)
	authService := auth.NewAuthService(Repositories.UserRepository, jwtService, cookieService)
	blogService := blogs.NewBlogService(Repositories.BlogRepository, imageService)

	return &Services{
		CookieService: cookieService,
		AuthService:   authService,
		JWTService:    jwtService,
		BlogService:   blogService,
		SystemService: systemService,
	}
}

func NewMiddlewares(Services *Services) *Middlewares {
	return &Middlewares{
		AuthMiddleware: middleware.NewAuthMiddleware(Services.AuthService, Services.JWTService, Services.CookieService),
	}
}
