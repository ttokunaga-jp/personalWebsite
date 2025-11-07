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

type projectRepository struct {
	db *sqlx.DB
}

// NewProjectRepository returns a MySQL-backed project repository.
func NewProjectRepository(db *sqlx.DB) repository.ProjectRepository {
	return &projectRepository{db: db}
}

const listProjectsQuery = `
SELECT
	p.id,
	p.year,
	p.link_url,
	p.title_ja,
	p.title_en,
	p.description_ja,
	p.description_en
FROM projects p
WHERE p.published = TRUE
ORDER BY COALESCE(p.sort_order, p.year * 1000), p.year DESC, p.id`

const listAdminProjectsQuery = `
SELECT
	p.id,
	p.year,
	p.link_url,
	p.title_ja,
	p.title_en,
	p.description_ja,
	p.description_en,
	p.published,
	p.sort_order,
	p.created_at,
	p.updated_at
FROM projects p
ORDER BY COALESCE(p.sort_order, p.year * 1000), p.year DESC, p.id`

const getAdminProjectQuery = `
SELECT
	p.id,
	p.year,
	p.link_url,
	p.title_ja,
	p.title_en,
	p.description_ja,
	p.description_en,
	p.published,
	p.sort_order,
	p.created_at,
	p.updated_at
FROM projects p
WHERE p.id = ?`

const insertProjectQuery = `
INSERT INTO projects (
	title_ja,
	title_en,
	description_ja,
	description_en,
	link_url,
	year,
	published,
	sort_order,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

const updateProjectQuery = `
UPDATE projects
SET
	title_ja = ?,
	title_en = ?,
	description_ja = ?,
	description_en = ?,
	link_url = ?,
	year = ?,
	published = ?,
	sort_order = ?,
	updated_at = NOW()
WHERE id = ?`

const (
	deleteProjectQuery = `DELETE FROM projects WHERE id = ?`

	projectEntityType      = "project"
	deleteProjectTechQuery = `DELETE FROM tech_relationships WHERE entity_type = ? AND entity_id = ?`
	insertProjectTechQuery = `INSERT INTO tech_relationships (entity_type, entity_id, tech_id, context, note, sort_order) VALUES (?, ?, ?, ?, ?, ?)`
	selectProjectTechQuery = `
SELECT
	tr.id              AS membership_id,
	tr.entity_id       AS entity_id,
	tr.context         AS context,
	tr.note            AS note,
	tr.sort_order      AS membership_sort_order,
	tc.id              AS tech_id,
	tc.slug            AS tech_slug,
	tc.display_name    AS tech_display_name,
	tc.category        AS tech_category,
	tc.level           AS tech_level,
	tc.icon            AS tech_icon,
	tc.sort_order      AS tech_sort_order,
	tc.is_active       AS tech_is_active,
	tc.created_at      AS tech_created_at,
	tc.updated_at      AS tech_updated_at
FROM tech_relationships tr
JOIN tech_catalog tc ON tc.id = tr.tech_id
WHERE tr.entity_type = ? AND tr.entity_id = ?
ORDER BY tr.sort_order, tr.id`
)

type projectRow struct {
	ID            int64          `db:"id"`
	Year          int            `db:"year"`
	LinkURL       sql.NullString `db:"link_url"`
	TitleJA       sql.NullString `db:"title_ja"`
	TitleEN       sql.NullString `db:"title_en"`
	DescriptionJA sql.NullString `db:"description_ja"`
	DescriptionEN sql.NullString `db:"description_en"`
	Published     sql.NullBool   `db:"published"`
	SortOrder     sql.NullInt64  `db:"sort_order"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
}

type projectTechRow struct {
	MembershipID    uint64         `db:"membership_id"`
	EntityID        int64          `db:"entity_id"`
	Context         string         `db:"context"`
	Note            sql.NullString `db:"note"`
	SortOrder       int            `db:"membership_sort_order"`
	TechID          uint64         `db:"tech_id"`
	TechSlug        string         `db:"tech_slug"`
	TechDisplayName string         `db:"tech_display_name"`
	TechCategory    sql.NullString `db:"tech_category"`
	TechLevel       string         `db:"tech_level"`
	TechIcon        sql.NullString `db:"tech_icon"`
	TechSortOrder   int            `db:"tech_sort_order"`
	TechActive      bool           `db:"tech_is_active"`
	TechCreatedAt   sql.NullTime   `db:"tech_created_at"`
	TechUpdatedAt   sql.NullTime   `db:"tech_updated_at"`
}

func (r *projectRepository) ListProjects(ctx context.Context) ([]model.Project, error) {
	var rows []projectRow
	if err := r.db.SelectContext(ctx, &rows, listProjectsQuery); err != nil {
		return nil, fmt.Errorf("select projects: %w", err)
	}

	projects := make([]model.Project, 0, len(rows))
	for _, row := range rows {
		memberships, err := r.fetchProjectTech(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		projects = append(projects, model.Project{
			ID:          row.ID,
			Title:       toLocalizedText(row.TitleJA, row.TitleEN),
			Description: toLocalizedText(row.DescriptionJA, row.DescriptionEN),
			Tech:        append([]model.TechMembership(nil), memberships...),
			TechStack:   extractTechDisplayNames(memberships),
			LinkURL:     nullableString(row.LinkURL),
			Year:        row.Year,
		})
	}

	return projects, nil
}

func (r *projectRepository) ListAdminProjects(ctx context.Context) ([]model.AdminProject, error) {
	var rows []projectRow
	if err := r.db.SelectContext(ctx, &rows, listAdminProjectsQuery); err != nil {
		return nil, fmt.Errorf("select admin projects: %w", err)
	}

	projects := make([]model.AdminProject, 0, len(rows))
	for _, row := range rows {
		memberships, err := r.fetchProjectTech(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, mapProjectRow(row, memberships))
	}

	return projects, nil
}

func (r *projectRepository) GetAdminProject(ctx context.Context, id int64) (*model.AdminProject, error) {
	var row projectRow
	if err := r.db.GetContext(ctx, &row, getAdminProjectQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get project %d: %w", id, err)
	}

	memberships, err := r.fetchProjectTech(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	project := mapProjectRow(row, memberships)
	return &project, nil
}

func (r *projectRepository) CreateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error) {
	if project == nil {
		return nil, repository.ErrInvalidInput
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer rollbackOnError(tx, &err)

	res, execErr := tx.ExecContext(ctx, insertProjectQuery,
		project.Title.Ja,
		project.Title.En,
		project.Description.Ja,
		project.Description.En,
		nullString(project.LinkURL),
		project.Year,
		project.Published,
		nullInt(project.SortOrder),
	)
	if execErr != nil {
		err = fmt.Errorf("insert project: %w", execErr)
		return nil, err
	}

	projectID, execErr := res.LastInsertId()
	if execErr != nil {
		err = fmt.Errorf("project last insert id: %w", execErr)
		return nil, err
	}

	if err = r.replaceProjectTechTx(ctx, tx, projectID, project.Tech); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit project insert: %w", err)
	}

	created, err := r.GetAdminProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *projectRepository) UpdateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error) {
	if project == nil {
		return nil, repository.ErrInvalidInput
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer rollbackOnError(tx, &err)

	res, execErr := tx.ExecContext(ctx, updateProjectQuery,
		project.Title.Ja,
		project.Title.En,
		project.Description.Ja,
		project.Description.En,
		nullString(project.LinkURL),
		project.Year,
		project.Published,
		nullInt(project.SortOrder),
		project.ID,
	)
	if execErr != nil {
		err = fmt.Errorf("update project %d: %w", project.ID, execErr)
		return nil, err
	}

	affected, execErr := res.RowsAffected()
	if execErr != nil {
		err = fmt.Errorf("rows affected project %d: %w", project.ID, execErr)
		return nil, err
	}
	if affected == 0 {
		err = repository.ErrNotFound
		return nil, err
	}

	if err = r.replaceProjectTechTx(ctx, tx, project.ID, project.Tech); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit project update: %w", err)
	}

	updated, err := r.GetAdminProject(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *projectRepository) DeleteAdminProject(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer rollbackOnError(tx, &err)

	if _, execErr := tx.ExecContext(ctx, deleteProjectTechQuery, projectEntityType, id); execErr != nil {
		err = fmt.Errorf("delete project tech %d: %w", id, execErr)
		return err
	}

	res, execErr := tx.ExecContext(ctx, deleteProjectQuery, id)
	if execErr != nil {
		err = fmt.Errorf("delete project %d: %w", id, execErr)
		return err
	}

	affected, execErr := res.RowsAffected()
	if execErr != nil {
		err = fmt.Errorf("rows affected delete project %d: %w", id, execErr)
		return err
	}
	if affected == 0 {
		err = repository.ErrNotFound
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit project delete: %w", err)
	}

	return nil
}

func (r *projectRepository) fetchProjectTech(ctx context.Context, projectID int64) ([]model.TechMembership, error) {
	var rows []projectTechRow
	if err := r.db.SelectContext(ctx, &rows, selectProjectTechQuery, projectEntityType, projectID); err != nil {
		return nil, fmt.Errorf("select project tech (id=%d): %w", projectID, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	memberships := make([]model.TechMembership, 0, len(rows))
	for _, row := range rows {
		entry := model.TechCatalogEntry{
			ID:          row.TechID,
			Slug:        strings.TrimSpace(row.TechSlug),
			DisplayName: strings.TrimSpace(row.TechDisplayName),
			Category:    strings.TrimSpace(row.TechCategory.String),
			Level:       model.TechLevel(strings.TrimSpace(row.TechLevel)),
			Icon:        strings.TrimSpace(row.TechIcon.String),
			SortOrder:   row.TechSortOrder,
			Active:      row.TechActive,
		}
		if row.TechCreatedAt.Valid {
			entry.CreatedAt = row.TechCreatedAt.Time.UTC()
		}
		if row.TechUpdatedAt.Valid {
			entry.UpdatedAt = row.TechUpdatedAt.Time.UTC()
		}

		memberships = append(memberships, model.TechMembership{
			MembershipID: row.MembershipID,
			EntityType:   projectEntityType,
			EntityID:     uint64(row.EntityID),
			Tech:         entry,
			Context:      model.TechContext(strings.TrimSpace(row.Context)),
			Note:         strings.TrimSpace(row.Note.String),
			SortOrder:    row.SortOrder,
		})
	}

	return memberships, nil
}

func (r *projectRepository) replaceProjectTechTx(ctx context.Context, tx *sqlx.Tx, projectID int64, tech []model.TechMembership) error {
	if _, err := tx.ExecContext(ctx, deleteProjectTechQuery, projectEntityType, projectID); err != nil {
		return fmt.Errorf("clear project tech %d: %w", projectID, err)
	}

	for _, membership := range tech {
		if membership.Tech.ID == 0 {
			continue
		}
		if _, err := tx.ExecContext(
			ctx,
			insertProjectTechQuery,
			projectEntityType,
			projectID,
			membership.Tech.ID,
			membership.Context,
			nullString(membership.Note),
			membership.SortOrder,
		); err != nil {
			return fmt.Errorf("insert project tech %d (tech=%d): %w", projectID, membership.Tech.ID, err)
		}
	}
	return nil
}

func extractTechDisplayNames(memberships []model.TechMembership) []string {
	if len(memberships) == 0 {
		return nil
	}
	names := make([]string, 0, len(memberships))
	for _, membership := range memberships {
		name := strings.TrimSpace(membership.Tech.DisplayName)
		if name != "" {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return nil
	}
	return names
}

func mapProjectRow(row projectRow, tech []model.TechMembership) model.AdminProject {
	var sortOrder *int
	if row.SortOrder.Valid {
		v := int(row.SortOrder.Int64)
		sortOrder = &v
	}

	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = time.Time{}
	}
	updatedAt := row.UpdatedAt.Time
	if !row.UpdatedAt.Valid {
		updatedAt = time.Time{}
	}

	return model.AdminProject{
		ID:          row.ID,
		Title:       toLocalizedText(row.TitleJA, row.TitleEN),
		Description: toLocalizedText(row.DescriptionJA, row.DescriptionEN),
		Tech:        append([]model.TechMembership(nil), tech...),
		LinkURL:     nullableString(row.LinkURL),
		Year:        row.Year,
		Published:   row.Published.Bool,
		SortOrder:   sortOrder,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

var _ repository.ProjectRepository = (*projectRepository)(nil)
var _ repository.AdminProjectRepository = (*projectRepository)(nil)
