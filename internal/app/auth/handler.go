package auth

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/api/cookies"
	"justcallmesu.com/rest-api/internal/api/response"
	"justcallmesu.com/rest-api/internal/app/users"
)

type AuthHandler struct {
	AuthService *AuthService
}

func NewAuthHandler(database *gorm.DB) *AuthHandler {
	UserRepository := users.NewUserRepository(database)
	authService := NewAuthService(UserRepository)

	return &AuthHandler{
		AuthService: authService,
	}
}

func (handler *AuthHandler) Login(context *gin.Context) {
	var credentials = LoginUser{}

	bindError := context.ShouldBindJSON(&credentials)

	if bindError != nil {
		context.JSON(http.StatusBadRequest, response.NewResponse(bindError.Error(), false, nil))
		return
	}

	jwtToken, loginError := handler.AuthService.Login(context, &credentials)

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

	setCookieError := cookies.SetCookie(context, os.Getenv("COOKIE_NAME"), "", -1)

	if setCookieError != nil {
		context.JSON(http.StatusInternalServerError, response.NewResponse(setCookieError.Error(), false, nil))
		return
	}

	context.JSON(http.StatusAccepted, response.NewResponse("Logout Success", true, nil))
}
