package provider

import (
	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/repository"
	repoFirestore "github.com/takumi/personal-website/internal/repository/firestore"
	"github.com/takumi/personal-website/internal/repository/inmemory"
)

// NewProfileRepository selects an appropriate profile repository implementation based on the Firestore client.
func NewProfileRepository(client *firestore.Client, cfg *config.AppConfig) repository.ProfileRepository {
	if client == nil {
		return inmemory.NewProfileRepository()
	}
	return repoFirestore.NewProfileRepository(client, prefix(cfg))
}

// NewProjectRepository selects an appropriate project repository implementation based on the Firestore client.
func NewProjectRepository(client *firestore.Client, cfg *config.AppConfig) repository.ProjectRepository {
	if client == nil {
		return inmemory.NewProjectRepository()
	}
	return repoFirestore.NewProjectRepository(client, prefix(cfg))
}

// NewAdminProjectRepository exposes the admin interface while reusing the concrete project repository instance.
func NewAdminProjectRepository(repo repository.ProjectRepository) repository.AdminProjectRepository {
	if adminRepo, ok := repo.(repository.AdminProjectRepository); ok {
		return adminRepo
	}
	panic("project repository does not implement admin interface")
}

// NewResearchRepository selects an appropriate research repository implementation based on the Firestore client.
func NewResearchRepository(client *firestore.Client, cfg *config.AppConfig) repository.ResearchRepository {
	if client == nil {
		return inmemory.NewResearchRepository()
	}
	return repoFirestore.NewResearchRepository(client, prefix(cfg))
}

// NewAdminResearchRepository exposes the admin interface for the concrete implementation.
func NewAdminResearchRepository(repo repository.ResearchRepository) repository.AdminResearchRepository {
	if adminRepo, ok := repo.(repository.AdminResearchRepository); ok {
		return adminRepo
	}
	panic("research repository does not implement admin interface")
}

// NewContactRepository keeps using the in-memory implementation until Firestore persistence is introduced.
func NewContactRepository() repository.ContactRepository {
	return inmemory.NewContactRepository()
}

// NewAvailabilityRepository selects the appropriate implementation for schedule computation.
func NewAvailabilityRepository(client *firestore.Client, cfg *config.AppConfig) repository.AvailabilityRepository {
	if client == nil {
		return inmemory.NewAvailabilityRepository()
	}
	return repoFirestore.NewAvailabilityRepository(client, prefix(cfg))
}

// NewBlogRepository selects an appropriate blog repository implementation based on the Firestore client.
func NewBlogRepository(client *firestore.Client, cfg *config.AppConfig) repository.BlogRepository {
	if client == nil {
		return inmemory.NewBlogRepository()
	}
	return repoFirestore.NewBlogRepository(client, prefix(cfg))
}

// NewMeetingRepository selects an appropriate meeting repository implementation based on the Firestore client.
func NewMeetingRepository(client *firestore.Client, cfg *config.AppConfig) repository.MeetingRepository {
	if client == nil {
		return inmemory.NewMeetingRepository()
	}
	return repoFirestore.NewMeetingRepository(client, prefix(cfg))
}

// NewBlacklistRepository selects an appropriate blacklist repository implementation based on the Firestore client.
func NewBlacklistRepository(client *firestore.Client, cfg *config.AppConfig) repository.BlacklistRepository {
	if client == nil {
		return inmemory.NewBlacklistRepository()
	}
	return repoFirestore.NewBlacklistRepository(client, prefix(cfg))
}

func prefix(cfg *config.AppConfig) string {
	if cfg == nil {
		return ""
	}
	return cfg.Firestore.CollectionPrefix
}
