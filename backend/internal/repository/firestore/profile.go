package firestore

import (
	"context"
	"fmt"
	"time"

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
	FocusAreas  []localizedDoc `firestore:"focusAreas"`
	UpdatedAt   time.Time      `firestore:"updatedAt"`
}

type profileRepository struct {
	base baseRepository
}

func NewProfileRepository(client *firestore.Client, prefix string) repository.ProfileRepository {
	return &profileRepository{base: newBaseRepository(client, prefix)}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	doc, err := r.load(ctx)
	if err != nil {
		return nil, err
	}

	profile := &model.Profile{
		Name:        fromLocalizedDoc(doc.Name),
		Title:       fromLocalizedDoc(doc.Title),
		Affiliation: fromLocalizedDoc(doc.Affiliation),
		Lab:         fromLocalizedDoc(doc.Lab),
		Summary:     fromLocalizedDoc(doc.Summary),
		Skills:      make([]model.LocalizedText, 0, len(doc.Skills)),
		FocusAreas:  make([]model.LocalizedText, 0, len(doc.FocusAreas)),
	}

	for _, skill := range doc.Skills {
		profile.Skills = append(profile.Skills, fromLocalizedDoc(skill))
	}
	for _, area := range doc.FocusAreas {
		profile.FocusAreas = append(profile.FocusAreas, fromLocalizedDoc(area))
	}

	return profile, nil
}

func (r *profileRepository) GetAdminProfile(ctx context.Context) (*model.AdminProfile, error) {
	doc, err := r.load(ctx)
	if err != nil {
		return nil, err
	}

	admin := &model.AdminProfile{
		Name:        fromLocalizedDoc(doc.Name),
		Title:       fromLocalizedDoc(doc.Title),
		Affiliation: fromLocalizedDoc(doc.Affiliation),
		Lab:         fromLocalizedDoc(doc.Lab),
		Summary:     fromLocalizedDoc(doc.Summary),
		Skills:      make([]model.LocalizedText, 0, len(doc.Skills)),
		FocusAreas:  make([]model.LocalizedText, 0, len(doc.FocusAreas)),
	}

	if !doc.UpdatedAt.IsZero() {
		t := doc.UpdatedAt.UTC()
		admin.UpdatedAt = &t
	}

	for _, skill := range doc.Skills {
		admin.Skills = append(admin.Skills, fromLocalizedDoc(skill))
	}
	for _, area := range doc.FocusAreas {
		admin.FocusAreas = append(admin.FocusAreas, fromLocalizedDoc(area))
	}

	return admin, nil
}

func (r *profileRepository) UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error) {
	if profile == nil {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(profileCollection, profileDocumentID)

	doc := profileDocument{
		Name:        toLocalizedDoc(profile.Name),
		Title:       toLocalizedDoc(profile.Title),
		Affiliation: toLocalizedDoc(profile.Affiliation),
		Lab:         toLocalizedDoc(profile.Lab),
		Summary:     toLocalizedDoc(profile.Summary),
		Skills:      make([]localizedDoc, 0, len(profile.Skills)),
		FocusAreas:  make([]localizedDoc, 0, len(profile.FocusAreas)),
		UpdatedAt:   time.Now().UTC(),
	}

	for _, skill := range profile.Skills {
		doc.Skills = append(doc.Skills, toLocalizedDoc(skill))
	}
	for _, area := range profile.FocusAreas {
		doc.FocusAreas = append(doc.FocusAreas, toLocalizedDoc(area))
	}

	if _, err := docRef.Set(ctx, doc); err != nil {
		return nil, fmt.Errorf("firestore profile: set %s: %w", profileDocumentID, err)
	}

	return r.GetAdminProfile(ctx)
}

func (r *profileRepository) load(ctx context.Context) (*profileDocument, error) {
	docRef := r.base.doc(profileCollection, profileDocumentID)
	snapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("firestore profile: get %s: %w", profileDocumentID, err)
	}

	var doc profileDocument
	if err := snapshot.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("firestore profile: decode %s: %w", profileDocumentID, err)
	}

	return &doc, nil
}

var _ repository.ProfileRepository = (*profileRepository)(nil)
var _ repository.AdminProfileRepository = (*profileRepository)(nil)
