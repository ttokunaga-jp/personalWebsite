package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type profileRepository struct {
	db      *sqlx.DB
	content repository.ContentProfileRepository
}

// NewProfileRepository returns a repository backed by MySQL. Caller must ensure db is non-nil.
func NewProfileRepository(db *sqlx.DB) repository.ProfileRepository {
	return &profileRepository{
		db:      db,
		content: NewContentProfileRepository(db),
	}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	document, err := r.content.GetProfileDocument(ctx)
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, repository.ErrNotFound
	}
	return convertDocumentToLegacy(document), nil
}

func convertDocumentToLegacy(doc *model.ProfileDocument) *model.Profile {
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
		affiliation = model.LocalizedText{
			Ja: doc.Affiliations[0].Name,
			En: doc.Affiliations[0].Name,
		}
	}

	lab := doc.Lab.Name

	return &model.Profile{
		Name:        name,
		Title:       title,
		Affiliation: affiliation,
		Lab:         lab,
		Summary:     summary,
		Skills:      nil,
		FocusAreas:  nil,
	}
}

func (r *profileRepository) GetAdminProfile(ctx context.Context) (*model.AdminProfile, error) {
	document, err := r.content.GetProfileDocument(ctx)
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, repository.ErrNotFound
	}
	return document, nil
}

func (r *profileRepository) UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error) {
	if profile == nil {
		return nil, repository.ErrInvalidInput
	}

	existing, err := r.content.GetProfileDocument(ctx)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, repository.ErrNotFound
	}

	profileID := existing.ID

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("profile update: begin tx: %w", err)
	}

	var opErr error
	defer rollbackOnError(tx, &opErr)

	if err := r.updateProfileRow(ctx, tx, profileID, profile); err != nil {
		opErr = err
		return nil, err
	}

	if err := r.replaceAffiliations(ctx, tx, profileID, model.ProfileAffiliationKindAffiliation, profile.Affiliations); err != nil {
		opErr = err
		return nil, err
	}

	if err := r.replaceAffiliations(ctx, tx, profileID, model.ProfileAffiliationKindCommunity, profile.Communities); err != nil {
		opErr = err
		return nil, err
	}

	if err := r.replaceWorkHistory(ctx, tx, profileID, profile.WorkHistory); err != nil {
		opErr = err
		return nil, err
	}

	if err := r.replaceSocialLinks(ctx, tx, profileID, profile.SocialLinks); err != nil {
		opErr = err
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		opErr = fmt.Errorf("profile update: commit: %w", err)
		return nil, opErr
	}
	opErr = nil

	updated, err := r.content.GetProfileDocument(ctx)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, repository.ErrNotFound
	}
	return updated, nil
}

func (r *profileRepository) updateProfileRow(ctx context.Context, tx *sqlx.Tx, profileID uint64, profile *model.AdminProfile) error {
	const query = `
UPDATE profiles SET
    display_name = ?,
    headline_ja = ?,
    headline_en = ?,
    summary_ja = ?,
    summary_en = ?,
    avatar_url = ?,
    location_ja = ?,
    location_en = ?,
    theme_mode = ?,
    theme_accent_color = ?,
    lab_name_ja = ?,
    lab_name_en = ?,
    lab_advisor_ja = ?,
    lab_advisor_en = ?,
    lab_room_ja = ?,
    lab_room_en = ?,
    lab_url = ?,
    updated_at = NOW(3)
WHERE id = ?`

	mode := strings.TrimSpace(string(profile.Theme.Mode))
	if mode == "" {
		mode = string(model.ProfileThemeModeSystem)
	}

	if _, err := tx.ExecContext(
		ctx,
		query,
		strings.TrimSpace(profile.DisplayName),
		strings.TrimSpace(profile.Headline.Ja),
		strings.TrimSpace(profile.Headline.En),
		strings.TrimSpace(profile.Summary.Ja),
		strings.TrimSpace(profile.Summary.En),
		strings.TrimSpace(profile.AvatarURL),
		strings.TrimSpace(profile.Location.Ja),
		strings.TrimSpace(profile.Location.En),
		mode,
		nullString(profile.Theme.AccentColor),
		strings.TrimSpace(profile.Lab.Name.Ja),
		strings.TrimSpace(profile.Lab.Name.En),
		strings.TrimSpace(profile.Lab.Advisor.Ja),
		strings.TrimSpace(profile.Lab.Advisor.En),
		strings.TrimSpace(profile.Lab.Room.Ja),
		strings.TrimSpace(profile.Lab.Room.En),
		strings.TrimSpace(profile.Lab.URL),
		profileID,
	); err != nil {
		return fmt.Errorf("profile update: update profiles %d: %w", profileID, err)
	}
	return nil
}

func (r *profileRepository) replaceAffiliations(
	ctx context.Context,
	tx *sqlx.Tx,
	profileID uint64,
	kind model.ProfileAffiliationKind,
	records []model.ProfileAffiliation,
) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM profile_affiliations WHERE profile_id = ? AND kind = ?`, profileID, string(kind)); err != nil {
		return fmt.Errorf("profile update: delete affiliations %s: %w", kind, err)
	}

	if len(records) == 0 {
		return nil
	}

	const query = `
INSERT INTO profile_affiliations (
    profile_id,
    kind,
    name,
    url,
    started_at,
    description_ja,
    description_en,
    sort_order
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	for _, record := range records {
		if record.Kind != kind {
			record.Kind = kind
		}
		if _, err := tx.ExecContext(
			ctx,
			query,
			profileID,
			string(kind),
			strings.TrimSpace(record.Name),
			nullString(record.URL),
			record.StartedAt.UTC(),
			strings.TrimSpace(record.Description.Ja),
			strings.TrimSpace(record.Description.En),
			record.SortOrder,
		); err != nil {
			return fmt.Errorf("profile update: insert affiliation %s: %w", kind, err)
		}
	}
	return nil
}

func (r *profileRepository) replaceWorkHistory(ctx context.Context, tx *sqlx.Tx, profileID uint64, history []model.ProfileWorkExperience) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM profile_work_history WHERE profile_id = ?`, profileID); err != nil {
		return fmt.Errorf("profile update: delete work history: %w", err)
	}

	if len(history) == 0 {
		return nil
	}

	const query = `
INSERT INTO profile_work_history (
    profile_id,
    organization_ja,
    organization_en,
    role_ja,
    role_en,
    summary_ja,
    summary_en,
    started_at,
    ended_at,
    external_url,
    sort_order
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, item := range history {
		var endedAt interface{}
		if item.EndedAt != nil && !item.EndedAt.IsZero() {
			endedAt = item.EndedAt.UTC()
		} else {
			endedAt = nil
		}

		if _, err := tx.ExecContext(
			ctx,
			query,
			profileID,
			strings.TrimSpace(item.Organization.Ja),
			strings.TrimSpace(item.Organization.En),
			strings.TrimSpace(item.Role.Ja),
			strings.TrimSpace(item.Role.En),
			strings.TrimSpace(item.Summary.Ja),
			strings.TrimSpace(item.Summary.En),
			item.StartedAt.UTC(),
			endedAt,
			nullString(item.ExternalURL),
			item.SortOrder,
		); err != nil {
			return fmt.Errorf("profile update: insert work history: %w", err)
		}
	}
	return nil
}

func (r *profileRepository) replaceSocialLinks(ctx context.Context, tx *sqlx.Tx, profileID uint64, links []model.ProfileSocialLink) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM profile_social_links WHERE profile_id = ?`, profileID); err != nil {
		return fmt.Errorf("profile update: delete social links: %w", err)
	}

	if len(links) == 0 {
		return nil
	}

	const query = `
INSERT INTO profile_social_links (
    profile_id,
    provider,
    label_ja,
    label_en,
    url,
    is_footer,
    sort_order
) VALUES (?, ?, ?, ?, ?, ?, ?)`

	for _, link := range links {
		isFooter := 0
		if link.IsFooter {
			isFooter = 1
		}
		if _, err := tx.ExecContext(
			ctx,
			query,
			profileID,
			string(link.Provider),
			strings.TrimSpace(link.Label.Ja),
			strings.TrimSpace(link.Label.En),
			strings.TrimSpace(link.URL),
			isFooter,
			link.SortOrder,
		); err != nil {
			return fmt.Errorf("profile update: insert social link: %w", err)
		}
	}
	return nil
}

var _ repository.ProfileRepository = (*profileRepository)(nil)
var _ repository.AdminProfileRepository = (*profileRepository)(nil)
