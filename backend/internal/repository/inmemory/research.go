package inmemory

import (
	"context"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type researchRepository struct {
	research []model.Research
}

func NewResearchRepository() repository.ResearchRepository {
	return &researchRepository{
		research: defaultResearch,
	}
}

func (r *researchRepository) ListResearch(ctx context.Context) ([]model.Research, error) {
	return r.research, nil
}
