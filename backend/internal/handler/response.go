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
	response := gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
	}
	if requestID := c.Writer.Header().Get("X-Request-ID"); requestID != "" {
		response["request_id"] = requestID
	}
	c.JSON(appErr.Status, response)
}
