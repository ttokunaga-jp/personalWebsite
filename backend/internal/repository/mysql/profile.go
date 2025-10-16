package mysql

import (
	"context"
	"database/sql"
	"fmt"

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
	summary_en
FROM profile
ORDER BY id
LIMIT 1`

	profileSkillsQuery = `
SELECT
	skill_ja,
	skill_en
FROM profile_skills
WHERE profile_id = ?
ORDER BY sort_order, id`
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
}

type profileSkillRow struct {
	SkillJA sql.NullString `db:"skill_ja"`
	SkillEN sql.NullString `db:"skill_en"`
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	var row profileRow
	if err := r.db.GetContext(ctx, &row, profileQuery); err != nil {
		return nil, fmt.Errorf("select profile: %w", err)
	}

	profile := &model.Profile{
		Name:        toLocalizedText(row.NameJA, row.NameEN),
		Title:       toLocalizedText(row.TitleJA, row.TitleEN),
		Affiliation: toLocalizedText(row.AffiliationJA, row.AffiliationEN),
		Lab:         toLocalizedText(row.LabJA, row.LabEN),
		Summary:     toLocalizedText(row.SummaryJA, row.SummaryEN),
	}

	var skills []profileSkillRow
	if err := r.db.SelectContext(ctx, &skills, profileSkillsQuery, row.ID); err != nil {
		return nil, fmt.Errorf("select profile skills: %w", err)
	}

	profile.Skills = make([]model.LocalizedText, 0, len(skills))
	for _, s := range skills {
		profile.Skills = append(profile.Skills, toLocalizedText(s.SkillJA, s.SkillEN))
	}

	return profile, nil
}
