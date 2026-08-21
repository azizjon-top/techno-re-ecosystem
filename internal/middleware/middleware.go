package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"techno-re-ecosystem/internal/logger"
	"time"
)

// RequestLogger middleware logs HTTP requests and responses
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		
		// Log request
		logger.Info("incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// Process request
		c.Next()

		// Log response
		duration := time.Since(startTime)
		logger.Info("request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status_code", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int("response_size", c.Writer.Size()),
		)
	}
}

// ErrorHandler middleware handles panics and logs errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Request.URL.Path),
					zap.Stack("stack"),
				)
				c.JSON(500, gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

// CORS middleware handles cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestTimeout middleware sets request timeout
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := c.Request.Context().WithTimeout(timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ContentTypeJSON middleware enforces JSON content type
func ContentTypeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" {
			c.Header("Content-Type", "application/json")
		}
		c.Next()
	}
}
