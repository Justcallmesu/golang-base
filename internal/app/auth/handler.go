package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/rest-api/internal/api/response"
	"justcallmesu.com/rest-api/internal/app/cookies"
	application_error "justcallmesu.com/rest-api/internal/utils/error"
)

type AuthHandler struct {
	AuthService   *AuthService
	CookieService *cookies.TokenCookieService
}

func NewAuthHandler(authService *AuthService, cookieService *cookies.TokenCookieService) *AuthHandler {

	return &AuthHandler{
		AuthService: authService,
		CookieService: cookieService,
	}
}

func (handler *AuthHandler) Login(context *gin.Context) {
	var credentials = LoginUser{}

	bindError := context.ShouldBindJSON(&credentials)

	if bindError != nil {
		context.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", application_error.FormatValidationError(bindError)))
		return
	}

	tokens, loginError := handler.AuthService.Login(context, &credentials)

	if loginError != nil {
		context.JSON(http.StatusUnauthorized, response.NewResponse(loginError.Error(), false, nil))
		return
	}

	handler.CookieService.GenerateTokenCookies(context, tokens)

	context.JSON(http.StatusCreated, response.NewResponse("Login Success", true, nil))
}

func (handler *AuthHandler) Logout(context *gin.Context) {

	resetTokenError :=	handler.CookieService.ResetTokenCookies(context)

	if(resetTokenError != nil) {
		context.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to reset cookies", resetTokenError.Error()))
		return
	}


	context.JSON(http.StatusAccepted, response.NewResponse("Logout Success", true, nil))
}
