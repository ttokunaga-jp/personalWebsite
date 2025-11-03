package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type profileRepository struct {
	db *sqlx.DB
}

// NewProfileRepository returns a repository backed by MySQL. Caller must ensure db is non-nil.
func NewProfileRepository(db *sqlx.DB) repository.ProfileRepository {
	return &profileRepository{db: db}
}

const (
	profileQuery = `
SELECT
	id,
	name_ja,
	name_en,
	title_ja,
	title_en,
	affiliation_ja,
	affiliation_en,
	lab_ja,
	lab_en,
	summary_ja,
	summary_en,
	updated_at
FROM profile
ORDER BY id
LIMIT 1`

	profileSkillsQuery = `
SELECT
	id,
	skill_ja,
	skill_en,
	sort_order
FROM profile_skills
WHERE profile_id = ?
ORDER BY sort_order, id`

	profileFocusAreasQuery = `
SELECT
	id,
	area_ja,
	area_en,
	sort_order
FROM profile_focus_areas
WHERE profile_id = ?
ORDER BY sort_order, id`

	updateProfileQuery = `
UPDATE profile SET
	name_ja = ?,
	name_en = ?,
	title_ja = ?,
	title_en = ?,
	affiliation_ja = ?,
	affiliation_en = ?,
	lab_ja = ?,
	lab_en = ?,
	summary_ja = ?,
	summary_en = ?,
	updated_at = NOW()
WHERE id = ?`

	deleteProfileSkillsQuery     = `DELETE FROM profile_skills WHERE profile_id = ?`
	insertProfileSkillQuery      = `INSERT INTO profile_skills (profile_id, skill_ja, skill_en, sort_order) VALUES (?, ?, ?, ?)`
	deleteProfileFocusAreasQuery = `DELETE FROM profile_focus_areas WHERE profile_id = ?`
	insertProfileFocusAreaQuery  = `INSERT INTO profile_focus_areas (profile_id, area_ja, area_en, sort_order) VALUES (?, ?, ?, ?)`
)

type profileRow struct {
	ID            int64          `db:"id"`
	NameJA        sql.NullString `db:"name_ja"`
	NameEN        sql.NullString `db:"name_en"`
	TitleJA       sql.NullString `db:"title_ja"`
	TitleEN       sql.NullString `db:"title_en"`
	AffiliationJA sql.NullString `db:"affiliation_ja"`
	AffiliationEN sql.NullString `db:"affiliation_en"`
	LabJA         sql.NullString `db:"lab_ja"`
	LabEN         sql.NullString `db:"lab_en"`
	SummaryJA     sql.NullString `db:"summary_ja"`
	SummaryEN     sql.NullString `db:"summary_en"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
}

type profileSkillRow struct {
	ID        int64          `db:"id"`
	SkillJA   sql.NullString `db:"skill_ja"`
	SkillEN   sql.NullString `db:"skill_en"`
	SortOrder sql.NullInt64  `db:"sort_order"`
}

type profileFocusAreaRow struct {
	ID        int64          `db:"id"`
	AreaJA    sql.NullString `db:"area_ja"`
	AreaEN    sql.NullString `db:"area_en"`
	SortOrder sql.NullInt64  `db:"sort_order"`
}

func (r *profileRepository) loadProfile(ctx context.Context) (*profileRow, []profileSkillRow, []profileFocusAreaRow, error) {
	var row profileRow
	if err := r.db.GetContext(ctx, &row, profileQuery); err != nil {
		return nil, nil, nil, fmt.Errorf("select profile: %w", err)
	}

	var skills []profileSkillRow
	if err := r.db.SelectContext(ctx, &skills, profileSkillsQuery, row.ID); err != nil {
		return nil, nil, nil, fmt.Errorf("select profile skills: %w", err)
	}

	var focusAreas []profileFocusAreaRow
	if err := r.db.SelectContext(ctx, &focusAreas, profileFocusAreasQuery, row.ID); err != nil {
		return nil, nil, nil, fmt.Errorf("select profile focus areas: %w", err)
	}

	return &row, skills, focusAreas, nil
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	row, skills, focusAreas, err := r.loadProfile(ctx)
	if err != nil {
		return nil, err
	}

	profile := &model.Profile{
		Name:        toLocalizedText(row.NameJA, row.NameEN),
		Title:       toLocalizedText(row.TitleJA, row.TitleEN),
		Affiliation: toLocalizedText(row.AffiliationJA, row.AffiliationEN),
		Lab:         toLocalizedText(row.LabJA, row.LabEN),
		Summary:     toLocalizedText(row.SummaryJA, row.SummaryEN),
		Skills:      make([]model.LocalizedText, 0, len(skills)),
		FocusAreas:  make([]model.LocalizedText, 0, len(focusAreas)),
	}

	for _, s := range skills {
		profile.Skills = append(profile.Skills, toLocalizedText(s.SkillJA, s.SkillEN))
	}

	for _, area := range focusAreas {
		profile.FocusAreas = append(profile.FocusAreas, toLocalizedText(area.AreaJA, area.AreaEN))
	}

	return profile, nil
}

func (r *profileRepository) GetAdminProfile(ctx context.Context) (*model.AdminProfile, error) {
	row, skills, focusAreas, err := r.loadProfile(ctx)
	if err != nil {
		return nil, err
	}

	admin := &model.AdminProfile{
		Name:        toLocalizedText(row.NameJA, row.NameEN),
		Title:       toLocalizedText(row.TitleJA, row.TitleEN),
		Affiliation: toLocalizedText(row.AffiliationJA, row.AffiliationEN),
		Lab:         toLocalizedText(row.LabJA, row.LabEN),
		Summary:     toLocalizedText(row.SummaryJA, row.SummaryEN),
		Skills:      make([]model.LocalizedText, 0, len(skills)),
		FocusAreas:  make([]model.LocalizedText, 0, len(focusAreas)),
		UpdatedAt:   nullableTime(row.UpdatedAt),
	}

	for _, s := range skills {
		admin.Skills = append(admin.Skills, toLocalizedText(s.SkillJA, s.SkillEN))
	}

	for _, area := range focusAreas {
		admin.FocusAreas = append(admin.FocusAreas, toLocalizedText(area.AreaJA, area.AreaEN))
	}

	return admin, nil
}

func (r *profileRepository) UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error) {
	if profile == nil {
		return nil, repository.ErrInvalidInput
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin profile tx: %w", err)
	}

	var opErr error
	defer rollbackOnError(tx, &opErr)

	var row profileRow
	if err = tx.GetContext(ctx, &row, profileQuery); err != nil {
		opErr = fmt.Errorf("select profile for update: %w", err)
		return nil, opErr
	}

	if _, err = tx.ExecContext(
		ctx,
		updateProfileQuery,
		strings.TrimSpace(profile.Name.Ja),
		strings.TrimSpace(profile.Name.En),
		strings.TrimSpace(profile.Title.Ja),
		strings.TrimSpace(profile.Title.En),
		strings.TrimSpace(profile.Affiliation.Ja),
		strings.TrimSpace(profile.Affiliation.En),
		strings.TrimSpace(profile.Lab.Ja),
		strings.TrimSpace(profile.Lab.En),
		strings.TrimSpace(profile.Summary.Ja),
		strings.TrimSpace(profile.Summary.En),
		row.ID,
	); err != nil {
		opErr = fmt.Errorf("update profile: %w", err)
		return nil, opErr
	}

	if _, err = tx.ExecContext(ctx, deleteProfileSkillsQuery, row.ID); err != nil {
		opErr = fmt.Errorf("truncate profile skills: %w", err)
		return nil, opErr
	}
	for idx, skill := range profile.Skills {
		if _, err = tx.ExecContext(
			ctx,
			insertProfileSkillQuery,
			row.ID,
			strings.TrimSpace(skill.Ja),
			strings.TrimSpace(skill.En),
			idx,
		); err != nil {
			opErr = fmt.Errorf("insert profile skill %d: %w", idx, err)
			return nil, opErr
		}
	}

	if _, err = tx.ExecContext(ctx, deleteProfileFocusAreasQuery, row.ID); err != nil {
		opErr = fmt.Errorf("truncate profile focus areas: %w", err)
		return nil, opErr
	}
	for idx, area := range profile.FocusAreas {
		if _, err = tx.ExecContext(
			ctx,
			insertProfileFocusAreaQuery,
			row.ID,
			strings.TrimSpace(area.Ja),
			strings.TrimSpace(area.En),
			idx,
		); err != nil {
			opErr = fmt.Errorf("insert profile focus area %d: %w", idx, err)
			return nil, opErr
		}
	}

	if err = tx.Commit(); err != nil {
		opErr = fmt.Errorf("commit profile update: %w", err)
		return nil, opErr
	}

	updated, err := r.GetAdminProfile(ctx)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

var _ repository.ProfileRepository = (*profileRepository)(nil)
var _ repository.AdminProfileRepository = (*profileRepository)(nil)
