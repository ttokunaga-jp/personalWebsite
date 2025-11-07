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

	"cloud.google.com/go/firestore"
	_ "github.com/go-sql-driver/mysql"

	firestoredb "github.com/takumi/personal-website/internal/infra/firestore"
)

type options struct {
	dsn     string
	project string
	prefix  string
	dryRun  bool
}

type catalogEntry struct {
	ID          uint64
	DisplayName string
}

type projectMembership struct {
	TechID      uint64
	DisplayName string
	SortOrder   int
}

func main() {
	opts := parseOptions()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := sql.Open("mysql", opts.dsn)
	if err != nil {
		log.Fatalf("failed to open mysql connection: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	catalog, err := loadTechCatalog(ctx, db)
	if err != nil {
		log.Fatalf("load tech catalog: %v", err)
	}
	log.Printf("loaded %d tech catalog entries", len(catalog))

	projectTech, err := loadProjectTech(ctx, db, catalog)
	if err != nil {
		log.Fatalf("load project tech stack: %v", err)
	}
	log.Printf("discovered %d projects with legacy tech stacks", len(projectTech))

	if len(projectTech) == 0 {
		log.Println("no legacy project tech stacks found; nothing to migrate")
		return
	}

	client, err := firestore.NewClient(ctx, opts.project)
	if err != nil {
		log.Fatalf("firestore client init failed: %v", err)
	}
	defer client.Close()

	if err := backfillFirestore(ctx, client, opts.prefix, projectTech, opts.dryRun); err != nil {
		log.Fatalf("firestore backfill failed: %v", err)
	}

	log.Println("project tech backfill completed successfully")
}

func parseOptions() *options {
	var (
		dsn     string
		project string
		prefix  string
		dryRun  bool
	)

	flag.StringVar(&dsn, "dsn", "", "MySQL DSN e.g. user:pass@tcp(localhost:3306)/portfolio?parseTime=true")
	flag.StringVar(&project, "project", "", "GCP project ID for Firestore")
	flag.StringVar(&prefix, "prefix", "", "Firestore collection prefix (environment)")
	flag.BoolVar(&dryRun, "dry-run", false, "Log intended mutations without executing them")
	flag.Parse()

	if dsn == "" {
		dsn = os.Getenv("APP_DATABASE_DSN")
	}
	if dsn == "" {
		log.Fatal("missing MySQL DSN: provide --dsn or APP_DATABASE_DSN")
	}
	if !strings.Contains(dsn, "parseTime") {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=true"
		} else {
			dsn += "?parseTime=true"
		}
	}

	if project == "" {
		project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if project == "" {
		log.Fatal("missing GCP project: provide --project or GOOGLE_CLOUD_PROJECT")
	}

	return &options{
		dsn:     dsn,
		project: project,
		prefix:  prefix,
		dryRun:  dryRun,
	}
}

func loadTechCatalog(ctx context.Context, db *sql.DB) (map[string]catalogEntry, error) {
	const query = `
SELECT
    id,
    display_name
FROM tech_catalog`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query tech_catalog: %w", err)
	}
	defer rows.Close()

	result := make(map[string]catalogEntry)
	for rows.Next() {
		var entry catalogEntry
		if err := rows.Scan(&entry.ID, &entry.DisplayName); err != nil {
			return nil, fmt.Errorf("scan tech_catalog: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(entry.DisplayName))
		if key == "" {
			continue
		}
		result[key] = entry
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tech_catalog: %w", err)
	}
	return result, nil
}

func loadProjectTech(ctx context.Context, db *sql.DB, catalog map[string]catalogEntry) (map[int64][]projectMembership, error) {
	const query = `
SELECT
    project_id,
    label,
    COALESCE(sort_order, 0) AS sort_order
FROM project_tech_stack
ORDER BY project_id, sort_order, label`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		if isTableMissing(err) {
			log.Println("project_tech_stack table not found; skipping legacy backfill")
			return map[int64][]projectMembership{}, nil
		}
		return nil, fmt.Errorf("query project_tech_stack: %w", err)
	}
	defer rows.Close()

	result := make(map[int64][]projectMembership)
	for rows.Next() {
		var (
			projectID int64
			label     string
			sortOrder int
		)
		if err := rows.Scan(&projectID, &label, &sortOrder); err != nil {
			return nil, fmt.Errorf("scan project_tech_stack: %w", err)
		}
		key := strings.ToLower(strings.TrimSpace(label))
		entry, ok := catalog[key]
		if !ok {
			log.Printf("warning: tech catalog entry not found for project %d label %q", projectID, label)
			continue
		}
		result[projectID] = append(result[projectID], projectMembership{
			TechID:      entry.ID,
			DisplayName: entry.DisplayName,
			SortOrder:   sortOrder,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project_tech_stack: %w", err)
	}
	return result, nil
}

func isTableMissing(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "doesn't exist")
}

func backfillFirestore(
	ctx context.Context,
	client *firestore.Client,
	prefix string,
	projectTech map[int64][]projectMembership,
	dryRun bool,
) error {
	base := firestoredb.CollectionName(prefix, "projects")
	now := time.Now().UTC()

	for projectID, memberships := range projectTech {
		if len(memberships) == 0 {
			continue
		}

		techDocs := make([]map[string]any, 0, len(memberships))
		labels := make([]string, 0, len(memberships))
		for idx, membership := range memberships {
			techDocs = append(techDocs, map[string]any{
				"id":        idx + 1,
				"techId":    membership.TechID,
				"context":   "primary",
				"note":      "",
				"sortOrder": membership.SortOrder,
			})
			labels = append(labels, membership.DisplayName)
		}

		docRef := client.Collection(base).Doc(fmt.Sprintf("%d", projectID))
		if dryRun {
			log.Printf("[dry-run] would update project %d: tech=%v techStack=%v", projectID, techDocs, labels)
			continue
		}

		_, err := docRef.Update(ctx, []firestore.Update{
			{Path: "tech", Value: techDocs},
			{Path: "techStack", Value: labels},
			{Path: "updatedAt", Value: now},
		})
		if err != nil {
			return fmt.Errorf("update firestore project %d: %w", projectID, err)
		}
		log.Printf("updated project %d tech relationships", projectID)
	}

	return nil
}
