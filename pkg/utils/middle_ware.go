package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func StructuredLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)

		// Create a logger for this request with the trace_id
		reqLogger := logger.With(zap.String("trace_id", traceID))
		c.Set("logger", reqLogger)

		c.Next()

		latency := time.Since(start)

		userID, _ := c.Get("user_id")

		reqLogger.Info("Request handled",
			zap.String("http_method", c.Request.Method),
			zap.String("http_path", c.Request.URL.Path),
			zap.Int("http_status_code", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Any("user_id", userID),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

func RequireAuth(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid Authorization format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			logger.Error("Invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info vào context để controller sử dụng
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		// Add user_id to the request logger
		if reqLogger, exists := c.Get("logger"); exists {
			if logger, ok := reqLogger.(*zap.Logger); ok {
				c.Set("logger", logger.With(zap.Int64("user_id", claims.UserID)))
			}
		}

		c.Next()
	}
}