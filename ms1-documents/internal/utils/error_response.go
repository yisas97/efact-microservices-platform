package utils

import (
	"ms1-documents/pkg/errors"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, appErr *errors.AppError) {
	c.Error(appErr)

	c.JSON(appErr.Code, appErr.AJson())
}
