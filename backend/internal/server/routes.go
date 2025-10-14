package server

import (
	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/middleware"
)

func registerRoutes(
	r *gin.Engine,
	healthHandler *handler.HealthHandler,
	profileHandler *handler.ProfileHandler,
	projectHandler *handler.ProjectHandler,
	researchHandler *handler.ResearchHandler,
	contactHandler *handler.ContactHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Ping)
		api.GET("/profile", profileHandler.GetProfile)
		api.GET("/projects", projectHandler.ListProjects)
		api.GET("/research", researchHandler.ListResearch)
		api.POST("/contact", contactHandler.SubmitContact)
	}

	admin := api.Group("/admin")
	admin.Use(jwtMiddleware.Handler())
	{
		admin.GET("/health", healthHandler.Ping)
	}
}
