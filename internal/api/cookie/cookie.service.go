package cookies_service

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	cookie_dtos "justcallmesu.com/golang-base/internal/app/dtos/cookie"
)

type Service struct {
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration

	refreshTokenName string
	accessTokenName  string
}

func NewService() *Service {
	refreshTokenExpiration, err := strconv.ParseInt(os.Getenv("COOKIE_REFRESH_EXPIRATION"), 10, 64)
	if err != nil {
		panic("Invalid COOKIE_REFRESH_EXPIRATION value")
	}

	accessTokenExpiration, err := strconv.ParseInt(os.Getenv("COOKIE_ACCESS_EXPIRATION"), 10, 64)
	if err != nil {
		panic("Invalid COOKIE_ACCESS_EXPIRATION value")
	}

	return &Service{
		refreshTokenExpiration: time.Duration(refreshTokenExpiration) * time.Hour,
		accessTokenExpiration:  time.Duration(accessTokenExpiration) * time.Minute,
		refreshTokenName:       os.Getenv("COOKIE_REFRESH_TOKEN"),
		accessTokenName:        os.Getenv("COOKIE_ACCESS_TOKEN"),
	}
}

func (service *Service) GenerateRefreshCookies(context *gin.Context, refreshToken string) {
	cookie := &http.Cookie{
		Name:     service.refreshTokenName,
		Value:    refreshToken,
		MaxAge:   int(service.refreshTokenExpiration.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   service.isSecureRequest(context), // Set to true if using HTTPS (recommended)
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(context.Writer, cookie)
}

func (service *Service) GenerateAccessCookies(context *gin.Context, accessToken string) {

	http.SetCookie(context.Writer, service.generateCookie(service.accessTokenName, accessToken, int(service.accessTokenExpiration.Seconds()), context))
}

func (service *Service) ResetTokenCookies(context *gin.Context) error {

	context.SetCookie(service.refreshTokenName, "", -1, "/", "", false, true)
	context.SetCookie(service.accessTokenName, "", -1, "/", "", false, true)

	return nil
}

func (service *Service) GenerateTokenCookies(context *gin.Context, tokenCookie *cookie_dtos.TokenCookie) {

	if tokenCookie.RefreshToken != "" {
		service.GenerateRefreshCookies(context, tokenCookie.RefreshToken)
	}

	if tokenCookie.AccessToken != "" {
		service.GenerateAccessCookies(context, tokenCookie.AccessToken)
	}
}

func NewTokenCookie(refreshToken, accessToken string) *cookie_dtos.TokenCookie {
	return &cookie_dtos.TokenCookie{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
}

func (service *Service) isSecureRequest(context *gin.Context) bool {
	return context.Request.URL.Scheme == "https" || context.Request.TLS != nil
}

func (service *Service) generateCookie(name, value string, maxAge int, context *gin.Context) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   service.isSecureRequest(context),
		SameSite: http.SameSiteLaxMode,
	}

}
