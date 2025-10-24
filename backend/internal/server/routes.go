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
	bookingHandler *handler.BookingHandler,
	authHandler *handler.AuthHandler,
	jwtMiddleware *middleware.JWTMiddleware,
	adminHandler *handler.AdminHandler,
	adminGuard *middleware.AdminGuard,
) {
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Ping)
		api.GET("/profile", profileHandler.GetProfile)
		api.GET("/projects", projectHandler.ListProjects)
		api.GET("/research", researchHandler.ListResearch)
		api.GET("/contact/availability", contactHandler.GetAvailability)
		api.POST("/contact", contactHandler.SubmitContact)
		api.POST("/contact/bookings", bookingHandler.CreateBooking)
		api.GET("/auth/login", authHandler.Login)
		api.GET("/auth/callback", authHandler.Callback)
	}

	admin := api.Group("/admin")
	admin.Use(jwtMiddleware.Handler(), adminGuard.RequireAdmin())
	{
		admin.GET("/health", healthHandler.Ping)
		admin.GET("/summary", adminHandler.Summary)

		admin.GET("/projects", adminHandler.ListProjects)
		admin.POST("/projects", adminHandler.CreateProject)
		admin.GET("/projects/:id", adminHandler.GetProject)
		admin.PUT("/projects/:id", adminHandler.UpdateProject)
		admin.DELETE("/projects/:id", adminHandler.DeleteProject)

		admin.GET("/research", adminHandler.ListResearch)
		admin.POST("/research", adminHandler.CreateResearch)
		admin.GET("/research/:id", adminHandler.GetResearch)
		admin.PUT("/research/:id", adminHandler.UpdateResearch)
		admin.DELETE("/research/:id", adminHandler.DeleteResearch)

		admin.GET("/blogs", adminHandler.ListBlogPosts)
		admin.POST("/blogs", adminHandler.CreateBlogPost)
		admin.GET("/blogs/:id", adminHandler.GetBlogPost)
		admin.PUT("/blogs/:id", adminHandler.UpdateBlogPost)
		admin.DELETE("/blogs/:id", adminHandler.DeleteBlogPost)

		admin.GET("/meetings", adminHandler.ListMeetings)
		admin.POST("/meetings", adminHandler.CreateMeeting)
		admin.GET("/meetings/:id", adminHandler.GetMeeting)
		admin.PUT("/meetings/:id", adminHandler.UpdateMeeting)
		admin.DELETE("/meetings/:id", adminHandler.DeleteMeeting)

		admin.GET("/blacklist", adminHandler.ListBlacklist)
		admin.POST("/blacklist", adminHandler.CreateBlacklist)
		admin.DELETE("/blacklist/:id", adminHandler.DeleteBlacklist)
	}
}
