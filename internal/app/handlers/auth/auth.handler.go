package auth_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/request"
	"justcallmesu.com/golang-base/internal/api/response"
	auth_dtos "justcallmesu.com/golang-base/internal/app/dtos/auth"
	auth_service "justcallmesu.com/golang-base/internal/app/services/auth"
	cookies_service "justcallmesu.com/golang-base/internal/api/cookie"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type AuthHandler struct {
	service       *auth_service.Service
	cookieService *cookies_service.Service
}

func NewAuthHandler(authService *auth_service.Service, cookieService *cookies_service.Service) *AuthHandler {
	return &AuthHandler{
		service:       authService,
		cookieService: cookieService,
	}
}

func (handler *AuthHandler) Login(context *gin.Context) {
	var credentials = auth_dtos.LoginUser{}

	if bindError := context.ShouldBindJSON(&credentials); bindError != nil {
		context.Error(request.BindError(bindError))
		return
	}

	tokens, userData, navigationMenu, loginError := handler.service.Login(context.Request.Context(), &credentials)

	if loginError != nil {
		context.Error(loginError)
		return
	}
	handler.cookieService.GenerateTokenCookies(context, tokens)

	context.JSON(http.StatusCreated, response.NewResponse("Login Success", true, gin.H{
		"user":           userData,
		"navigationMenu": navigationMenu,
	}))
}

func (handler *AuthHandler) Me(context *gin.Context) {

	jwtClaims, getError := utils.GetUserData(context)

	if getError != nil {
		context.Error(getError)
		return
	}

	userData, navigationMenu, getError := handler.service.GetMe(context.Request.Context(), jwtClaims)

	if getError != nil {
		context.Error(getError)
		return
	}

	context.JSON(http.StatusOK, response.NewResponse("User fetched successfully", true, gin.H{
		"user":           userData,
		"navigationMenu": navigationMenu,
	}))

}

func (handler *AuthHandler) ChangePassword(context *gin.Context) {
	userData, getError := utils.GetUserData(context)

	if getError != nil {
		context.Error(getError)
		return
	}

	var changePasswordDto auth_dtos.ChangePassword

	bindingError := context.BindJSON(&changePasswordDto)

	if bindingError != nil {
		context.Error(bindingError)
		return
	}

	updatePasswordError := handler.service.ChangePassword(context.Request.Context(), &changePasswordDto, userData)

	if updatePasswordError != nil {
		context.Error(updatePasswordError)
		return
	}

	context.JSON(http.StatusOK, response.NewResponse("Password Changed Successfully", true, nil))
}

func (handler *AuthHandler) Logout(context *gin.Context) {

	resetTokenError := handler.cookieService.ResetTokenCookies(context)

	if resetTokenError != nil {
		context.Error(resetTokenError)
		return
	}

	context.JSON(http.StatusAccepted, response.NewResponse("Logout Success", true, nil))
}
