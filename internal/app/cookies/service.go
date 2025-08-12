package cookies

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenCookie struct {
	RefreshToken string
	AccessToken  string
}

func NewTokenCookie(refreshToken, accessToken string) *TokenCookie {
	return &TokenCookie{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
}

type TokenCookieService struct {
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration

	refreshTokenName string
	accessTokenName  string
}

func NewTokenCookieService() *TokenCookieService {
	refreshTokenExpiration, err := strconv.ParseInt(os.Getenv("COOKIE_REFRESH_EXPIRATION"), 10, 64) 
	if err != nil {
		panic("Invalid COOKIE_REFRESH_EXPIRATION value")
	}

	accessTokenExpiration, err := strconv.ParseInt(os.Getenv("COOKIE_ACCESS_EXPIRATION"), 10, 64)
	if err != nil {
		panic("Invalid COOKIE_ACCESS_EXPIRATION value")
	}

	return &TokenCookieService{
		refreshTokenExpiration: time.Duration(refreshTokenExpiration) * time.Hour,
		accessTokenExpiration:  time.Duration(accessTokenExpiration) * time.Minute,
		refreshTokenName:       os.Getenv("COOKIE_REFRESH_TOKEN"),
		accessTokenName:        os.Getenv("COOKIE_ACCESS_TOKEN"),
	}
}

func (service *TokenCookieService) GenerateRefreshCookies(context *gin.Context, refreshToken string) {
	context.SetCookie(service.refreshTokenName, refreshToken, int(service.refreshTokenExpiration.Seconds()), "/", "", false, true)
}

func (service *TokenCookieService) GenerateAccessCookies(context *gin.Context, accessToken string)  {
	context.SetCookie(service.accessTokenName, accessToken, int(service.accessTokenExpiration.Seconds()), "/", "", false, true)
}

func (service *TokenCookieService) ResetTokenCookies(context *gin.Context) error {

	context.SetCookie(service.refreshTokenName, "", -1, "/", "", false, true)
	context.SetCookie(service.accessTokenName, "", -1, "/", "", false, true)

	return nil
}

func (service *TokenCookieService) GenerateTokenCookies(context *gin.Context, tokenCookie *TokenCookie) {
	if tokenCookie.RefreshToken != "" {
		service.GenerateRefreshCookies(context, tokenCookie.RefreshToken)
	}

	if tokenCookie.AccessToken != "" {
		service.GenerateAccessCookies(context, tokenCookie.AccessToken)
	}
}
