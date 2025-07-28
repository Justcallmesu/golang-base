package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/cookies"
	"justcallmesu.com/rest-api/internal/api/response"
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
		context.JSON(http.StatusBadRequest, response.NewResponse(fmt.Sprintf("Error Signing You Up: %s", createError.Error()), false, nil))
		return
	}

	context.JSON(http.StatusCreated, response.NewResponse("User Created", true, nil))
}

func (handler *AuthHandler) Login(context *gin.Context) {
	jwtToken, loginError := handler.AuthService.Login(context)

	if loginError != nil {
		context.JSON(http.StatusUnauthorized, response.NewResponse(loginError.Error(), false, nil))
		return
	}

	cookieMaxAge, parseError := strconv.ParseInt(os.Getenv("COOKIE_EXPIRATION"), 10, 64)

	if parseError != nil {
		context.JSON(http.StatusUnauthorized, response.NewResponse(parseError.Error(), false, nil))
		return
	}

	setCookieError := cookies.SetCookie(context, os.Getenv("COOKIE_NAME"), jwtToken, cookieMaxAge)

	if setCookieError != nil {
		context.JSON(http.StatusUnauthorized, response.NewResponse(setCookieError.Error(), false, nil))
		return
	}

	context.JSON(http.StatusCreated, response.NewResponse("Login Success", true, nil))
}

func (handler *AuthHandler) Logout(context *gin.Context) {
	handler.AuthService.Logout(context)
	context.JSON(http.StatusAccepted, response.NewResponse("Logout Success", true, nil))
}
