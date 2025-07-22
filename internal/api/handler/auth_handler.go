package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api"
	"justcallmesu.com/rest-api/internal/api/cookies"
	"justcallmesu.com/rest-api/internal/app/auth"
	"justcallmesu.com/rest-api/internal/app/users"
)

type AuthHandler struct {
	AuthService *auth.AuthService
}

func NewAuthHandler(database *sql.DB) *AuthHandler {
	UserRepository := users.NewUserRepository(database)
	authService := auth.NewAuthService(UserRepository)

	return &AuthHandler{
		AuthService: authService,
	}
}

func (handler *AuthHandler) SignUp(context *gin.Context) {
	_, createError := handler.AuthService.SignUp(context)

	if createError != nil {
		context.JSON(http.StatusBadRequest, api.NewResponse(fmt.Sprintf("Error Signing You Up: %s", createError.Error()), false, nil))
		return
	}

	context.JSON(http.StatusCreated, api.NewResponse("User Created", true, nil))
}

func (handler *AuthHandler) Login(context *gin.Context) {
	jwtToken, loginError := handler.AuthService.Login(context)

	if loginError != nil {
		context.JSON(http.StatusUnauthorized, api.NewResponse(loginError.Error(), false, nil))
		return
	}

	cookieMaxAge, parseError := strconv.ParseInt(os.Getenv("COOKIE_EXPIRATION"), 10, 64)

	if parseError != nil {
		context.JSON(http.StatusUnauthorized, api.NewResponse(parseError.Error(), false, nil))
		return
	}

	cookies.SetCookie(context, os.Getenv("COOKIE_NAME"), jwtToken, cookieMaxAge)
	context.JSON(http.StatusCreated, api.NewResponse("Login Success", true, nil))
}

func (handler *AuthHandler) Logout(context *gin.Context) {
	handler.AuthService.Logout(context)
	context.JSON(http.StatusAccepted, api.NewResponse("Logout Success", true, nil))
}
