package middleware

import (
	"ms1-documents/internal/config"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		config.Logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		if status >= 400 {
			errorMessages := []string{}
			for _, err := range c.Errors {
				errorMessages = append(errorMessages, err.Error())
			}

			config.Logger.Error("HTTP Error",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.String("client_ip", c.ClientIP()),
				zap.Strings("error_messages", errorMessages),
			)
		}
	}
}
