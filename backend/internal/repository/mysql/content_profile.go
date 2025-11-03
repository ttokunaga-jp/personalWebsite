package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// contentProfileRepository implements repository.ContentProfileRepository.
type contentProfileRepository struct {
	db *sqlx.DB
}

// NewContentProfileRepository returns a MySQL-backed implementation for ProfileDocument retrieval.
func NewContentProfileRepository(db *sqlx.DB) repository.ContentProfileRepository {
	return &contentProfileRepository{db: db}
}

const (
	profileDocumentQuery = `
SELECT
    id,
    display_name,
    headline_ja,
    headline_en,
    summary_ja,
    summary_en,
    avatar_url,
    location_ja,
    location_en,
    theme_mode,
    theme_accent_color,
    lab_name_ja,
    lab_name_en,
    lab_advisor_ja,
    lab_advisor_en,
    lab_room_ja,
    lab_room_en,
    lab_url,
    updated_at
FROM profiles
ORDER BY id
LIMIT 1`

	profileAffiliationsQuery = `
SELECT
    id,
    profile_id,
    kind,
    name,
    url,
    description_ja,
    description_en,
    started_at,
    sort_order
FROM profile_affiliations
WHERE profile_id = ?
ORDER BY sort_order, id`

	profileWorkHistoryQuery = `
SELECT
    id,
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
FROM profile_work_history
WHERE profile_id = ?
ORDER BY sort_order, id`

	profileSocialLinksQuery = `
SELECT
    id,
    profile_id,
    provider,
    label_ja,
    label_en,
    url,
    is_footer,
    sort_order
FROM profile_social_links
WHERE profile_id = ?
ORDER BY sort_order, id`

	profileTechSectionsQuery = `
SELECT
    id,
    profile_id,
    title_ja,
    title_en,
    layout,
    breakpoint,
    sort_order
FROM profile_tech_sections
WHERE profile_id = ?
ORDER BY sort_order, id`

	profileTechMembershipsQuery = `
SELECT
    tr.id            AS membership_id,
    tr.entity_id     AS section_id,
    tr.context       AS context,
    tr.note          AS note,
    tr.sort_order    AS membership_sort_order,
    tc.id            AS tech_id,
    tc.slug          AS tech_slug,
    tc.display_name  AS tech_display_name,
    tc.category      AS tech_category,
    tc.level         AS tech_level,
    tc.icon          AS tech_icon,
    tc.sort_order    AS tech_sort_order,
    tc.is_active     AS tech_is_active,
    tc.created_at    AS tech_created_at,
    tc.updated_at    AS tech_updated_at
FROM tech_relationships tr
JOIN tech_catalog tc ON tc.id = tr.tech_id
WHERE tr.entity_type = 'profile_section' AND tr.entity_id IN (?)
ORDER BY tr.sort_order, tr.id`
)

type contentProfileRow struct {
	ID           uint64         `db:"id"`
	DisplayName  sql.NullString `db:"display_name"`
	HeadlineJA   sql.NullString `db:"headline_ja"`
	HeadlineEN   sql.NullString `db:"headline_en"`
	SummaryJA    sql.NullString `db:"summary_ja"`
	SummaryEN    sql.NullString `db:"summary_en"`
	AvatarURL    sql.NullString `db:"avatar_url"`
	LocationJA   sql.NullString `db:"location_ja"`
	LocationEN   sql.NullString `db:"location_en"`
	ThemeMode    sql.NullString `db:"theme_mode"`
	ThemeAccent  sql.NullString `db:"theme_accent_color"`
	LabNameJA    sql.NullString `db:"lab_name_ja"`
	LabNameEN    sql.NullString `db:"lab_name_en"`
	LabAdvisorJA sql.NullString `db:"lab_advisor_ja"`
	LabAdvisorEN sql.NullString `db:"lab_advisor_en"`
	LabRoomJA    sql.NullString `db:"lab_room_ja"`
	LabRoomEN    sql.NullString `db:"lab_room_en"`
	LabURL       sql.NullString `db:"lab_url"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
}

type contentProfileAffiliationRow struct {
	ID            uint64         `db:"id"`
	ProfileID     uint64         `db:"profile_id"`
	Kind          string         `db:"kind"`
	Name          sql.NullString `db:"name"`
	URL           sql.NullString `db:"url"`
	DescriptionJA sql.NullString `db:"description_ja"`
	DescriptionEN sql.NullString `db:"description_en"`
	StartedAt     time.Time      `db:"started_at"`
	SortOrder     int            `db:"sort_order"`
}

type contentProfileWorkHistoryRow struct {
	ID             uint64         `db:"id"`
	ProfileID      uint64         `db:"profile_id"`
	OrganizationJA sql.NullString `db:"organization_ja"`
	OrganizationEN sql.NullString `db:"organization_en"`
	RoleJA         sql.NullString `db:"role_ja"`
	RoleEN         sql.NullString `db:"role_en"`
	SummaryJA      sql.NullString `db:"summary_ja"`
	SummaryEN      sql.NullString `db:"summary_en"`
	StartedAt      time.Time      `db:"started_at"`
	EndedAt        sql.NullTime   `db:"ended_at"`
	ExternalURL    sql.NullString `db:"external_url"`
	SortOrder      int            `db:"sort_order"`
}

type contentProfileSocialLinkRow struct {
	ID        uint64         `db:"id"`
	ProfileID uint64         `db:"profile_id"`
	Provider  string         `db:"provider"`
	LabelJA   sql.NullString `db:"label_ja"`
	LabelEN   sql.NullString `db:"label_en"`
	URL       sql.NullString `db:"url"`
	IsFooter  bool           `db:"is_footer"`
	SortOrder int            `db:"sort_order"`
}

type contentProfileTechSectionRow struct {
	ID         uint64         `db:"id"`
	ProfileID  uint64         `db:"profile_id"`
	TitleJA    sql.NullString `db:"title_ja"`
	TitleEN    sql.NullString `db:"title_en"`
	Layout     sql.NullString `db:"layout"`
	Breakpoint sql.NullString `db:"breakpoint"`
	SortOrder  int            `db:"sort_order"`
}

type contentProfileTechMembershipRow struct {
	MembershipID        uint64         `db:"membership_id"`
	SectionID           uint64         `db:"section_id"`
	Context             sql.NullString `db:"context"`
	Note                sql.NullString `db:"note"`
	MembershipSortOrder int            `db:"membership_sort_order"`
	TechID              uint64         `db:"tech_id"`
	TechSlug            sql.NullString `db:"tech_slug"`
	TechDisplayName     sql.NullString `db:"tech_display_name"`
	TechCategory        sql.NullString `db:"tech_category"`
	TechLevel           sql.NullString `db:"tech_level"`
	TechIcon            sql.NullString `db:"tech_icon"`
	TechSortOrder       int            `db:"tech_sort_order"`
	TechIsActive        bool           `db:"tech_is_active"`
	TechCreatedAt       time.Time      `db:"tech_created_at"`
	TechUpdatedAt       time.Time      `db:"tech_updated_at"`
}

func (r *contentProfileRepository) GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error) {
	var row contentProfileRow
	if err := r.db.GetContext(ctx, &row, profileDocumentQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("select profile document: %w", err)
	}

	profileID := row.ID

	affiliations, err := r.loadAffiliations(ctx, profileID)
	if err != nil {
		return nil, err
	}

	workHistory, err := r.loadWorkHistory(ctx, profileID)
	if err != nil {
		return nil, err
	}

	socialLinks, err := r.loadSocialLinks(ctx, profileID)
	if err != nil {
		return nil, err
	}

	sections, err := r.loadTechSections(ctx, profileID)
	if err != nil {
		return nil, err
	}

	document := &model.ProfileDocument{
		ID:          profileID,
		DisplayName: strings.TrimSpace(row.DisplayName.String),
		Headline:    toLocalizedText(row.HeadlineJA, row.HeadlineEN),
		Summary:     toLocalizedText(row.SummaryJA, row.SummaryEN),
		AvatarURL:   nullableString(row.AvatarURL),
		Location:    toLocalizedText(row.LocationJA, row.LocationEN),
		Theme: model.ProfileTheme{
			Mode:        model.ProfileThemeMode(strings.ToLower(strings.TrimSpace(row.ThemeMode.String))),
			AccentColor: nullableString(row.ThemeAccent),
		},
		Lab: model.ProfileLab{
			Name:    toLocalizedText(row.LabNameJA, row.LabNameEN),
			Advisor: toLocalizedText(row.LabAdvisorJA, row.LabAdvisorEN),
			Room:    toLocalizedText(row.LabRoomJA, row.LabRoomEN),
			URL:     nullableString(row.LabURL),
		},
		Affiliations: affiliations.affiliations,
		Communities:  affiliations.communities,
		WorkHistory:  workHistory,
		SocialLinks:  socialLinks,
		TechSections: sections,
	}

	if row.UpdatedAt.Valid {
		document.UpdatedAt = row.UpdatedAt.Time
	} else {
		document.UpdatedAt = time.Now().UTC()
	}

	return document, nil
}

type affiliationResult struct {
	affiliations []model.ProfileAffiliation
	communities  []model.ProfileAffiliation
}

func (r *contentProfileRepository) loadAffiliations(ctx context.Context, profileID uint64) (affiliationResult, error) {
	var rows []contentProfileAffiliationRow
	if err := r.db.SelectContext(ctx, &rows, profileAffiliationsQuery, profileID); err != nil {
		return affiliationResult{}, fmt.Errorf("profile affiliations: %w", err)
	}

	result := affiliationResult{
		affiliations: make([]model.ProfileAffiliation, 0, len(rows)),
		communities:  make([]model.ProfileAffiliation, 0, len(rows)),
	}

	for _, row := range rows {
		affiliation := model.ProfileAffiliation{
			ID:          row.ID,
			ProfileID:   row.ProfileID,
			Kind:        model.ProfileAffiliationKind(row.Kind),
			Name:        strings.TrimSpace(row.Name.String),
			URL:         nullableString(row.URL),
			Description: toLocalizedText(row.DescriptionJA, row.DescriptionEN),
			StartedAt:   row.StartedAt,
			SortOrder:   row.SortOrder,
		}

		if affiliation.Kind == model.ProfileAffiliationKindCommunity {
			result.communities = append(result.communities, affiliation)
		} else {
			result.affiliations = append(result.affiliations, affiliation)
		}
	}

	return result, nil
}

func (r *contentProfileRepository) loadWorkHistory(ctx context.Context, profileID uint64) ([]model.ProfileWorkExperience, error) {
	var rows []contentProfileWorkHistoryRow
	if err := r.db.SelectContext(ctx, &rows, profileWorkHistoryQuery, profileID); err != nil {
		return nil, fmt.Errorf("profile work history: %w", err)
	}

	items := make([]model.ProfileWorkExperience, 0, len(rows))
	for _, row := range rows {
		item := model.ProfileWorkExperience{
			ID:           row.ID,
			ProfileID:    row.ProfileID,
			Organization: toLocalizedText(row.OrganizationJA, row.OrganizationEN),
			Role:         toLocalizedText(row.RoleJA, row.RoleEN),
			Summary:      toLocalizedText(row.SummaryJA, row.SummaryEN),
			StartedAt:    row.StartedAt,
			SortOrder:    row.SortOrder,
			ExternalURL:  nullableString(row.ExternalURL),
		}
		if row.EndedAt.Valid {
			ended := row.EndedAt.Time
			item.EndedAt = &ended
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *contentProfileRepository) loadSocialLinks(ctx context.Context, profileID uint64) ([]model.ProfileSocialLink, error) {
	var rows []contentProfileSocialLinkRow
	if err := r.db.SelectContext(ctx, &rows, profileSocialLinksQuery, profileID); err != nil {
		return nil, fmt.Errorf("profile social links: %w", err)
	}

	links := make([]model.ProfileSocialLink, 0, len(rows))
	for _, row := range rows {
		link := model.ProfileSocialLink{
			ID:        row.ID,
			ProfileID: row.ProfileID,
			Provider:  model.ProfileSocialProvider(strings.TrimSpace(row.Provider)),
			Label:     toLocalizedText(row.LabelJA, row.LabelEN),
			URL:       nullableString(row.URL),
			IsFooter:  row.IsFooter,
			SortOrder: row.SortOrder,
		}
		links = append(links, link)
	}
	return links, nil
}

func (r *contentProfileRepository) loadTechSections(ctx context.Context, profileID uint64) ([]model.ProfileTechSection, error) {
	var rows []contentProfileTechSectionRow
	if err := r.db.SelectContext(ctx, &rows, profileTechSectionsQuery, profileID); err != nil {
		return nil, fmt.Errorf("profile tech sections: %w", err)
	}

	if len(rows) == 0 {
		return []model.ProfileTechSection{}, nil
	}

	sectionIDs := make([]interface{}, 0, len(rows))
	for _, row := range rows {
		sectionIDs = append(sectionIDs, row.ID)
	}

	query, args, err := sqlx.In(profileTechMembershipsQuery, sectionIDs)
	if err != nil {
		return nil, fmt.Errorf("profile tech membership query compose: %w", err)
	}
	query = r.db.Rebind(query)

	var membershipRows []contentProfileTechMembershipRow
	if err := r.db.SelectContext(ctx, &membershipRows, query, args...); err != nil {
		return nil, fmt.Errorf("profile tech memberships: %w", err)
	}

	membershipMap := make(map[uint64][]model.TechMembership)
	for _, mem := range membershipRows {
		membershipMap[mem.SectionID] = append(membershipMap[mem.SectionID], model.TechMembership{
			MembershipID: mem.MembershipID,
			EntityType:   "profile_section",
			EntityID:     mem.SectionID,
			Tech: model.TechCatalogEntry{
				ID:          mem.TechID,
				Slug:        strings.TrimSpace(mem.TechSlug.String),
				DisplayName: strings.TrimSpace(mem.TechDisplayName.String),
				Category:    strings.TrimSpace(mem.TechCategory.String),
				Level:       model.TechLevel(strings.TrimSpace(mem.TechLevel.String)),
				Icon:        nullableString(mem.TechIcon),
				SortOrder:   mem.TechSortOrder,
				Active:      mem.TechIsActive,
				CreatedAt:   mem.TechCreatedAt,
				UpdatedAt:   mem.TechUpdatedAt,
			},
			Context:   model.TechContext(strings.TrimSpace(mem.Context.String)),
			Note:      nullableString(mem.Note),
			SortOrder: mem.MembershipSortOrder,
		})
	}

	sections := make([]model.ProfileTechSection, 0, len(rows))
	for _, row := range rows {
		section := model.ProfileTechSection{
			ID:         row.ID,
			ProfileID:  row.ProfileID,
			Title:      toLocalizedText(row.TitleJA, row.TitleEN),
			Layout:     strings.TrimSpace(row.Layout.String),
			Breakpoint: strings.TrimSpace(row.Breakpoint.String),
			SortOrder:  row.SortOrder,
			Members:    membershipMap[row.ID],
		}
		sections = append(sections, section)
	}

	return sections, nil
}

var _ repository.ContentProfileRepository = (*contentProfileRepository)(nil)
