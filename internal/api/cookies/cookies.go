package cookies

import (
	"fmt"
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

	RefreshTokenName string
	AccessTokenName  string
}

func NewService() (*Service, error) {
	refreshTokenExpiration, err := strconv.ParseInt(os.Getenv("COOKIE_REFRESH_EXPIRATION"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid COOKIE_REFRESH_EXPIRATION value: %w", err)
	}

	accessTokenExpiration, err := strconv.ParseInt(os.Getenv("COOKIE_ACCESS_EXPIRATION"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid COOKIE_ACCESS_EXPIRATION value: %w", err)
	}

	refreshTokenName := os.Getenv("COOKIE_REFRESH_TOKEN")
	if refreshTokenName == "" {
		return nil, fmt.Errorf("environment variable COOKIE_REFRESH_TOKEN is not set")
	}

	accessTokenName := os.Getenv("COOKIE_ACCESS_TOKEN")
	if accessTokenName == "" {
		return nil, fmt.Errorf("environment variable COOKIE_ACCESS_TOKEN is not set")
	}

	return &Service{
		refreshTokenExpiration: time.Duration(refreshTokenExpiration) * time.Hour,
		accessTokenExpiration:  time.Duration(accessTokenExpiration) * time.Minute,
		RefreshTokenName:       os.Getenv("COOKIE_REFRESH_TOKEN"),
		AccessTokenName:        os.Getenv("COOKIE_ACCESS_TOKEN"),
	}, nil
}

func (service *Service) GenerateRefreshCookies(context *gin.Context, refreshToken string) {
	cookie := &http.Cookie{
		Name:     service.RefreshTokenName,
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

	cookie := &http.Cookie{
		Name:     service.AccessTokenName,
		Value:    accessToken,
		MaxAge:   int(service.accessTokenExpiration.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   service.isSecureRequest(context), // Set to true if using HTTPS (recommended)
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(context.Writer, cookie)
}

func (service *Service) ResetTokenCookies(context *gin.Context) error {

	context.SetCookie(service.RefreshTokenName, "", -1, "/", "", false, true)
	context.SetCookie(service.AccessTokenName, "", -1, "/", "", false, true)

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
