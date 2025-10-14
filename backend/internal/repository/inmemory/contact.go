package inmemory

import (
	"context"
	"fmt"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactRepository struct{}

func NewContactRepository() repository.ContactRepository {
	return &contactRepository{}
}

func (r *contactRepository) CreateSubmission(ctx context.Context, payload *model.ContactRequest) (*model.ContactSubmission, error) {
	// Use deterministic placeholder while the real persistence layer is prepared.
	id := fmt.Sprintf("submission-%d", time.Now().UnixNano())
	return &model.ContactSubmission{
		ID:      id,
		Status:  "accepted",
		Comment: fmt.Sprintf("queued at %s", time.Now().UTC().Format(time.RFC3339)),
	}, nil
}
