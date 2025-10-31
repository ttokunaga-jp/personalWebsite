package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

const profileCollection = "profiles"
const profileDocumentID = "primary"

type profileDocument struct {
	Name        localizedDoc   `firestore:"name"`
	Title       localizedDoc   `firestore:"title"`
	Affiliation localizedDoc   `firestore:"affiliation"`
	Lab         localizedDoc   `firestore:"lab"`
	Summary     localizedDoc   `firestore:"summary"`
	Skills      []localizedDoc `firestore:"skills"`
}

type profileRepository struct {
	base baseRepository
}

func NewProfileRepository(client *firestore.Client, prefix string) repository.ProfileRepository {
	return &profileRepository{base: newBaseRepository(client, prefix)}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	docRef := r.base.doc(profileCollection, profileDocumentID)
	snapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("firestore profile: get %s: %w", profileDocumentID, err)
	}

	var doc profileDocument
	if err := snapshot.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("firestore profile: decode %s: %w", profileDocumentID, err)
	}

	profile := &model.Profile{
		Name:        fromLocalizedDoc(doc.Name),
		Title:       fromLocalizedDoc(doc.Title),
		Affiliation: fromLocalizedDoc(doc.Affiliation),
		Lab:         fromLocalizedDoc(doc.Lab),
		Summary:     fromLocalizedDoc(doc.Summary),
		Skills:      make([]model.LocalizedText, 0, len(doc.Skills)),
	}

	for _, skill := range doc.Skills {
		profile.Skills = append(profile.Skills, fromLocalizedDoc(skill))
	}

	return profile, nil
}

var _ repository.ProfileRepository = (*profileRepository)(nil)
