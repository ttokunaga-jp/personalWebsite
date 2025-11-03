package firestore

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactRepository struct {
	base baseRepository
}

const contactCollection = "contact_messages"

type contactDocument struct {
	Name      string              `firestore:"name"`
	Email     string              `firestore:"email"`
	Topic     string              `firestore:"topic"`
	Message   string              `firestore:"message"`
	Status    model.ContactStatus `firestore:"status"`
	AdminNote string              `firestore:"adminNote"`
	CreatedAt time.Time           `firestore:"createdAt"`
	UpdatedAt time.Time           `firestore:"updatedAt"`
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
		Name:      stringsTrim(payload.Name),
		Email:     stringsTrim(payload.Email),
		Topic:     stringsTrim(payload.Topic),
		Message:   stringsTrim(payload.Message),
		Status:    model.ContactStatusPending,
		AdminNote: "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ref, _, err := r.base.collection(contactCollection).Add(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("firestore contact: add submission: %w", err)
	}

	return &model.ContactSubmission{
		ID:      ref.ID,
		Status:  string(model.ContactStatusPending),
		Comment: fmt.Sprintf("stored at %s", now.Format(time.RFC3339)),
	}, nil
}

func (r *contactRepository) ListContactMessages(ctx context.Context) ([]model.ContactMessage, error) {
	iter := r.base.collection(contactCollection).Documents(ctx)
	defer iter.Stop()

	messages := make([]model.ContactMessage, 0)
	for {
		docSnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("firestore contact: list messages: %w", err)
		}
		msg, err := decodeContactDocument(docSnap.Ref.ID, docSnap)
		if err != nil {
			return nil, err
		}
		messages = append(messages, *msg)
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
	if stringsTrim(id) == "" {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(contactCollection, id)
	snapshot, err := docRef.Get(ctx)
	if err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore contact: get %s: %w", id, err)
	}

	return decodeContactDocument(id, snapshot)
}

func (r *contactRepository) UpdateContactMessage(ctx context.Context, message *model.ContactMessage) (*model.ContactMessage, error) {
	if message == nil {
		return nil, repository.ErrInvalidInput
	}
	id := stringsTrim(message.ID)
	if id == "" {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(contactCollection, id)
	updates := []firestore.Update{
		{Path: "topic", Value: stringsTrim(message.Topic)},
		{Path: "message", Value: stringsTrim(message.Message)},
		{Path: "status", Value: message.Status},
		{Path: "adminNote", Value: stringsTrim(message.AdminNote)},
		{Path: "updatedAt", Value: time.Now().UTC()},
	}
	if _, err := docRef.Update(ctx, updates); err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore contact: update %s: %w", id, err)
	}

	return r.GetContactMessage(ctx, id)
}

func (r *contactRepository) DeleteContactMessage(ctx context.Context, id string) error {
	if stringsTrim(id) == "" {
		return repository.ErrInvalidInput
	}
	docRef := r.base.doc(contactCollection, id)
	if _, err := docRef.Delete(ctx); err != nil {
		if notFound(err) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("firestore contact: delete %s: %w", id, err)
	}
	return nil
}

func decodeContactDocument(id string, snapshot *firestore.DocumentSnapshot) (*model.ContactMessage, error) {
	var doc contactDocument
	if err := snapshot.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("firestore contact: decode %s: %w", id, err)
	}

	createdAt := doc.CreatedAt
	if createdAt.IsZero() {
		createdAt = snapshot.CreateTime
	}
	updatedAt := doc.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = snapshot.UpdateTime
	}

	status := doc.Status
	if status == "" {
		status = model.ContactStatusPending
	}

	message := &model.ContactMessage{
		ID:        id,
		Name:      doc.Name,
		Email:     doc.Email,
		Topic:     doc.Topic,
		Message:   doc.Message,
		Status:    status,
		AdminNote: doc.AdminNote,
		CreatedAt: createdAt.UTC(),
		UpdatedAt: updatedAt.UTC(),
	}
	return message, nil
}

func stringsTrim(value string) string {
	return strings.TrimSpace(value)
}

var _ repository.ContactRepository = (*contactRepository)(nil)
var _ repository.AdminContactRepository = (*contactRepository)(nil)
