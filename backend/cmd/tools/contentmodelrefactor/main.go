package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type options struct {
	dsn    string
	dryRun bool
}

type migrationStep struct {
	name string
	fn   func(context.Context, *sql.DB, *options) error
}

func main() {
	opts := parseOptions()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	db, err := sql.Open("mysql", opts.dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	steps := []migrationStep{
		{name: "profiles", fn: migrateProfiles},
		{name: "tech_catalog", fn: migrateTechCatalog},
		{name: "projects", fn: migrateProjects},
		{name: "research_blog_entries", fn: migrateResearchAndBlog},
		{name: "meeting_reservations", fn: migrateMeetings},
	}

	for _, step := range steps {
		log.Printf("=== running step: %s (dryRun=%t) ===", step.name, opts.dryRun)
		if err := step.fn(ctx, db, opts); err != nil {
			log.Fatalf("step %s failed: %v", step.name, err)
		}
	}

	log.Println("migration scaffolding completed successfully")
}

func parseOptions() *options {
	var (
		dsn    string
		dryRun bool
	)

	flag.StringVar(&dsn, "dsn", "", "MySQL DSN (e.g. user:pass@tcp(localhost:3306)/portfolio?parseTime=true)")
	flag.BoolVar(&dryRun, "dry-run", false, "Simulate migration without mutating data")
	flag.Parse()

	if dsn == "" {
		dsn = os.Getenv("APP_DATABASE_DSN")
	}
	if dsn == "" {
		log.Fatal("dsn is required (provide --dsn or APP_DATABASE_DSN)")
	}

	if !strings.Contains(dsn, "parseTime") {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=true"
		} else {
			dsn += "?parseTime=true"
		}
	}

	return &options{
		dsn:    dsn,
		dryRun: dryRun,
	}
}

func migrateProfiles(ctx context.Context, db *sql.DB, opts *options) error {
	const countQuery = `SELECT COUNT(*) FROM legacy_profile`
	var count int
	if err := db.QueryRowContext(ctx, countQuery).Scan(&count); err != nil {
		if isTableMissing(err) {
			log.Println("legacy_profile not found, skipping profile migration")
			return nil
		}
		return fmt.Errorf("count legacy_profile: %w", err)
	}
	log.Printf("legacy_profile rows: %d", count)
	if opts.dryRun || count == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertProfile := `
INSERT INTO profiles (
  id, display_name, headline_ja, headline_en, summary_ja, summary_en,
  avatar_url, location_ja, location_en, theme_mode, theme_accent_color,
  lab_name_ja, lab_name_en, lab_advisor_ja, lab_advisor_en, lab_room_ja, lab_room_en, lab_url,
  created_at, updated_at
) SELECT
  id,
  COALESCE(NULLIF(name_ja, ''), NULLIF(name_en, ''), CONCAT('Profile#', id)),
  title_ja,
  title_en,
  summary_ja,
  summary_en,
  NULL,
  affiliation_ja,
  affiliation_en,
  'system',
  NULL,
  lab_ja,
  lab_en,
  NULL,
  NULL,
  NULL,
  NULL,
  NULL,
  NOW(3),
  NOW(3)
FROM legacy_profile
ON DUPLICATE KEY UPDATE
  display_name = VALUES(display_name),
  summary_ja = VALUES(summary_ja),
  summary_en = VALUES(summary_en),
  updated_at = VALUES(updated_at);`

	if _, err = tx.ExecContext(ctx, insertProfile); err != nil {
		return fmt.Errorf("insert profiles: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit profile migration: %w", err)
	}

	return nil
}

func migrateTechCatalog(ctx context.Context, db *sql.DB, opts *options) error {
	const legacySkillCount = `SELECT COUNT(*) FROM legacy_profile_skills`
	var skillCount int
	if err := db.QueryRowContext(ctx, legacySkillCount).Scan(&skillCount); err != nil {
		if isTableMissing(err) {
			log.Println("legacy_profile_skills not found, skipping tech catalog scaffolding")
			return nil
		}
		return fmt.Errorf("count legacy_profile_skills: %w", err)
	}
	log.Printf("legacy_profile_skills rows: %d", skillCount)

	if opts.dryRun || skillCount == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertCatalog := `
INSERT IGNORE INTO tech_catalog (slug, display_name, level, sort_order, created_at, updated_at)
SELECT
  LOWER(REPLACE(REPLACE(REPLACE(skill_en, ' ', '-'), '/', '-'), '#', 'sharp')),
  COALESCE(NULLIF(skill_en, ''), skill_ja),
  'intermediate',
  COALESCE(sort_order, 0),
  NOW(3),
  NOW(3)
FROM legacy_profile_skills
WHERE COALESCE(NULLIF(skill_en, ''), NULLIF(skill_ja, '')) IS NOT NULL;`

	if _, err = tx.ExecContext(ctx, insertCatalog); err != nil {
		return fmt.Errorf("insert tech_catalog: %w", err)
	}

	createSection := `
INSERT INTO profile_tech_sections (profile_id, title_ja, title_en, layout, breakpoint, sort_order)
SELECT DISTINCT
  lps.profile_id,
  'スキルセット',
  'Skills',
  'chips',
  'lg',
  0
FROM legacy_profile_skills lps
LEFT JOIN profile_tech_sections pts ON pts.profile_id = lps.profile_id
WHERE pts.id IS NULL;`

	if _, err = tx.ExecContext(ctx, createSection); err != nil {
		return fmt.Errorf("insert profile_tech_sections: %w", err)
	}

	linkTech := `
INSERT IGNORE INTO tech_relationships (entity_type, entity_id, tech_id, context, sort_order, created_at)
SELECT
  'profile_section',
  pts.id,
  tc.id,
  'primary',
  lps.sort_order,
  NOW(3)
FROM legacy_profile_skills lps
JOIN profile_tech_sections pts ON pts.profile_id = lps.profile_id
JOIN tech_catalog tc ON tc.display_name = COALESCE(NULLIF(lps.skill_en, ''), lps.skill_ja);`

	if _, err = tx.ExecContext(ctx, linkTech); err != nil {
		return fmt.Errorf("link tech to profile sections: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tech catalog migration: %w", err)
	}

	return nil
}

func migrateProjects(ctx context.Context, db *sql.DB, opts *options) error {
	const legacyProjectCount = `SELECT COUNT(*) FROM legacy_projects`
	var count int
	if err := db.QueryRowContext(ctx, legacyProjectCount).Scan(&count); err != nil {
		if isTableMissing(err) {
			log.Println("legacy_projects not found, skipping project migration")
			return nil
		}
		return fmt.Errorf("count legacy_projects: %w", err)
	}
	log.Printf("legacy_projects rows: %d", count)
	if opts.dryRun || count == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertProjects := `
INSERT INTO projects (
  id, slug, title_ja, title_en, summary_ja, summary_en, description_ja, description_en,
  primary_link_url, created_at, updated_at, published, highlight, sort_order
)
SELECT
  id,
  IFNULL(
    NULLIF(LOWER(REPLACE(REPLACE(title_en, ' ', '-'), '/', '-')), ''),
    CONCAT('project-', id)
  ),
  title_ja,
  title_en,
  description_ja,
  description_en,
  description_ja,
  description_en,
  link_url,
  created_at,
  updated_at,
  published,
  0,
  COALESCE(sort_order, 0)
FROM legacy_projects
ON DUPLICATE KEY UPDATE
  title_ja = VALUES(title_ja),
  title_en = VALUES(title_en),
  summary_ja = VALUES(summary_ja),
  summary_en = VALUES(summary_en),
  updated_at = VALUES(updated_at);`

	if _, err = tx.ExecContext(ctx, insertProjects); err != nil {
		return fmt.Errorf("insert projects: %w", err)
	}

	linkProjectTech := `
INSERT IGNORE INTO tech_relationships (entity_type, entity_id, tech_id, context, sort_order, created_at)
SELECT
  'project',
  p.id,
  tc.id,
  'primary',
  pts.sort_order,
  NOW(3)
FROM legacy_project_tech_stack pts
JOIN projects p ON p.id = pts.project_id
JOIN tech_catalog tc ON tc.display_name = pts.label;`

	if _, err = tx.ExecContext(ctx, linkProjectTech); err != nil {
		return fmt.Errorf("link project tech: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit project migration: %w", err)
	}

	return nil
}

func migrateResearchAndBlog(ctx context.Context, db *sql.DB, opts *options) error {
	var (
		researchCount int
		blogCount     int
	)
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM legacy_research").Scan(&researchCount); err != nil {
		if !isTableMissing(err) {
			return fmt.Errorf("count legacy_research: %w", err)
		}
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM legacy_blog_posts").Scan(&blogCount); err != nil {
		if !isTableMissing(err) {
			return fmt.Errorf("count legacy_blog_posts: %w", err)
		}
	}
	log.Printf("legacy_research rows: %d, legacy_blog_posts rows: %d", researchCount, blogCount)

	if opts.dryRun || (researchCount+blogCount) == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertResearch := `
INSERT INTO research_blog_entries (
  slug, kind, title_ja, title_en, overview_ja, overview_en,
  external_url, published_at, created_at, updated_at, is_draft
)
SELECT
  LOWER(REPLACE(REPLACE(title_en, ' ', '-'), '/', '-')),
  'research',
  title_ja,
  title_en,
  summary_ja,
  summary_en,
  '' AS external_url,
  CONCAT(year, '-01-01 00:00:00'),
  created_at,
  updated_at,
  NOT published
FROM legacy_research;`

	if _, err = tx.ExecContext(ctx, insertResearch); err != nil {
		return fmt.Errorf("insert research entries: %w", err)
	}

	insertBlog := `
INSERT INTO research_blog_entries (
  slug, kind, title_ja, title_en, overview_ja, overview_en,
  external_url, published_at, created_at, updated_at, is_draft
)
SELECT
  LOWER(REPLACE(REPLACE(title_en, ' ', '-'), '/', '-')),
  'blog',
  title_ja,
  title_en,
  summary_ja,
  summary_en,
  '' AS external_url,
  COALESCE(published_at, created_at),
  created_at,
  updated_at,
  NOT published
FROM legacy_blog_posts;`

	if _, err = tx.ExecContext(ctx, insertBlog); err != nil {
		return fmt.Errorf("insert blog entries: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit research/blog migration: %w", err)
	}

	return nil
}

func migrateMeetings(ctx context.Context, db *sql.DB, opts *options) error {
	const legacyMeetingCount = `SELECT COUNT(*) FROM legacy_meetings`
	var count int
	if err := db.QueryRowContext(ctx, legacyMeetingCount).Scan(&count); err != nil {
		if isTableMissing(err) {
			log.Println("legacy_meetings not found, skipping meeting migration")
			return nil
		}
		return fmt.Errorf("count legacy_meetings: %w", err)
	}
	log.Printf("legacy_meetings rows: %d", count)
	if opts.dryRun || count == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertMeetings := `
INSERT INTO meeting_reservations (
  id, name, email, topic, message, start_at, end_at, duration_minutes,
  google_event_id, google_calendar_status, status, confirmation_sent_at,
  last_notification_sent_at, lookup_hash, created_at, updated_at
)
SELECT
  id,
  name,
  email,
  topic,
  notes,
  meeting_at,
  DATE_ADD(meeting_at, INTERVAL duration_minutes MINUTE),
  duration_minutes,
  calendar_event_id,
  status,
  status,
  updated_at,
  NULL,
  SHA2(CONCAT(LOWER(email), ':', LOWER(name)), 256),
  created_at,
  updated_at
FROM legacy_meetings
ON DUPLICATE KEY UPDATE
  google_event_id = VALUES(google_event_id),
  status = VALUES(status),
  updated_at = VALUES(updated_at);`

	if _, err = tx.ExecContext(ctx, insertMeetings); err != nil {
		return fmt.Errorf("insert meeting_reservations: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit meeting migration: %w", err)
	}

	return nil
}

func isTableMissing(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "doesn't exist")
}
