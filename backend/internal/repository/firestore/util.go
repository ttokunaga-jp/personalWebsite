package firestore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	firestoredb "github.com/takumi/personal-website/internal/infra/firestore"
	"github.com/takumi/personal-website/internal/model"
)

type baseRepository struct {
	client *firestore.Client
	prefix string
}

func newBaseRepository(client *firestore.Client, prefix string) baseRepository {
	return baseRepository{
		client: client,
		prefix: prefix,
	}
}

func (b baseRepository) collection(name string) *firestore.CollectionRef {
	return b.client.Collection(firestoredb.CollectionName(b.prefix, name))
}

func (b baseRepository) doc(name, id string) *firestore.DocumentRef {
	return b.collection(name).Doc(id)
}

func toLocalizedDoc(text model.LocalizedText) localizedDoc {
	return localizedDoc{
		Ja: text.Ja,
		En: text.En,
	}
}

func fromLocalizedDoc(doc localizedDoc) model.LocalizedText {
	return model.LocalizedText{
		Ja: doc.Ja,
		En: doc.En,
	}
}

type localizedDoc struct {
	Ja string `firestore:"ja"`
	En string `firestore:"en"`
}

func toStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func fromStringPtr(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func toIntPtr(value int) *int {
	return &value
}

func fromIntPtr(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func copyStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	cp := make([]string, len(values))
	copy(cp, values)
	return cp
}

func nextID(ctx context.Context, client *firestore.Client, prefix, key string) (int64, error) {
	col := firestoredb.CollectionName(prefix, "counters")
	docRef := client.Collection(col).Doc(key)

	var id int64
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(docRef)
		switch status.Code(err) {
		case codes.NotFound:
			id = 1
		case codes.OK:
			var payload struct {
				Next int64 `firestore:"next"`
			}
			if err := snap.DataTo(&payload); err != nil {
				return fmt.Errorf("decode counter %s: %w", key, err)
			}
			id = payload.Next + 1
		default:
			if err != nil {
				return fmt.Errorf("read counter %s: %w", key, err)
			}
		}
		return tx.Set(docRef, map[string]any{"next": id}, firestore.MergeAll)
	})

	return id, err
}

func mergeCreatedAt(data map[string]any, create bool, now time.Time) map[string]any {
	if create {
		data["createdAt"] = now
	}
	data["updatedAt"] = now
	return data
}

func notFound(err error) bool {
	return status.Code(err) == codes.NotFound
}
