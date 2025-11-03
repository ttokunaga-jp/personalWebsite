package inmemory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactRepository struct {
	messages map[string]*model.ContactMessage
}

func NewContactRepository() repository.ContactRepository {
	repo := &contactRepository{
		messages: make(map[string]*model.ContactMessage, len(defaultContactMessages)),
	}
	for i := range defaultContactMessages {
		msg := defaultContactMessages[i]
		repo.messages[msg.ID] = cloneContactMessage(&msg)
	}
	return repo
}

func (r *contactRepository) CreateSubmission(ctx context.Context, payload *model.ContactRequest) (*model.ContactSubmission, error) {
	id := fmt.Sprintf("contact-%d", time.Now().UnixNano())
	now := time.Now().UTC()
	message := &model.ContactMessage{
		ID:        id,
		Name:      strings.TrimSpace(payload.Name),
		Email:     strings.TrimSpace(payload.Email),
		Topic:     strings.TrimSpace(payload.Topic),
		Message:   strings.TrimSpace(payload.Message),
		Status:    model.ContactStatusPending,
		AdminNote: "",
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.messages[id] = message
	return &model.ContactSubmission{
		ID:      id,
		Status:  string(model.ContactStatusPending),
		Comment: fmt.Sprintf("queued at %s", now.Format(time.RFC3339)),
	}, nil
}

func (r *contactRepository) ListContactMessages(ctx context.Context) ([]model.ContactMessage, error) {
	messages := make([]model.ContactMessage, 0, len(r.messages))
	for _, msg := range r.messages {
		messages = append(messages, *cloneContactMessage(msg))
	}
	sort.Slice(messages, func(i, j int) bool {
		if messages[i].CreatedAt.Equal(messages[j].CreatedAt) {
			return messages[i].ID > messages[j].ID
		}
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	return messages, nil
}

func (r *contactRepository) GetContactMessage(ctx context.Context, id string) (*model.ContactMessage, error) {
	msg, ok := r.messages[strings.TrimSpace(id)]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneContactMessage(msg), nil
}

func (r *contactRepository) UpdateContactMessage(ctx context.Context, message *model.ContactMessage) (*model.ContactMessage, error) {
	if message == nil {
		return nil, repository.ErrInvalidInput
	}
	id := strings.TrimSpace(message.ID)
	msg, ok := r.messages[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	msg.Topic = strings.TrimSpace(message.Topic)
	msg.Message = strings.TrimSpace(message.Message)
	msg.Status = normalizeContactStatus(message.Status)
	msg.AdminNote = strings.TrimSpace(message.AdminNote)
	msg.UpdatedAt = time.Now().UTC()

	return cloneContactMessage(msg), nil
}

func (r *contactRepository) DeleteContactMessage(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if _, ok := r.messages[id]; !ok {
		return repository.ErrNotFound
	}
	delete(r.messages, id)
	return nil
}

func cloneContactMessage(msg *model.ContactMessage) *model.ContactMessage {
	if msg == nil {
		return nil
	}
	clone := *msg
	return &clone
}

func normalizeContactStatus(status model.ContactStatus) model.ContactStatus {
	switch status {
	case model.ContactStatusPending,
		model.ContactStatusInReview,
		model.ContactStatusResolved,
		model.ContactStatusArchived:
		return status
	default:
		return model.ContactStatusPending
	}
}

var _ repository.ContactRepository = (*contactRepository)(nil)
var _ repository.AdminContactRepository = (*contactRepository)(nil)
