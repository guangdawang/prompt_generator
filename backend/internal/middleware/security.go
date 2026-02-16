package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultAllowedOrigins = "*"

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}

func CORSFromEnv() gin.HandlerFunc {
	allowedOrigins := parseListEnv("ALLOWED_ORIGINS", defaultAllowedOrigins)
	allowCredentials := strings.EqualFold(os.Getenv("ALLOW_CREDENTIALS"), "true")
	allowedMethods := "GET, POST, PUT, DELETE, OPTIONS"
	allowedHeaders := "Content-Type, Authorization, Accept, Origin, Cache-Control, X-Requested-With"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			if allowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			c.Header("Access-Control-Allow-Headers", allowedHeaders)
			c.Header("Access-Control-Allow-Methods", allowedMethods)
		}

		if c.Request.Method == http.MethodOptions {
			if origin != "" && !isOriginAllowed(origin, allowedOrigins) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RequestSizeLimitFromEnv() gin.HandlerFunc {
	maxBytes := parseInt64Env("MAX_BODY_BYTES", 1<<20)
	return RequestSizeLimit(maxBytes)
}

func RequestSizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request too large"})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func RateLimitFromEnv() gin.HandlerFunc {
	limit := parseIntEnv("RATE_LIMIT_PER_MINUTE", 60)
	if limit <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	limiter := newRateLimiter(limit, time.Minute)
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !limiter.Allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

type rateLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	items  map[string]*rateEntry
}

type rateEntry struct {
	count int
	reset time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:  limit,
		window: window,
		items:  make(map[string]*rateEntry),
	}
}

func (r *rateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	entry, ok := r.items[key]
	if !ok || now.After(entry.reset) {
		r.items[key] = &rateEntry{count: 1, reset: now.Add(r.window)}
		return true
	}

	if entry.count >= r.limit {
		return false
	}

	entry.count++
	return true
}

func isOriginAllowed(origin string, allowed []string) bool {
	for _, allowedOrigin := range allowed {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
	}
	return false
}

func parseListEnv(key, defaultValue string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		value = defaultValue
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{defaultValue}
	}
	return result
}

func parseIntEnv(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func parseInt64Env(key string, defaultValue int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return parsed
}
