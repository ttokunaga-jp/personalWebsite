package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/service"
)

type ContactHandler struct {
	service service.ContactService
}

func NewContactHandler(service service.ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

func (h *ContactHandler) SubmitContact(c *gin.Context) {
	var req model.ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid contact payload", err))
		return
	}

	submission, err := h.service.SubmitContact(c.Request.Context(), &req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": submission,
	})
}
