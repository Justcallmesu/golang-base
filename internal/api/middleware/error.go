package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"justcallmesu.com/golang-base/internal/api/response"
	"justcallmesu.com/golang-base/internal/common/apperror"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.Error

		if !errors.As(err, &appErr) {
			appErr = apperror.Internal(err)
		}

		if appErr.Kind == apperror.KindInternal {
			slog.Error("internal error", "method", c.Request.Method, "path", c.FullPath(), "error", err)
		}

		c.JSON(statusOf(appErr.Kind), response.NewErrorResponse(appErr.Message, appErr.Details))
	}
}

func statusOf(k apperror.Kind) int {
	switch k {
	case apperror.KindBadRequest:
		return http.StatusBadRequest
	case apperror.KindUnauthorized:
		return http.StatusUnauthorized
	case apperror.KindForbidden:
		return http.StatusForbidden
	case apperror.KindNotFound:
		return http.StatusNotFound
	case apperror.KindConflict:
		return http.StatusConflict
	case apperror.KindUnprocessableEntity:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
