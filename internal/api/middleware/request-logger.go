package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/url"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"justcallmesu.com/golang-base/internal/common/utils"
)

var sensitiveKeys = regexp.MustCompile(`(?i)password|token|secret|authorization|cookie`)

const redactedValue = "[REDACTED]"

type LoggerConfig struct {
	LogBody      bool                // False in Prod
	MaxBodyBytes int64               // e.g 4 Kb
	SkipPaths    map[string]struct{} // Paths to skip logging
	Pretty       bool                // Pretty print JSON logs
}

func RequestLogger(cfg LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, skip := cfg.SkipPaths[c.Request.URL.Path]; skip {
			c.Next()
			return
		}

		startTime := time.Now()
		requestID := c.GetHeader("X-Request-ID")

		requestID = uuid.NewString()

		c.Set(utils.RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		var body any
		if cfg.LogBody && c.Request.Body != nil && c.ContentType() == "application/json" {
			rawBody, _ := io.ReadAll(io.LimitReader(c.Request.Body, cfg.MaxBodyBytes))
			c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(rawBody), c.Request.Body))
			body = redactJSON(rawBody)
		}

		c.Next()

		// After Processing Request Sections
		status := c.Writer.Status()
		entry := requestLog{
			RequestID: requestID,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Route:     c.FullPath(),
			Status:    status,
			IP:        c.ClientIP(),
			Duration:  time.Since(startTime),
			Username:  "",
			UserID:    "",
			Error:     "",
			Query:     map[string]any{},
			Body:      body,
		}

		if claims, err := utils.GetUserData(c); err == nil {
			entry.UserID = claims.UserID.String()
			entry.Username = claims.Username
		} else {
			entry.UserID = "N/A"
			entry.Username = "GUEST"
		}

		if query := c.Request.URL.Query(); len(query) > 0 {
			entry.Query = redactQuery(query)
		}

		if body != nil {
			entry.Body = body
		}

		if len(c.Errors) > 0 {
			entry.Error = c.Errors.Last().Error()
		}

		if cfg.Pretty {
			printPretty(entry)
			return
		}

		logLevel := slog.LevelInfo

		switch {
		case status >= 500:
			logLevel = slog.LevelError
		case status >= 400:
			logLevel = slog.LevelWarn
		}

		attributes := entry.attributes()
		slog.Log(c.Request.Context(), logLevel, "http request", attributes...)
	}
}

// Helper Utils

func redactJSON(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}

	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return "[INVALID JSON]"
	}

	return redactValue(data)
}

func redactQuery(q url.Values) map[string]any {
	output := make(map[string]any, len(q))

	for k, v := range q {
		switch {
		case sensitiveKeys.MatchString(k):
			output[k] = redactedValue // Redact sensitive keys
		case len(v) == 1:
			output[k] = redactValue(v[0]) // Redact single value
		default:
			output[k] = redactValue(v) // Redact multiple values (array)
		}
	}
	return output
}

func redactValue(value any) any {
	switch v := value.(type) {
	case map[string]any: // Handle maps / objects
		for key, val := range v {
			if sensitiveKeys.MatchString(key) {
				v[key] = redactedValue
			} else {
				v[key] = redactValue(val)
			}
		}
		return v
	case []any: // Handle arrays/slices
		for i, val := range v {
			v[i] = redactValue(val)
		}
		return v
	default: // Primitive types (string, number, bool, etc.)
		return v
	}
}
