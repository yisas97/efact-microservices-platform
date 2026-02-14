package handler

import (
	"ms1-documents/internal/utils"
	"ms1-documents/pkg/errors"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, appErr *errors.AppError) {
	utils.RespondWithError(c, appErr)
}
