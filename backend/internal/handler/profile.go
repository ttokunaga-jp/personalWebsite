package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/service"
)

type ProfileHandler struct {
	service service.ProfileService
}

func NewProfileHandler(service service.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	profile, err := h.service.GetProfile(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}
