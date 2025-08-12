package app

import (
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/api/middleware"
	"justcallmesu.com/rest-api/internal/app/auth"
	"justcallmesu.com/rest-api/internal/app/cookies"
	"justcallmesu.com/rest-api/internal/app/users"
)

type Services struct {
	CookieService *cookies.TokenCookieService
	AuthService   *auth.AuthService
	JWTService    *auth.JWTService
}

type Repositories struct {
	UserRepository *users.UserRepository
}

type Middlewares struct {
	AuthMiddleware *middleware.AuthMiddleware
}

func NewRepositories(database *gorm.DB) *Repositories {
	userRepository := users.NewUserRepository(database)

	return &Repositories{
		UserRepository: userRepository,
	}
}

func NewServices(Repositories *Repositories) *Services {
	cookieService := cookies.NewTokenCookieService()
	jwtService := auth.NewJWTService()
	authService := auth.NewAuthService(Repositories.UserRepository,jwtService, cookieService)

	return &Services{
		CookieService: cookieService,
		AuthService:   authService,
		JWTService:    jwtService,
	}
}


func NewMiddlewares(Services *Services) *Middlewares {
	return &Middlewares{
		AuthMiddleware: middleware.NewAuthMiddleware(Services.AuthService, Services.JWTService),
	}
}