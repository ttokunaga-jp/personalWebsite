package provider

import (
	"cloud.google.com/go/firestore"
	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/repository"
	repoFirestore "github.com/takumi/personal-website/internal/repository/firestore"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	repoMySQL "github.com/takumi/personal-website/internal/repository/mysql"
)

// NewTechCatalogRepository selects the appropriate implementation for the tech catalog.
func NewTechCatalogRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.TechCatalogRepository {
	switch {
	case db != nil:
		return repoMySQL.NewTechCatalogRepository(db)
	case client != nil:
		// TODO: add dedicated Firestore implementation once schema is finalised.
		return inmemory.NewTechCatalogRepository()
	default:
		return inmemory.NewTechCatalogRepository()
	}
}

// NewProfileRepository selects the appropriate implementation with a preference for MySQL.
func NewProfileRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ProfileRepository {
	switch {
	case db != nil:
		return repoMySQL.NewProfileRepository(db)
	case client != nil:
		return repoFirestore.NewProfileRepository(client, prefix(cfg))
	default:
		return inmemory.NewProfileRepository()
	}
}

// NewContentProfileRepository selects the appropriate implementation for the new profile document schema.
func NewContentProfileRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ContentProfileRepository {
	switch {
	case db != nil:
		return repoMySQL.NewContentProfileRepository(db)
	case client != nil:
		return repoFirestore.NewContentProfileRepository(client, prefix(cfg))
	default:
		return inmemory.NewContentProfileRepository()
	}
}

// NewProjectDocumentRepository selects the implementation for project aggregates.
func NewProjectDocumentRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ProjectDocumentRepository {
	switch {
	case db != nil:
		return repoMySQL.NewProjectDocumentRepository(db)
	case client != nil:
		// TODO: add dedicated Firestore implementation.
		return inmemory.NewProjectDocumentRepository()
	default:
		return inmemory.NewProjectDocumentRepository()
	}
}

// NewResearchDocumentRepository selects the implementation for research/blog aggregates.
func NewResearchDocumentRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ResearchDocumentRepository {
	switch {
	case db != nil:
		return repoMySQL.NewResearchDocumentRepository(db)
	case client != nil:
		// TODO: add dedicated Firestore implementation.
		return inmemory.NewResearchDocumentRepository()
	default:
		return inmemory.NewResearchDocumentRepository()
	}
}

// NewContactFormSettingsRepository provides contact settings accessors.
func NewContactFormSettingsRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ContactFormSettingsRepository {
	switch {
	case db != nil:
		return repoMySQL.NewContactFormSettingsRepository(db)
	case client != nil:
		// TODO: add dedicated Firestore implementation.
		return inmemory.NewContactFormSettingsRepository()
	default:
		return inmemory.NewContactFormSettingsRepository()
	}
}

// NewAdminContactFormSettingsRepository exposes administrative operations for contact settings.
func NewAdminContactFormSettingsRepository(repo repository.ContactFormSettingsRepository) repository.AdminContactSettingsRepository {
	if adminRepo, ok := repo.(repository.AdminContactSettingsRepository); ok {
		return adminRepo
	}
	panic("contact form settings repository does not implement admin interface")
}

// NewHomePageConfigRepository loads home page configuration aggregates.
func NewHomePageConfigRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.HomePageConfigRepository {
	switch {
	case db != nil:
		return repoMySQL.NewHomePageConfigRepository(db)
	case client != nil:
		// TODO: add dedicated Firestore implementation.
		return inmemory.NewHomePageConfigRepository()
	default:
		return inmemory.NewHomePageConfigRepository()
	}
}

// NewAdminHomePageConfigRepository exposes administrative operations for home settings.
func NewAdminHomePageConfigRepository(repo repository.HomePageConfigRepository) repository.AdminHomePageConfigRepository {
	if adminRepo, ok := repo.(repository.AdminHomePageConfigRepository); ok {
		return adminRepo
	}
	panic("home page config repository does not implement admin interface")
}

// NewAdminProfileRepository exposes administrative profile capabilities.
func NewAdminProfileRepository(repo repository.ProfileRepository) repository.AdminProfileRepository {
	if adminRepo, ok := repo.(repository.AdminProfileRepository); ok {
		return adminRepo
	}
	panic("profile repository does not implement admin interface")
}

// NewProjectRepository selects an appropriate project repository implementation.
func NewProjectRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ProjectRepository {
	switch {
	case db != nil:
		return repoMySQL.NewProjectRepository(db)
	case client != nil:
		return repoFirestore.NewProjectRepository(client, prefix(cfg))
	default:
		return inmemory.NewProjectRepository()
	}
}

// NewAdminProjectRepository exposes the admin interface while reusing the concrete project repository instance.
func NewAdminProjectRepository(repo repository.ProjectRepository) repository.AdminProjectRepository {
	if adminRepo, ok := repo.(repository.AdminProjectRepository); ok {
		return adminRepo
	}
	panic("project repository does not implement admin interface")
}

// NewAdminSessionRepository constructs the administrative session repository.
func NewAdminSessionRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.AdminSessionRepository {
	switch {
	case db != nil:
		return repoMySQL.NewAdminSessionRepository(db)
	case client != nil:
		return repoFirestore.NewAdminSessionRepository(client, prefix(cfg))
	default:
		return nil
	}
}

// NewResearchRepository selects an appropriate research repository implementation.
func NewResearchRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ResearchRepository {
	switch {
	case db != nil:
		return repoMySQL.NewResearchRepository(db)
	case client != nil:
		return repoFirestore.NewResearchRepository(client, prefix(cfg))
	default:
		return inmemory.NewResearchRepository()
	}
}

// NewAdminResearchRepository exposes the admin interface for the concrete implementation.
func NewAdminResearchRepository(repo repository.ResearchRepository) repository.AdminResearchRepository {
	if adminRepo, ok := repo.(repository.AdminResearchRepository); ok {
		return adminRepo
	}
	panic("research repository does not implement admin interface")
}

// NewContactRepository selects the backing store for contact messages.
func NewContactRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.ContactRepository {
	if db != nil {
		return repoMySQL.NewContactRepository(db)
	}
	if client != nil {
		return repoFirestore.NewContactRepository(client, prefix(cfg))
	}
	return inmemory.NewContactRepository()
}

// NewAdminContactRepository exposes administrative contact moderation capabilities.
func NewAdminContactRepository(repo repository.ContactRepository) repository.AdminContactRepository {
	if adminRepo, ok := repo.(repository.AdminContactRepository); ok {
		return adminRepo
	}
	panic("contact repository does not implement admin interface")
}

// NewAvailabilityRepository selects the appropriate implementation for schedule computation.
func NewAvailabilityRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.AvailabilityRepository {
	switch {
	case db != nil:
		return repoMySQL.NewAvailabilityRepository(db)
	case client != nil:
		return repoFirestore.NewAvailabilityRepository(client, prefix(cfg))
	default:
		return inmemory.NewAvailabilityRepository()
	}
}

// NewBlogRepository selects an appropriate blog repository implementation based on the Firestore client.
func NewBlogRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.BlogRepository {
	switch {
	case db != nil:
		return repoMySQL.NewBlogRepository(db)
	case client != nil:
		return repoFirestore.NewBlogRepository(client, prefix(cfg))
	default:
		return inmemory.NewBlogRepository()
	}
}

// NewMeetingReservationRepository selects an appropriate reservation repository implementation.
func NewMeetingReservationRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.MeetingReservationRepository {
	if db != nil {
		return repoMySQL.NewMeetingReservationRepository(db)
	}
	// Meeting reservations rely on SQL features; fall back to in-memory when a database is unavailable.
	return inmemory.NewMeetingReservationRepository()
}

// NewMeetingNotificationRepository selects an appropriate notification repository implementation.
func NewMeetingNotificationRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.MeetingNotificationRepository {
	if db != nil {
		return repoMySQL.NewMeetingNotificationRepository(db)
	}
	return inmemory.NewMeetingNotificationRepository()
}

// NewBlacklistRepository selects an appropriate blacklist repository implementation based on the Firestore client.
func NewBlacklistRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.BlacklistRepository {
	switch {
	case db != nil:
		return repoMySQL.NewBlacklistRepository(db)
	case client != nil:
		return repoFirestore.NewBlacklistRepository(client, prefix(cfg))
	default:
		return inmemory.NewBlacklistRepository()
	}
}

func prefix(cfg *config.AppConfig) string {
	if cfg == nil {
		return ""
	}
	return cfg.Firestore.CollectionPrefix
}
