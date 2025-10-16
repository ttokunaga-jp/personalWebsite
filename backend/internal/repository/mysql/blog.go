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

type blogRepository struct {
	db *sqlx.DB
}

// NewBlogRepository returns a MySQL-backed blog repository.
func NewBlogRepository(db *sqlx.DB) repository.BlogRepository {
	return &blogRepository{db: db}
}

const listBlogPostsQuery = `
SELECT
	b.id,
	b.title_ja,
	b.title_en,
	b.summary_ja,
	b.summary_en,
	b.content_md_ja,
	b.content_md_en,
	b.published,
	b.published_at,
	b.created_at,
	b.updated_at
FROM blog_posts b
ORDER BY COALESCE(b.published_at, b.updated_at) DESC, b.id DESC`

const getBlogPostQuery = `
SELECT
	b.id,
	b.title_ja,
	b.title_en,
	b.summary_ja,
	b.summary_en,
	b.content_md_ja,
	b.content_md_en,
	b.published,
	b.published_at,
	b.created_at,
	b.updated_at
FROM blog_posts b
WHERE b.id = ?`

const insertBlogPostQuery = `
INSERT INTO blog_posts (
	title_ja,
	title_en,
	summary_ja,
	summary_en,
	content_md_ja,
	content_md_en,
	published,
	published_at,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

const updateBlogPostQuery = `
UPDATE blog_posts
SET
	title_ja = ?,
	title_en = ?,
	summary_ja = ?,
	summary_en = ?,
	content_md_ja = ?,
	content_md_en = ?,
	published = ?,
	published_at = ?,
	updated_at = NOW()
WHERE id = ?`

const deleteBlogPostQuery = `DELETE FROM blog_posts WHERE id = ?`

const listBlogTagsQuery = `SELECT tag FROM blog_post_tags WHERE post_id = ? ORDER BY sort_order, id`
const deleteBlogTagsQuery = `DELETE FROM blog_post_tags WHERE post_id = ?`
const insertBlogTagQuery = `INSERT INTO blog_post_tags (post_id, tag, sort_order) VALUES (?, ?, ?)`

type blogPostRow struct {
	ID          int64          `db:"id"`
	TitleJA     sql.NullString `db:"title_ja"`
	TitleEN     sql.NullString `db:"title_en"`
	SummaryJA   sql.NullString `db:"summary_ja"`
	SummaryEN   sql.NullString `db:"summary_en"`
	ContentJA   sql.NullString `db:"content_md_ja"`
	ContentEN   sql.NullString `db:"content_md_en"`
	Published   sql.NullBool   `db:"published"`
	PublishedAt sql.NullTime   `db:"published_at"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

func (r *blogRepository) ListBlogPosts(ctx context.Context) ([]model.BlogPost, error) {
	var rows []blogPostRow
	if err := r.db.SelectContext(ctx, &rows, listBlogPostsQuery); err != nil {
		return nil, fmt.Errorf("select blog posts: %w", err)
	}

	posts := make([]model.BlogPost, 0, len(rows))
	for _, row := range rows {
		tags, err := r.fetchTags(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, mapBlogPostRow(row, tags))
	}
	return posts, nil
}

func (r *blogRepository) GetBlogPost(ctx context.Context, id int64) (*model.BlogPost, error) {
	var row blogPostRow
	if err := r.db.GetContext(ctx, &row, getBlogPostQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get blog post %d: %w", id, err)
	}

	tags, err := r.fetchTags(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	post := mapBlogPostRow(row, tags)
	return &post, nil
}

func (r *blogRepository) CreateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error) {
	if post == nil {
		return nil, repository.ErrInvalidInput
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer rollbackOnError(tx, &err)

	res, execErr := tx.ExecContext(ctx, insertBlogPostQuery,
		post.Title.Ja,
		post.Title.En,
		post.Summary.Ja,
		post.Summary.En,
		post.ContentMD.Ja,
		post.ContentMD.En,
		post.Published,
		nullTime(post.PublishedAt),
	)
	if execErr != nil {
		err = fmt.Errorf("insert blog post: %w", execErr)
		return nil, err
	}

	postID, execErr := res.LastInsertId()
	if execErr != nil {
		err = fmt.Errorf("blog post last insert id: %w", execErr)
		return nil, err
	}

	if err = r.replaceTagsTx(ctx, tx, postID, post.Tags); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit blog post insert: %w", err)
	}

	created, err := r.GetBlogPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *blogRepository) UpdateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error) {
	if post == nil {
		return nil, repository.ErrInvalidInput
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer rollbackOnError(tx, &err)

	res, execErr := tx.ExecContext(ctx, updateBlogPostQuery,
		post.Title.Ja,
		post.Title.En,
		post.Summary.Ja,
		post.Summary.En,
		post.ContentMD.Ja,
		post.ContentMD.En,
		post.Published,
		nullTime(post.PublishedAt),
		post.ID,
	)
	if execErr != nil {
		err = fmt.Errorf("update blog post %d: %w", post.ID, execErr)
		return nil, err
	}

	affected, execErr := res.RowsAffected()
	if execErr != nil {
		err = fmt.Errorf("rows affected blog post %d: %w", post.ID, execErr)
		return nil, err
	}
	if affected == 0 {
		err = repository.ErrNotFound
		return nil, err
	}

	if err = r.replaceTagsTx(ctx, tx, post.ID, post.Tags); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit blog post update: %w", err)
	}

	updated, err := r.GetBlogPost(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *blogRepository) DeleteBlogPost(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer rollbackOnError(tx, &err)

	if _, execErr := tx.ExecContext(ctx, deleteBlogTagsQuery, id); execErr != nil {
		err = fmt.Errorf("delete blog tags %d: %w", id, execErr)
		return err
	}

	res, execErr := tx.ExecContext(ctx, deleteBlogPostQuery, id)
	if execErr != nil {
		err = fmt.Errorf("delete blog post %d: %w", id, execErr)
		return err
	}

	affected, execErr := res.RowsAffected()
	if execErr != nil {
		err = fmt.Errorf("rows affected blog post %d: %w", id, execErr)
		return err
	}
	if affected == 0 {
		err = repository.ErrNotFound
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit blog post delete: %w", err)
	}

	return nil
}

func (r *blogRepository) fetchTags(ctx context.Context, postID int64) ([]string, error) {
	var rows []struct {
		Tag sql.NullString `db:"tag"`
	}
	if err := r.db.SelectContext(ctx, &rows, listBlogTagsQuery, postID); err != nil {
		return nil, fmt.Errorf("select blog tags %d: %w", postID, err)
	}

	tags := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.Tag.Valid {
			tags = append(tags, row.Tag.String)
		}
	}
	return tags, nil
}

func (r *blogRepository) replaceTagsTx(ctx context.Context, tx *sqlx.Tx, postID int64, tags []string) error {
	if _, err := tx.ExecContext(ctx, deleteBlogTagsQuery, postID); err != nil {
		return fmt.Errorf("clear blog tags %d: %w", postID, err)
	}

	for idx, tag := range tags {
		if _, err := tx.ExecContext(ctx, insertBlogTagQuery, postID, tag, idx); err != nil {
			return fmt.Errorf("insert blog tag %d (%s): %w", postID, tag, err)
		}
	}
	return nil
}

func mapBlogPostRow(row blogPostRow, tags []string) model.BlogPost {
	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = time.Time{}
	}
	updatedAt := row.UpdatedAt.Time
	if !row.UpdatedAt.Valid {
		updatedAt = time.Time{}
	}

	return model.BlogPost{
		ID:          row.ID,
		Title:       toLocalizedText(row.TitleJA, row.TitleEN),
		Summary:     toLocalizedText(row.SummaryJA, row.SummaryEN),
		ContentMD:   toLocalizedText(row.ContentJA, row.ContentEN),
		Tags:        append([]string(nil), tags...),
		Published:   row.Published.Bool,
		PublishedAt: nullableTime(row.PublishedAt),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

var _ repository.BlogRepository = (*blogRepository)(nil)
