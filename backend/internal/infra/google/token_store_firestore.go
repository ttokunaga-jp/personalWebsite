package google

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	firestoredb "github.com/takumi/personal-website/internal/infra/firestore"
)

type firestoreTokenStore struct {
	client     *firestore.Client
	collection string
	key        []byte
}

// NewFirestoreTokenStore persists OAuth tokens inside Firestore, encrypting
// values at rest using the provided secret.
func NewFirestoreTokenStore(client *firestore.Client, prefix, secret string) (TokenStore, error) {
	if client == nil {
		return nil, fmt.Errorf("firestore token store: client is nil")
	}
	if secret == "" {
		return nil, fmt.Errorf("firestore token store: secret is empty")
	}
	sum := sha256.Sum256([]byte(secret))
	return &firestoreTokenStore{
		client:     client,
		collection: firestoredb.CollectionName(prefix, "google_oauth_tokens"),
		key:        sum[:],
	}, nil
}

type tokenDocument struct {
	AccessToken  string    `firestore:"accessToken"`
	RefreshToken string    `firestore:"refreshToken"`
	Expiry       time.Time `firestore:"expiry"`
	CreatedAt    time.Time `firestore:"createdAt"`
	UpdatedAt    time.Time `firestore:"updatedAt"`
}

func (s *firestoreTokenStore) Load(ctx context.Context, provider string) (*TokenRecord, error) {
	docRef := s.client.Collection(s.collection).Doc(provider)
	snapshot, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("firestore token store: get %s: %w", provider, err)
	}

	var doc tokenDocument
	if err := snapshot.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("firestore token store: decode %s: %w", provider, err)
	}

	access, err := decryptToken(doc.AccessToken, s.key)
	if err != nil {
		return nil, fmt.Errorf("firestore token store: decrypt access token: %w", err)
	}

	refresh := ""
	if doc.RefreshToken != "" {
		refresh, err = decryptToken(doc.RefreshToken, s.key)
		if err != nil {
			return nil, fmt.Errorf("firestore token store: decrypt refresh token: %w", err)
		}
	}

	return &TokenRecord{
		AccessToken:  access,
		RefreshToken: refresh,
		Expiry:       doc.Expiry.UTC(),
	}, nil
}

func (s *firestoreTokenStore) Save(ctx context.Context, provider string, record *TokenRecord) error {
	if record == nil {
		return fmt.Errorf("firestore token store: record is nil")
	}

	encAccess, err := encryptToken(record.AccessToken, s.key)
	if err != nil {
		return fmt.Errorf("firestore token store: encrypt access token: %w", err)
	}

	encRefresh := ""
	if record.RefreshToken != "" {
		encRefresh, err = encryptToken(record.RefreshToken, s.key)
		if err != nil {
			return fmt.Errorf("firestore token store: encrypt refresh token: %w", err)
		}
	}

	docRef := s.client.Collection(s.collection).Doc(provider)
	now := time.Now().UTC()

	return s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := tx.Get(docRef)
		create := false
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return fmt.Errorf("firestore token store: read existing %s: %w", provider, err)
			}
			create = true
		}

		data := map[string]any{
			"accessToken":  encAccess,
			"refreshToken": encRefresh,
			"expiry":       record.Expiry.UTC(),
			"updatedAt":    now,
		}
		if create {
			data["createdAt"] = now
		}

		return tx.Set(docRef, data, firestore.MergeAll)
	})
}

var _ TokenStore = (*firestoreTokenStore)(nil)
