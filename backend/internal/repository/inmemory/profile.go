package inmemory

import (
	"context"
	"strings"
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
	return documentToLegacy(r.profile), nil
}

func (r *profileRepository) GetAdminProfile(ctx context.Context) (*model.AdminProfile, error) {
	return cloneAdminProfile(r.profile), nil
}

func (r *profileRepository) UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error) {
	if profile == nil {
		return nil, repository.ErrInvalidInput
	}

	clone := cloneAdminProfile(profile)
	clone.UpdatedAt = time.Now().UTC()
	r.profile = clone

	return cloneAdminProfile(r.profile), nil
}

func cloneAdminProfile(src *model.AdminProfile) *model.AdminProfile {
	if src == nil {
		return nil
	}

	clone := *src
	clone.Affiliations = append([]model.ProfileAffiliation(nil), src.Affiliations...)
	clone.Communities = append([]model.ProfileAffiliation(nil), src.Communities...)
	clone.WorkHistory = append([]model.ProfileWorkExperience(nil), src.WorkHistory...)
	clone.TechSections = cloneTechSections(src.TechSections)
	clone.SocialLinks = append([]model.ProfileSocialLink(nil), src.SocialLinks...)

	if src.Home != nil {
		homeClone := *src.Home
		homeClone.QuickLinks = append([]model.HomeQuickLink(nil), src.Home.QuickLinks...)
		homeClone.ChipSources = append([]model.HomeChipSource(nil), src.Home.ChipSources...)
		clone.Home = &homeClone
	}

	return &clone
}

func cloneTechSections(sections []model.ProfileTechSection) []model.ProfileTechSection {
	if len(sections) == 0 {
		return nil
	}
	out := make([]model.ProfileTechSection, len(sections))
	for i, section := range sections {
		out[i] = section
		out[i].Members = append([]model.TechMembership(nil), section.Members...)
	}
	return out
}

func documentToLegacy(doc *model.AdminProfile) *model.Profile {
	if doc == nil {
		return nil
	}

	name := model.LocalizedText{
		Ja: strings.TrimSpace(doc.DisplayName),
		En: strings.TrimSpace(doc.DisplayName),
	}
	title := doc.Headline
	summary := doc.Summary

	var affiliation model.LocalizedText
	if len(doc.Affiliations) > 0 {
		primary := doc.Affiliations[0]
		affiliation = model.NewLocalizedText(primary.Name, primary.Name)
	}

	lab := doc.Lab.Name

	return &model.Profile{
		Name:        name,
		Title:       title,
		Affiliation: affiliation,
		Lab:         lab,
		Summary:     summary,
	}
}

var _ repository.ProfileRepository = (*profileRepository)(nil)
var _ repository.AdminProfileRepository = (*profileRepository)(nil)
