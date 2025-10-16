package mysql

import (
	"context"
	"database/sql"
	"fmt"
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

const deleteProjectQuery = `DELETE FROM projects WHERE id = ?`
const deleteProjectTechStackQuery = `DELETE FROM project_tech_stack WHERE project_id = ?`
const insertProjectTechStackQuery = `INSERT INTO project_tech_stack (project_id, label, sort_order) VALUES (?, ?, ?)`

const projectTechStackQuery = `
SELECT
	label
FROM project_tech_stack
WHERE project_id = ?
ORDER BY sort_order, id`

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

type techStackRow struct {
	Label sql.NullString `db:"label"`
}

func (r *projectRepository) ListProjects(ctx context.Context) ([]model.Project, error) {
	var rows []projectRow
	if err := r.db.SelectContext(ctx, &rows, listProjectsQuery); err != nil {
		return nil, fmt.Errorf("select projects: %w", err)
	}

	projects := make([]model.Project, 0, len(rows))
	for _, row := range rows {
		stack, err := r.fetchTechStack(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		projects = append(projects, model.Project{
			ID:          row.ID,
			Title:       toLocalizedText(row.TitleJA, row.TitleEN),
			Description: toLocalizedText(row.DescriptionJA, row.DescriptionEN),
			TechStack:   stack,
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
		stack, err := r.fetchTechStack(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, mapProjectRow(row, stack))
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

	stack, err := r.fetchTechStack(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	project := mapProjectRow(row, stack)
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

	if err = r.replaceTechStackTx(ctx, tx, projectID, project.TechStack); err != nil {
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

	if err = r.replaceTechStackTx(ctx, tx, project.ID, project.TechStack); err != nil {
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

	if _, execErr := tx.ExecContext(ctx, deleteProjectTechStackQuery, id); execErr != nil {
		err = fmt.Errorf("delete project tech stack %d: %w", id, execErr)
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

func (r *projectRepository) fetchTechStack(ctx context.Context, projectID int64) ([]string, error) {
	var stackRows []techStackRow
	if err := r.db.SelectContext(ctx, &stackRows, projectTechStackQuery, projectID); err != nil {
		return nil, fmt.Errorf("select project tech stack (id=%d): %w", projectID, err)
	}

	stack := make([]string, 0, len(stackRows))
	for _, sr := range stackRows {
		if sr.Label.Valid {
			stack = append(stack, sr.Label.String)
		}
	}
	return stack, nil
}

func (r *projectRepository) replaceTechStackTx(ctx context.Context, tx *sqlx.Tx, projectID int64, stack []string) error {
	if _, err := tx.ExecContext(ctx, deleteProjectTechStackQuery, projectID); err != nil {
		return fmt.Errorf("clear project tech stack %d: %w", projectID, err)
	}

	for index, label := range stack {
		if _, err := tx.ExecContext(ctx, insertProjectTechStackQuery, projectID, label, index); err != nil {
			return fmt.Errorf("insert project tech stack %d (%s): %w", projectID, label, err)
		}
	}
	return nil
}

func mapProjectRow(row projectRow, stack []string) model.AdminProject {
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
		TechStack:   append([]string(nil), stack...),
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
