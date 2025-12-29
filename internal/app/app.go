package app

import (
	"gorm.io/gorm"
	"justcallmesu.com/golang-base/internal/api/middleware"
	"justcallmesu.com/golang-base/internal/app/auth"
	"justcallmesu.com/golang-base/internal/app/cookies"
	"justcallmesu.com/golang-base/internal/app/tokens"
	"justcallmesu.com/golang-base/internal/app/users"
)

type Services struct {
	CookieService *cookies.TokenCookieService
	AuthService   *auth.Service
	JWTService    *tokens.JWTService
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
	jwtService := tokens.NewJWTService()
	authService := auth.NewService(Repositories.UserRepository, jwtService, cookieService)

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
