package repository

import (
	"context"

	"github.com/takumi/personal-website/internal/model"
)

// ContentProfileRepository exposes read operations for the redesigned profile aggregate.
type ContentProfileRepository interface {
	GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error)
}

// TechCatalogRepository retrieves canonical technology definitions.
type TechCatalogRepository interface {
	ListTechCatalog(ctx context.Context, includeInactive bool) ([]model.TechCatalogEntry, error)
	GetTechCatalogEntry(ctx context.Context, id uint64) (*model.TechCatalogEntry, error)
	CreateTechCatalogEntry(ctx context.Context, entry *model.TechCatalogEntry) (*model.TechCatalogEntry, error)
	UpdateTechCatalogEntry(ctx context.Context, entry *model.TechCatalogEntry) (*model.TechCatalogEntry, error)
}

// ProjectDocumentRepository retrieves project aggregates compliant with the new schema.
type ProjectDocumentRepository interface {
	ListProjectDocuments(ctx context.Context, includeDrafts bool) ([]model.ProjectDocument, error)
}

// ResearchDocumentRepository retrieves research/blog aggregates.
type ResearchDocumentRepository interface {
	ListResearchDocuments(ctx context.Context, includeDrafts bool) ([]model.ResearchDocument, error)
}

// ContactFormSettingsRepository provides access to contact form configuration v2.
type ContactFormSettingsRepository interface {
	GetContactFormSettings(ctx context.Context) (*model.ContactFormSettingsV2, error)
}

// HomePageConfigRepository returns home page configuration derived from home_page_config tables.
type HomePageConfigRepository interface {
	GetHomePageConfig(ctx context.Context) (*model.HomePageConfigDocument, error)
}
