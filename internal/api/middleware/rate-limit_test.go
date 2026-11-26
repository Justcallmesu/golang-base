package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func newRateLimitEngine(limit int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ErrorHandler())
	engine.GET("/limited", RateLimit(limit, time.Minute), func(c *gin.Context) {
		c.String(200, "OK")
	})

	return engine
}

func sendFrom(engine *gin.Engine, ip string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)

	request.RemoteAddr = ip + ":1234" // Set the RemoteAddr to simulate the client IP
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, request)

	return recorder
}

func TestRateLimitBlocking(t *testing.T) {
	engine := newRateLimitEngine(2)

	for attempt := 1; attempt <= 2; attempt++ {
		if response := sendFrom(engine, "203.0.113.7"); response.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for attempt %d, got %d", attempt, response.Code)
		}
	}

	response := sendFrom(engine, "203.0.113.7")
	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected status 429 for the third attempt, got %d", response.Code)
	}

	if response.Header().Get("Retry-After") == "" {
		t.Fatalf("Expected Retry-After header to be set, but it was empty")
	}
}

func TestRateLimitIsPerIP(t *testing.T) {
	engine := newRateLimitEngine(1)

	sendFrom(engine, "203.0.113.7")

	if response := sendFrom(engine, "203.0.113.7"); response.Code != http.StatusTooManyRequests {
		t.Fatalf("same IP second request: status %d, want 429", response.Code)
	}

	if response := sendFrom(engine, "198.51.100.9"); response.Code != http.StatusOK {
		t.Fatalf("different IP first request: status %d, want 200", response.Code)
	}
}
