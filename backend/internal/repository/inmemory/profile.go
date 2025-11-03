package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type profileRepository struct {
	profile *model.AdminProfile
}

func NewProfileRepository() repository.ProfileRepository {
	return &profileRepository{
		profile: cloneAdminProfile(defaultAdminProfile),
	}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	return adminProfileToPublic(r.profile), nil
}

func (r *profileRepository) GetAdminProfile(ctx context.Context) (*model.AdminProfile, error) {
	return cloneAdminProfile(r.profile), nil
}

func (r *profileRepository) UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error) {
	if profile == nil {
		return nil, repository.ErrInvalidInput
	}
	updated := cloneAdminProfile(profile)
	now := time.Now().UTC()
	updated.UpdatedAt = &now
	r.profile = updated
	return cloneAdminProfile(r.profile), nil
}

func cloneAdminProfile(src *model.AdminProfile) *model.AdminProfile {
	if src == nil {
		return nil
	}
	clone := &model.AdminProfile{
		Name:        src.Name,
		Title:       src.Title,
		Affiliation: src.Affiliation,
		Lab:         src.Lab,
		Summary:     src.Summary,
		Skills:      append([]model.LocalizedText(nil), src.Skills...),
		FocusAreas:  append([]model.LocalizedText(nil), src.FocusAreas...),
		UpdatedAt:   src.UpdatedAt,
	}
	return clone
}

func adminProfileToPublic(admin *model.AdminProfile) *model.Profile {
	if admin == nil {
		return nil
	}
	return &model.Profile{
		Name:        admin.Name,
		Title:       admin.Title,
		Affiliation: admin.Affiliation,
		Lab:         admin.Lab,
		Summary:     admin.Summary,
		Skills:      append([]model.LocalizedText(nil), admin.Skills...),
		FocusAreas:  append([]model.LocalizedText(nil), admin.FocusAreas...),
	}
}

var _ repository.ProfileRepository = (*profileRepository)(nil)
var _ repository.AdminProfileRepository = (*profileRepository)(nil)
