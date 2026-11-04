package app

import (
	"gorm.io/gorm"
	cookies_service "justcallmesu.com/golang-base/internal/api/cookie"
	"justcallmesu.com/golang-base/internal/api/middleware"
	role_repository "justcallmesu.com/golang-base/internal/app/repositories/role"
	user_repository "justcallmesu.com/golang-base/internal/app/repositories/user"
	auth_service "justcallmesu.com/golang-base/internal/app/services/auth"
	navigation_menu_service "justcallmesu.com/golang-base/internal/app/services/navigation-menu"
	role_service "justcallmesu.com/golang-base/internal/app/services/role"
	tokens "justcallmesu.com/golang-base/internal/app/services/token"
	user_service "justcallmesu.com/golang-base/internal/app/services/user"
)

type Services struct {
	CookieService         *cookies_service.Service
	AuthService           *auth_service.Service
	JWTService            *tokens.Service
	NavigationMenuService *navigation_menu_service.Service

	UserService *user_service.Service

	RoleService *role_service.Service
}

type Repositories struct {
	UserRepository *user_repository.UserRepository
	RoleRepository *role_repository.Repository
}

type Middlewares struct {
	AuthMiddleware *middleware.AuthMiddleware
}

func NewRepositories(database *gorm.DB) *Repositories {
	userRepository := user_repository.NewUserRepository(database)
	roleRepository := role_repository.NewRepository(database)

	return &Repositories{
		UserRepository: userRepository,
		RoleRepository: roleRepository,
	}
}

func NewServices(Repositories *Repositories) *Services {
	cookieService := cookies_service.NewService()
	jwtService := tokens.NewService()
	navigationMenuService := navigation_menu_service.NewService()
	authService := auth_service.NewService(Repositories.UserRepository, jwtService, cookieService, navigationMenuService)

	// User Service
	userService := user_service.NewService(Repositories.UserRepository)

	// Role Service
	roleService := role_service.NewService(Repositories.RoleRepository)

	return &Services{
		CookieService:         cookieService,
		AuthService:           authService,
		JWTService:            jwtService,
		NavigationMenuService: navigationMenuService,
		UserService:           userService,
		RoleService:           roleService,
	}
}

func NewMiddlewares(Services *Services) *Middlewares {
	return &Middlewares{
		AuthMiddleware: middleware.NewAuthMiddleware(Services.AuthService, Services.JWTService),
	}
}
