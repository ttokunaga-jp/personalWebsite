package firestore

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type blacklistRepository struct {
	base baseRepository
}

const blacklistCollection = "blacklist"

type blacklistDocument struct {
	ID        int64     `firestore:"id"`
	Email     string    `firestore:"email"`
	Reason    string    `firestore:"reason"`
	CreatedAt time.Time `firestore:"createdAt"`
}

func NewBlacklistRepository(client *firestore.Client, prefix string) repository.BlacklistRepository {
	return &blacklistRepository{base: newBaseRepository(client, prefix)}
}

func (r *blacklistRepository) ListBlacklistEntries(ctx context.Context) ([]model.BlacklistEntry, error) {
	docs, err := r.base.collection(blacklistCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore blacklist: list: %w", err)
	}

	entries, err := r.decodeEntries(docs)
	if err != nil {
		return nil, err
	}

	result := make([]model.BlacklistEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, mapBlacklistDocument(entry))
	}
	return result, nil
}

func (r *blacklistRepository) AddBlacklistEntry(ctx context.Context, entry *model.BlacklistEntry) (*model.BlacklistEntry, error) {
	if entry == nil {
		return nil, repository.ErrInvalidInput
	}

	email := normalizeEmail(entry.Email)
	if email == "" {
		return nil, repository.ErrInvalidInput
	}

	if existing, err := r.FindBlacklistEntryByEmail(ctx, email); err == nil && existing != nil {
		return nil, repository.ErrDuplicate
	} else if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	id, err := nextID(ctx, r.base.client, r.base.prefix, blacklistCollection)
	if err != nil {
		return nil, fmt.Errorf("firestore blacklist: next id: %w", err)
	}

	now := time.Now().UTC()
	docRef := r.base.doc(blacklistCollection, strconv.FormatInt(id, 10))
	payload := blacklistDocument{
		ID:        id,
		Email:     email,
		Reason:    strings.TrimSpace(entry.Reason),
		CreatedAt: now,
	}

	if _, err := docRef.Create(ctx, payload); err != nil {
		return nil, fmt.Errorf("firestore blacklist: create %d: %w", id, err)
	}

	created, err := r.FindBlacklistEntryByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *blacklistRepository) RemoveBlacklistEntry(ctx context.Context, id int64) error {
	docRef := r.base.doc(blacklistCollection, strconv.FormatInt(id, 10))
	if _, err := docRef.Delete(ctx); err != nil {
		if notFound(err) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("firestore blacklist: delete %d: %w", id, err)
	}
	return nil
}

func (r *blacklistRepository) FindBlacklistEntryByEmail(ctx context.Context, email string) (*model.BlacklistEntry, error) {
	normalized := normalizeEmail(email)
	if normalized == "" {
		return nil, repository.ErrInvalidInput
	}

	query := r.base.collection(blacklistCollection).
		Where("email", "==", normalized).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore blacklist: query email %s: %w", normalized, err)
	}
	if len(docs) == 0 {
		return nil, repository.ErrNotFound
	}

	var payload blacklistDocument
	if err := docs[0].DataTo(&payload); err != nil {
		return nil, fmt.Errorf("firestore blacklist: decode email %s: %w", normalized, err)
	}
	entry := mapBlacklistDocument(payload)
	return &entry, nil
}

func (r *blacklistRepository) decodeEntries(docs []*firestore.DocumentSnapshot) ([]blacklistDocument, error) {
	items := make([]blacklistDocument, 0, len(docs))
	for _, doc := range docs {
		var entry blacklistDocument
		if err := doc.DataTo(&entry); err != nil {
			return nil, fmt.Errorf("firestore blacklist: decode %s: %w", doc.Ref.ID, err)
		}
		if entry.ID == 0 {
			if id, err := strconv.ParseInt(doc.Ref.ID, 10, 64); err == nil {
				entry.ID = id
			}
		}
		items = append(items, entry)
	}

	sort.Slice(items, func(i, j int) bool {
		left, right := items[i], items[j]
		if !left.CreatedAt.Equal(right.CreatedAt) {
			return left.CreatedAt.After(right.CreatedAt)
		}
		return left.ID > right.ID
	})

	return items, nil
}

func mapBlacklistDocument(doc blacklistDocument) model.BlacklistEntry {
	return model.BlacklistEntry{
		ID:        doc.ID,
		Email:     doc.Email,
		Reason:    doc.Reason,
		CreatedAt: doc.CreatedAt,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

var _ repository.BlacklistRepository = (*blacklistRepository)(nil)
