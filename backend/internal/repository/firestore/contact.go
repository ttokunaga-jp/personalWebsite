package firestore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactRepository struct {
	base baseRepository
}

const contactCollection = "contact_messages"

type contactDocument struct {
	Name      string    `firestore:"name"`
	Email     string    `firestore:"email"`
	Topic     string    `firestore:"topic"`
	Message   string    `firestore:"message"`
	CreatedAt time.Time `firestore:"createdAt"`
}

// NewContactRepository returns a Firestore-backed implementation for contact submissions.
func NewContactRepository(client *firestore.Client, prefix string) repository.ContactRepository {
	return &contactRepository{
		base: newBaseRepository(client, prefix),
	}
}

func (r *contactRepository) CreateSubmission(ctx context.Context, payload *model.ContactRequest) (*model.ContactSubmission, error) {
	if payload == nil {
		return nil, repository.ErrInvalidInput
	}

	now := time.Now().UTC()
	doc := contactDocument{
		Name:      payload.Name,
		Email:     payload.Email,
		Topic:     payload.Topic,
		Message:   payload.Message,
		CreatedAt: now,
	}

	ref, _, err := r.base.collection(contactCollection).Add(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("firestore contact: add submission: %w", err)
	}

	return &model.ContactSubmission{
		ID:      ref.ID,
		Status:  "queued",
		Comment: fmt.Sprintf("stored at %s", now.Format(time.RFC3339)),
	}, nil
}

var _ repository.ContactRepository = (*contactRepository)(nil)
