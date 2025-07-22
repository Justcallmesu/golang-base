package cookies

import (
	"time"

	"github.com/gin-gonic/gin"
)

func SetCookie(context *gin.Context, cookieName string, cookiePayload string, cookieMaxAge int64) error {
	context.SetCookie(cookieName, cookiePayload, int(cookieMaxAge*int64(time.Hour/time.Second)), "/", "", false, true)
	return nil
}
