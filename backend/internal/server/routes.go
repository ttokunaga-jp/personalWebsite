package server

import (
	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/handler"
)

func registerRoutes(r *gin.Engine, healthHandler *handler.HealthHandler) {
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Ping)
	}
}
