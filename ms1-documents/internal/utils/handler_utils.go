package utils

import (
	"context"
	"ms1-documents/pkg/errors"
	"time"

	"github.com/gin-gonic/gin"
)

func ManejarErrorServicio(c *gin.Context, err error, mensajeFallback string) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*errors.AppError); ok {
		RespondWithError(c, appErr)
	} else {
		RespondWithError(c, errors.NuevoErrorServidorInterno(mensajeFallback))
	}
	return true
}

func CrearContextoConTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), DefaultOperationTimeout)
}

func CrearContextoConTimeoutPersonalizado(c *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), timeout)
}

func CrearContextoConTimeoutDB(contextoBase context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(contextoBase, DatabaseTimeout)
}

func ValidarJSON(c *gin.Context, objeto interface{}) bool {
	if err := c.ShouldBindJSON(objeto); err != nil {
		RespondWithError(c, errors.NuevoErrorValidacion(ErrorInvalidJSON))
		return true
	}
	return false
}
