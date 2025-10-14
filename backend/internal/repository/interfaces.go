package repository

import (
	"context"

	"github.com/takumi/personal-website/internal/model"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context) (*model.Profile, error)
}

type ProjectRepository interface {
	ListProjects(ctx context.Context) ([]model.Project, error)
}

type ResearchRepository interface {
	ListResearch(ctx context.Context) ([]model.Research, error)
}

type ContactRepository interface {
	CreateSubmission(ctx context.Context, payload *model.ContactRequest) (*model.ContactSubmission, error)
}
