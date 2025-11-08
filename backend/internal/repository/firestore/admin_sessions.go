package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

const adminSessionsCollection = "admin_sessions"

type adminSessionDocument struct {
	Subject        string     `firestore:"subject"`
	Email          string     `firestore:"email"`
	Roles          []string   `firestore:"roles"`
	UserAgent      *string    `firestore:"userAgent,omitempty"`
	IPAddress      *string    `firestore:"ipAddress,omitempty"`
	ExpiresAt      time.Time  `firestore:"expiresAt"`
	LastAccessedAt time.Time  `firestore:"lastAccessedAt"`
	CreatedAt      time.Time  `firestore:"createdAt"`
	UpdatedAt      time.Time  `firestore:"updatedAt"`
	RevokedAt      *time.Time `firestore:"revokedAt,omitempty"`
}

type adminSessionRepository struct {
	base baseRepository
}

// NewAdminSessionRepository returns a Firestore-backed admin session repository.
func NewAdminSessionRepository(client *firestore.Client, prefix string) repository.AdminSessionRepository {
	return &adminSessionRepository{
		base: newBaseRepository(client, prefix),
	}
}

func (r *adminSessionRepository) doc(hash string) *firestore.DocumentRef {
	return r.base.doc(adminSessionsCollection, hash)
}

func (r *adminSessionRepository) CreateSession(ctx context.Context, session *model.AdminSession) (*model.AdminSession, error) {
	if session == nil {
		return nil, fmt.Errorf("firestore admin session: %w", repository.ErrInvalidInput)
	}

	doc := adminSessionDocument{
		Subject:        session.Subject,
		Email:          session.Email,
		Roles:          append([]string(nil), session.Roles...),
		UserAgent:      toStringPtr(session.UserAgent),
		IPAddress:      toStringPtr(session.IPAddress),
		ExpiresAt:      session.ExpiresAt.UTC(),
		LastAccessedAt: session.LastAccessedAt.UTC(),
		CreatedAt:      session.CreatedAt.UTC(),
		UpdatedAt:      session.UpdatedAt.UTC(),
		RevokedAt:      normalizeTimePtr(session.RevokedAt),
	}

	if _, err := r.doc(session.TokenHash).Create(ctx, doc); err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, fmt.Errorf("firestore admin session: duplicate session hash")
		}
		return nil, fmt.Errorf("firestore admin session: create %s: %w", session.TokenHash, err)
	}

	return r.FindSessionByHash(ctx, session.TokenHash)
}

func (r *adminSessionRepository) FindSessionByHash(ctx context.Context, hash string) (*model.AdminSession, error) {
	snap, err := r.doc(hash).Get(ctx)
	switch status.Code(err) {
	case codes.NotFound:
		return nil, repository.ErrNotFound
	case codes.OK:
		// continue
	default:
		if err != nil {
			return nil, fmt.Errorf("firestore admin session: get %s: %w", hash, err)
		}
	}

	var doc adminSessionDocument
	if err := snap.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("firestore admin session: decode %s: %w", hash, err)
	}

	return mapAdminSessionDoc(hash, doc), nil
}

func (r *adminSessionRepository) UpdateSessionActivity(ctx context.Context, hash string, lastAccessed time.Time, expiresAt time.Time) (*model.AdminSession, error) {
	ref := r.doc(hash)
	var updated *model.AdminSession

	err := r.base.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(ref)
		switch status.Code(err) {
		case codes.NotFound:
			return repository.ErrNotFound
		case codes.OK:
			// continue
		default:
			if err != nil {
				return fmt.Errorf("firestore admin session: get %s: %w", hash, err)
			}
		}

		var doc adminSessionDocument
		if err := snap.DataTo(&doc); err != nil {
			return fmt.Errorf("firestore admin session: decode %s: %w", hash, err)
		}
		if doc.RevokedAt != nil {
			return repository.ErrNotFound
		}

		now := time.Now().UTC()
		doc.LastAccessedAt = lastAccessed.UTC()
		doc.ExpiresAt = expiresAt.UTC()
		doc.UpdatedAt = now

		update := map[string]any{
			"lastAccessedAt": doc.LastAccessedAt,
			"expiresAt":      doc.ExpiresAt,
			"updatedAt":      doc.UpdatedAt,
		}
		if err := tx.Set(ref, update, firestore.MergeAll); err != nil {
			return fmt.Errorf("firestore admin session: update %s: %w", hash, err)
		}

		updated = mapAdminSessionDoc(hash, doc)
		return nil
	})

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return updated, nil
}

func (r *adminSessionRepository) RevokeSession(ctx context.Context, hash string) error {
	ref := r.doc(hash)

	err := r.base.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(ref)
		switch status.Code(err) {
		case codes.NotFound:
			return repository.ErrNotFound
		case codes.OK:
			// continue
		default:
			if err != nil {
				return fmt.Errorf("firestore admin session: get %s: %w", hash, err)
			}
		}

		var doc adminSessionDocument
		if err := snap.DataTo(&doc); err != nil {
			return fmt.Errorf("firestore admin session: decode %s: %w", hash, err)
		}
		if doc.RevokedAt != nil {
			return repository.ErrNotFound
		}

		now := time.Now().UTC()
		update := map[string]any{
			"revokedAt": now,
			"updatedAt": now,
		}

		if err := tx.Set(ref, update, firestore.MergeAll); err != nil {
			return fmt.Errorf("firestore admin session: revoke %s: %w", hash, err)
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	return nil
}

func mapAdminSessionDoc(hash string, doc adminSessionDocument) *model.AdminSession {
	session := &model.AdminSession{
		TokenHash:      hash,
		Subject:        doc.Subject,
		Email:          doc.Email,
		Roles:          copyStringSlice(doc.Roles),
		UserAgent:      fromStringPtr(doc.UserAgent),
		IPAddress:      fromStringPtr(doc.IPAddress),
		ExpiresAt:      doc.ExpiresAt,
		LastAccessedAt: doc.LastAccessedAt,
		CreatedAt:      doc.CreatedAt,
		UpdatedAt:      doc.UpdatedAt,
	}

	if doc.RevokedAt != nil {
		value := doc.RevokedAt.UTC()
		session.RevokedAt = &value
	}

	return session
}

func normalizeTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	val := t.UTC()
	return &val
}

var _ repository.AdminSessionRepository = (*adminSessionRepository)(nil)
