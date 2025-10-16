package provider

import (
	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	"github.com/takumi/personal-website/internal/repository/mysql"
)

// NewProfileRepository selects an appropriate profile repository implementation based on the database handle.
func NewProfileRepository(db *sqlx.DB) repository.ProfileRepository {
	if db == nil {
		return inmemory.NewProfileRepository()
	}
	return mysql.NewProfileRepository(db)
}

// NewProjectRepository selects an appropriate project repository implementation based on the database handle.
func NewProjectRepository(db *sqlx.DB) repository.ProjectRepository {
	if db == nil {
		return inmemory.NewProjectRepository()
	}
	return mysql.NewProjectRepository(db)
}

// NewAdminProjectRepository exposes the admin interface while reusing the concrete project repository instance.
func NewAdminProjectRepository(repo repository.ProjectRepository) repository.AdminProjectRepository {
	if adminRepo, ok := repo.(repository.AdminProjectRepository); ok {
		return adminRepo
	}
	panic("project repository does not implement admin interface")
}

// NewResearchRepository selects an appropriate research repository implementation based on the database handle.
func NewResearchRepository(db *sqlx.DB) repository.ResearchRepository {
	if db == nil {
		return inmemory.NewResearchRepository()
	}
	return mysql.NewResearchRepository(db)
}

// NewAdminResearchRepository exposes the admin interface for the concrete implementation.
func NewAdminResearchRepository(repo repository.ResearchRepository) repository.AdminResearchRepository {
	if adminRepo, ok := repo.(repository.AdminResearchRepository); ok {
		return adminRepo
	}
	panic("research repository does not implement admin interface")
}

// NewContactRepository keeps using the in-memory implementation until persistence is available.
func NewContactRepository() repository.ContactRepository {
	return inmemory.NewContactRepository()
}

// NewAvailabilityRepository selects the appropriate implementation for schedule computation.
func NewAvailabilityRepository(db *sqlx.DB) repository.AvailabilityRepository {
	if db == nil {
		return inmemory.NewAvailabilityRepository()
	}
	return mysql.NewAvailabilityRepository(db)
}

// NewBlogRepository selects an appropriate blog repository implementation based on the database handle.
func NewBlogRepository(db *sqlx.DB) repository.BlogRepository {
	if db == nil {
		return inmemory.NewBlogRepository()
	}
	return mysql.NewBlogRepository(db)
}

// NewMeetingRepository selects an appropriate meeting repository implementation based on the database handle.
func NewMeetingRepository(db *sqlx.DB) repository.MeetingRepository {
	if db == nil {
		return inmemory.NewMeetingRepository()
	}
	return mysql.NewMeetingRepository(db)
}

// NewBlacklistRepository selects an appropriate blacklist repository implementation based on the database handle.
func NewBlacklistRepository(db *sqlx.DB) repository.BlacklistRepository {
	if db == nil {
		return inmemory.NewBlacklistRepository()
	}
	return mysql.NewBlacklistRepository(db)
}
