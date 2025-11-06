package utils

import (
	"homemie/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := uuid.New().String()

		// Set trace_id in the request context
		c.Request = c.Request.WithContext(logger.ContextWithTraceID(c.Request.Context(), traceID))

		c.Next()

		latency := time.Since(start)

		userID, _ := c.Get("user_id")

		logger.Infow(c.Request.Context(), "Request handled",
			"http_method", c.Request.Method,
			"http_path", c.Request.URL.Path,
			"http_status_code", c.Writer.Status(),
			"latency", latency,
			"user_id", userID,
			"client_ip", c.ClientIP(),
		)
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warnw(c.Request.Context(), "Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warnw(c.Request.Context(), "Invalid Authorization format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			logger.Errorw(c.Request.Context(), "Invalid token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info vào context để controller sử dụng
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		// Add user_id to the request logger in context
		c.Request = c.Request.WithContext(logger.ContextWithFields(c.Request.Context(), zap.Int64("user_id", claims.UserID)))

		c.Next()
	}
}