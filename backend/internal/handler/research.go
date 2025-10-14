package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/service"
)

type ResearchHandler struct {
	service service.ResearchService
}

func NewResearchHandler(service service.ResearchService) *ResearchHandler {
	return &ResearchHandler{service: service}
}

func (h *ResearchHandler) ListResearch(c *gin.Context) {
	research, err := h.service.ListResearch(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": research,
	})
}
