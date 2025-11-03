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
	includeDrafts := c.Query("includeDrafts") == "true"

	research, err := h.service.ListResearchDocuments(c.Request.Context(), includeDrafts)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": research,
	})
}
