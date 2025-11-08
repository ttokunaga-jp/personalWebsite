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
	adminAuthHandler *handler.AdminAuthHandler,
	sessionMiddleware *middleware.AdminSessionMiddleware,
	adminHandler *handler.AdminHandler,
	adminGuard *middleware.AdminGuard,
	securityHandler *handler.SecurityHandler,
) {
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Ping)
		api.HEAD("/health", healthHandler.Ping)
		api.GET("/profile", profileHandler.GetProfile)
		api.GET("/projects", projectHandler.ListProjects)
		api.GET("/research", researchHandler.ListResearch)
		api.GET("/contact/availability", contactHandler.GetAvailability)
		api.GET("/contact/config", contactHandler.GetConfig)
		api.POST("/contact", contactHandler.SubmitContact)
		api.POST("/contact/bookings", bookingHandler.CreateBooking)
		api.GET("/contact/bookings/:lookupHash", bookingHandler.GetReservation)
		api.GET("/auth/login", authHandler.Login)
		api.GET("/auth/callback", authHandler.Callback)
		if securityHandler != nil {
			api.GET("/security/csrf", securityHandler.IssueCSRFToken)
		}
	}

	adminAuth := api.Group("/admin/auth")
	{
		adminAuth.GET("/login", adminAuthHandler.Login)
		adminAuth.GET("/callback", adminAuthHandler.Callback)
		adminAuth.GET("/session", adminAuthHandler.Session)
	}

	publicV1 := api.Group("/v1/public")
	{
		publicV1.GET("/profile", profileHandler.GetProfile)
		publicV1.GET("/projects", projectHandler.ListProjects)
		publicV1.GET("/research", researchHandler.ListResearch)
		publicV1.GET("/contact/availability", contactHandler.GetAvailability)
		publicV1.GET("/contact/config", contactHandler.GetConfig)
		publicV1.POST("/contact", contactHandler.SubmitContact)
		publicV1.POST("/contact/bookings", bookingHandler.CreateBooking)
		publicV1.GET("/contact/bookings/:lookupHash", bookingHandler.GetReservation)
	}

	admin := api.Group("/admin")
	admin.Use(sessionMiddleware.Handler(), adminGuard.RequireAdmin())
	{
		admin.GET("/health", healthHandler.Ping)
		admin.GET("/summary", adminHandler.Summary)
		admin.GET("/tech-catalog", adminHandler.ListTechCatalog)

		admin.GET("/profile", adminHandler.GetProfile)
		admin.PUT("/profile", adminHandler.UpdateProfile)

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

		admin.GET("/contacts", adminHandler.ListContacts)
		admin.GET("/contacts/:id", adminHandler.GetContact)
		admin.PUT("/contacts/:id", adminHandler.UpdateContact)
		admin.DELETE("/contacts/:id", adminHandler.DeleteContact)

		admin.GET("/blacklist", adminHandler.ListBlacklist)
		admin.POST("/blacklist", adminHandler.CreateBlacklist)
		admin.PUT("/blacklist/:id", adminHandler.UpdateBlacklist)
		admin.DELETE("/blacklist/:id", adminHandler.DeleteBlacklist)
	}
}
