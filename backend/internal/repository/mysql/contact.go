package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactRepository struct {
	db *sqlx.DB
}

// NewContactRepository persists contact submissions to the contact_messages table.
func NewContactRepository(db *sqlx.DB) repository.ContactRepository {
	return &contactRepository{db: db}
}

const insertContactQuery = `
INSERT INTO contact_messages (
	name,
	email,
	topic,
	message,
	created_at
) VALUES (?, ?, ?, ?, NOW())`

func (r *contactRepository) CreateSubmission(ctx context.Context, payload *model.ContactRequest) (*model.ContactSubmission, error) {
	if payload == nil {
		return nil, repository.ErrInvalidInput
	}

	result, err := r.db.ExecContext(ctx, insertContactQuery,
		payload.Name,
		payload.Email,
		payload.Topic,
		payload.Message,
	)
	if err != nil {
		return nil, fmt.Errorf("insert contact message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("contact message last insert id: %w", err)
	}

	return &model.ContactSubmission{
		ID:      fmt.Sprintf("%d", id),
		Status:  "queued",
		Comment: fmt.Sprintf("stored at %s", time.Now().UTC().Format(time.RFC3339)),
	}, nil
}
