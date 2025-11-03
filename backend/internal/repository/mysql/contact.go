package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

const (
	insertContactQuery = `
INSERT INTO contact_messages (
	name,
	email,
	topic,
	message,
	status,
	admin_note,
	created_at,
	updated_at
) VALUES (?, ?, ?, ?, 'pending', '', NOW(), NOW())`

	listContactMessagesQuery = `
SELECT
	id,
	name,
	email,
	topic,
	message,
	status,
	admin_note,
	created_at,
	updated_at
FROM contact_messages
ORDER BY created_at DESC, id DESC`

	getContactMessageQuery = `
SELECT
	id,
	name,
	email,
	topic,
	message,
	status,
	admin_note,
	created_at,
	updated_at
FROM contact_messages
WHERE id = ?`

	updateContactMessageQuery = `
UPDATE contact_messages
SET topic = ?,
	message = ?,
	status = ?,
	admin_note = ?,
	updated_at = NOW()
WHERE id = ?`

	deleteContactMessageQuery = `DELETE FROM contact_messages WHERE id = ?`
)

type contactRow struct {
	ID        int64          `db:"id"`
	Name      sql.NullString `db:"name"`
	Email     sql.NullString `db:"email"`
	Topic     sql.NullString `db:"topic"`
	Message   sql.NullString `db:"message"`
	Status    sql.NullString `db:"status"`
	AdminNote sql.NullString `db:"admin_note"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

func (r *contactRepository) CreateSubmission(ctx context.Context, payload *model.ContactRequest) (*model.ContactSubmission, error) {
	if payload == nil {
		return nil, repository.ErrInvalidInput
	}

	email := strings.TrimSpace(payload.Email)
	if email == "" {
		return nil, repository.ErrInvalidInput
	}

	result, err := r.db.ExecContext(ctx, insertContactQuery,
		strings.TrimSpace(payload.Name),
		email,
		strings.TrimSpace(payload.Topic),
		strings.TrimSpace(payload.Message),
	)
	if err != nil {
		return nil, fmt.Errorf("insert contact message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("contact message last insert id: %w", err)
	}

	return &model.ContactSubmission{
		ID:      strconv.FormatInt(id, 10),
		Status:  "pending",
		Comment: fmt.Sprintf("stored at %s", time.Now().UTC().Format(time.RFC3339)),
	}, nil
}

func (r *contactRepository) ListContactMessages(ctx context.Context) ([]model.ContactMessage, error) {
	var rows []contactRow
	if err := r.db.SelectContext(ctx, &rows, listContactMessagesQuery); err != nil {
		return nil, fmt.Errorf("list contact messages: %w", err)
	}

	messages := make([]model.ContactMessage, 0, len(rows))
	for _, row := range rows {
		messages = append(messages, mapContactRow(row))
	}
	return messages, nil
}

func (r *contactRepository) GetContactMessage(ctx context.Context, id string) (*model.ContactMessage, error) {
	internalID, err := parseContactID(id)
	if err != nil {
		return nil, err
	}

	var row contactRow
	if err := r.db.GetContext(ctx, &row, getContactMessageQuery, internalID); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get contact message %s: %w", id, err)
	}
	message := mapContactRow(row)
	return &message, nil
}

func (r *contactRepository) UpdateContactMessage(ctx context.Context, message *model.ContactMessage) (*model.ContactMessage, error) {
	if message == nil {
		return nil, repository.ErrInvalidInput
	}
	internalID, err := parseContactID(message.ID)
	if err != nil {
		return nil, err
	}

	if _, err := r.db.ExecContext(
		ctx,
		updateContactMessageQuery,
		strings.TrimSpace(message.Topic),
		strings.TrimSpace(message.Message),
		strings.TrimSpace(string(message.Status)),
		strings.TrimSpace(message.AdminNote),
		internalID,
	); err != nil {
		return nil, fmt.Errorf("update contact message %s: %w", message.ID, err)
	}

	updated, err := r.GetContactMessage(ctx, message.ID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *contactRepository) DeleteContactMessage(ctx context.Context, id string) error {
	internalID, err := parseContactID(id)
	if err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, deleteContactMessageQuery, internalID)
	if err != nil {
		return fmt.Errorf("delete contact message %s: %w", id, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected delete contact message %s: %w", id, err)
	}
	if affected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func parseContactID(id string) (int64, error) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return 0, repository.ErrInvalidInput
	}
	internalID, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, repository.ErrInvalidInput
	}
	return internalID, nil
}

func mapContactRow(row contactRow) model.ContactMessage {
	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = timeNowUTC()
	}
	updatedAt := row.UpdatedAt.Time
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time
	} else {
		updatedAt = createdAt
	}

	status := model.ContactStatus(strings.TrimSpace(row.Status.String))
	if status == "" {
		status = model.ContactStatusPending
	}

	return model.ContactMessage{
		ID:        strconv.FormatInt(row.ID, 10),
		Name:      nullableString(row.Name),
		Email:     nullableString(row.Email),
		Topic:     nullableString(row.Topic),
		Message:   nullableString(row.Message),
		Status:    status,
		AdminNote: nullableString(row.AdminNote),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

var _ repository.ContactRepository = (*contactRepository)(nil)
var _ repository.AdminContactRepository = (*contactRepository)(nil)
