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

// NewMeetingRepository selects an appropriate meeting repository implementation based on the Firestore client.
func NewMeetingRepository(db *sqlx.DB, client *firestore.Client, cfg *config.AppConfig) repository.MeetingRepository {
	switch {
	case db != nil:
		return repoMySQL.NewMeetingRepository(db)
	case client != nil:
		return repoFirestore.NewMeetingRepository(client, prefix(cfg))
	default:
		return inmemory.NewMeetingRepository()
	}
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
