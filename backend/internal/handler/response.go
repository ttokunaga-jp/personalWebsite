package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
)

func respondError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	appErr := errs.From(err)
	c.JSON(appErr.Status, gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
	})
}
