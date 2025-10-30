package app

import (
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/api/middleware"
	"justcallmesu.com/rest-api/internal/app/auth"
	"justcallmesu.com/rest-api/internal/app/blogs"
	"justcallmesu.com/rest-api/internal/app/cookies"
	"justcallmesu.com/rest-api/internal/app/users"
)

type Services struct {
	CookieService *cookies.TokenCookieService
	AuthService   *auth.AuthService
	JWTService    *auth.JWTService
	BlogService   *blogs.BlogService
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
	cookieService := cookies.NewTokenCookieService()
	jwtService := auth.NewJWTService()
	authService := auth.NewAuthService(Repositories.UserRepository, jwtService, cookieService)
	blogService := blogs.NewBlogService(Repositories.BlogRepository)

	return &Services{
		CookieService: cookieService,
		AuthService:   authService,
		JWTService:    jwtService,
		BlogService:   blogService,
	}
}

func NewMiddlewares(Services *Services) *Middlewares {
	return &Middlewares{
		AuthMiddleware: middleware.NewAuthMiddleware(Services.AuthService, Services.JWTService, Services.CookieService),
	}
}
