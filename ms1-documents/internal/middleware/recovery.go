package middleware

import (
	"ms1-documents/internal/config"
	"ms1-documents/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				config.Logger.Error("PANIC Recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Stack("stacktrace"),
				)

				appErr := errors.ErrorInterno("Ocurrió un error inesperado en el servidor")
				c.JSON(http.StatusInternalServerError, appErr.AJson())
				c.Abort()
			}
		}()
		c.Next()
	}
}
