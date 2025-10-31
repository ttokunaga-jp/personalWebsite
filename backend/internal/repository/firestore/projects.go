package firestore

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type projectRepository struct {
	base baseRepository
}

const projectsCollection = "projects"

type projectDocument struct {
	ID          int64        `firestore:"id"`
	Title       localizedDoc `firestore:"title"`
	Description localizedDoc `firestore:"description"`
	TechStack   []string     `firestore:"techStack"`
	LinkURL     string       `firestore:"linkUrl"`
	Year        int          `firestore:"year"`
	Published   bool         `firestore:"published"`
	SortOrder   *int         `firestore:"sortOrder,omitempty"`
	CreatedAt   time.Time    `firestore:"createdAt"`
	UpdatedAt   time.Time    `firestore:"updatedAt"`
}

func NewProjectRepository(client *firestore.Client, prefix string) repository.ProjectRepository {
	return &projectRepository{base: newBaseRepository(client, prefix)}
}

func (r *projectRepository) ListProjects(ctx context.Context) ([]model.Project, error) {
	docs, err := r.base.collection(projectsCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore projects: list: %w", err)
	}

	entries, err := r.decodeProjects(docs)
	if err != nil {
		return nil, err
	}

	var result []model.Project
	for _, entry := range entries {
		if !entry.Published {
			continue
		}
		result = append(result, model.Project{
			ID:          entry.ID,
			Title:       fromLocalizedDoc(entry.Title),
			Description: fromLocalizedDoc(entry.Description),
			TechStack:   copyStringSlice(entry.TechStack),
			LinkURL:     entry.LinkURL,
			Year:        entry.Year,
		})
	}

	return result, nil
}

func (r *projectRepository) ListAdminProjects(ctx context.Context) ([]model.AdminProject, error) {
	docs, err := r.base.collection(projectsCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore projects: list admin: %w", err)
	}

	entries, err := r.decodeProjects(docs)
	if err != nil {
		return nil, err
	}

	admin := make([]model.AdminProject, 0, len(entries))
	for _, entry := range entries {
		admin = append(admin, mapProjectDocument(entry))
	}
	return admin, nil
}

func (r *projectRepository) GetAdminProject(ctx context.Context, id int64) (*model.AdminProject, error) {
	doc, err := r.base.doc(projectsCollection, strconv.FormatInt(id, 10)).Get(ctx)
	if err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore projects: get %d: %w", id, err)
	}

	var entry projectDocument
	if err := doc.DataTo(&entry); err != nil {
		return nil, fmt.Errorf("firestore projects: decode %d: %w", id, err)
	}
	return ptr(mapProjectDocument(entry)), nil
}

func (r *projectRepository) CreateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error) {
	if project == nil {
		return nil, repository.ErrInvalidInput
	}

	id, err := nextID(ctx, r.base.client, r.base.prefix, projectsCollection)
	if err != nil {
		return nil, fmt.Errorf("firestore projects: next id: %w", err)
	}

	now := time.Now().UTC()
	entry := projectDocument{
		ID:          id,
		Title:       toLocalizedDoc(project.Title),
		Description: toLocalizedDoc(project.Description),
		TechStack:   copyStringSlice(project.TechStack),
		LinkURL:     project.LinkURL,
		Year:        project.Year,
		Published:   project.Published,
		SortOrder:   project.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	docRef := r.base.doc(projectsCollection, strconv.FormatInt(id, 10))
	if _, err := docRef.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("firestore projects: create %d: %w", id, err)
	}

	return r.GetAdminProject(ctx, id)
}

func (r *projectRepository) UpdateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error) {
	if project == nil {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(projectsCollection, strconv.FormatInt(project.ID, 10))
	now := time.Now().UTC()

	data := map[string]any{
		"title":       toLocalizedDoc(project.Title),
		"description": toLocalizedDoc(project.Description),
		"techStack":   copyStringSlice(project.TechStack),
		"linkUrl":     project.LinkURL,
		"year":        project.Year,
		"published":   project.Published,
		"sortOrder":   project.SortOrder,
		"updatedAt":   now,
	}

	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "title", Value: data["title"]},
		{Path: "description", Value: data["description"]},
		{Path: "techStack", Value: data["techStack"]},
		{Path: "linkUrl", Value: data["linkUrl"]},
		{Path: "year", Value: data["year"]},
		{Path: "published", Value: data["published"]},
		{Path: "sortOrder", Value: data["sortOrder"]},
		{Path: "updatedAt", Value: data["updatedAt"]},
	})
	if err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore projects: update %d: %w", project.ID, err)
	}

	return r.GetAdminProject(ctx, project.ID)
}

func (r *projectRepository) DeleteAdminProject(ctx context.Context, id int64) error {
	docRef := r.base.doc(projectsCollection, strconv.FormatInt(id, 10))
	_, err := docRef.Delete(ctx)
	if err != nil {
		if notFound(err) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("firestore projects: delete %d: %w", id, err)
	}
	return nil
}

func (r *projectRepository) decodeProjects(docs []*firestore.DocumentSnapshot) ([]projectDocument, error) {
	result := make([]projectDocument, 0, len(docs))
	for _, doc := range docs {
		var entry projectDocument
		if err := doc.DataTo(&entry); err != nil {
			return nil, fmt.Errorf("firestore projects: decode %s: %w", doc.Ref.ID, err)
		}
		if entry.ID == 0 {
			// Ensure ID consistency even if the field is missing.
			if id, err := strconv.ParseInt(doc.Ref.ID, 10, 64); err == nil {
				entry.ID = id
			}
		}
		result = append(result, entry)
	}

	sort.Slice(result, func(i, j int) bool {
		left, right := result[i], result[j]
		leftOrder := sortKey(left)
		rightOrder := sortKey(right)
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		if left.Year != right.Year {
			return left.Year > right.Year
		}
		return left.ID < right.ID
	})

	return result, nil
}

func sortKey(item projectDocument) int {
	if item.SortOrder != nil {
		return *item.SortOrder
	}
	return item.Year * 1000
}

func mapProjectDocument(doc projectDocument) model.AdminProject {
	return model.AdminProject{
		ID:          doc.ID,
		Title:       fromLocalizedDoc(doc.Title),
		Description: fromLocalizedDoc(doc.Description),
		TechStack:   copyStringSlice(doc.TechStack),
		LinkURL:     doc.LinkURL,
		Year:        doc.Year,
		Published:   doc.Published,
		SortOrder:   doc.SortOrder,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}

func ptr(value model.AdminProject) *model.AdminProject {
	return &value
}

var _ repository.ProjectRepository = (*projectRepository)(nil)
var _ repository.AdminProjectRepository = (*projectRepository)(nil)
