package firestore

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type projectRepository struct {
	base baseRepository
}

const (
	projectsCollection = "projects"
	projectEntityType  = "project"
)

type projectDocument struct {
	ID          int64            `firestore:"id"`
	Title       localizedDoc     `firestore:"title"`
	Description localizedDoc     `firestore:"description"`
	TechStack   []string         `firestore:"techStack"`
	Tech        []projectTechDoc `firestore:"tech"`
	LinkURL     string           `firestore:"linkUrl"`
	Year        int              `firestore:"year"`
	Published   bool             `firestore:"published"`
	SortOrder   *int             `firestore:"sortOrder,omitempty"`
	CreatedAt   time.Time        `firestore:"createdAt"`
	UpdatedAt   time.Time        `firestore:"updatedAt"`
}

type projectTechDoc struct {
	ID        uint64 `firestore:"id"`
	TechID    uint64 `firestore:"techId"`
	Context   string `firestore:"context"`
	Note      string `firestore:"note"`
	SortOrder int    `firestore:"sortOrder"`
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
		tech := mapProjectTech(entry)
		result = append(result, model.Project{
			ID:          entry.ID,
			Title:       fromLocalizedDoc(entry.Title),
			Description: fromLocalizedDoc(entry.Description),
			Tech:        tech,
			TechStack:   techDisplayNames(entry, tech),
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
		TechStack:   techStackFromMemberships(project.Tech),
		Tech:        toProjectTechDocs(project.Tech),
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
		"techStack":   techStackFromMemberships(project.Tech),
		"tech":        toProjectTechDocs(project.Tech),
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
		{Path: "tech", Value: data["tech"]},
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

func mapProjectTech(doc projectDocument) []model.TechMembership {
	if len(doc.Tech) == 0 {
		return nil
	}
	memberships := make([]model.TechMembership, 0, len(doc.Tech))
	for _, membership := range doc.Tech {
		memberships = append(memberships, model.TechMembership{
			MembershipID: membership.ID,
			EntityType:   projectEntityType,
			EntityID:     uint64(doc.ID),
			Tech: model.TechCatalogEntry{
				ID: membership.TechID,
			},
			Context:   model.TechContext(strings.TrimSpace(membership.Context)),
			Note:      strings.TrimSpace(membership.Note),
			SortOrder: membership.SortOrder,
		})
	}
	return memberships
}

func techDisplayNames(doc projectDocument, tech []model.TechMembership) []string {
	if len(doc.TechStack) > 0 {
		return trimStringSlice(doc.TechStack)
	}
	if len(tech) == 0 {
		return nil
	}
	names := make([]string, 0, len(tech))
	for _, membership := range tech {
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

func trimStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func techStackFromMemberships(tech []model.TechMembership) []string {
	if len(tech) == 0 {
		return nil
	}
	names := make([]string, 0, len(tech))
	for _, membership := range tech {
		if membership.Tech.DisplayName == "" {
			continue
		}
		if trimmed := strings.TrimSpace(membership.Tech.DisplayName); trimmed != "" {
			names = append(names, trimmed)
		}
	}
	if len(names) == 0 {
		return nil
	}
	return names
}

func toProjectTechDocs(tech []model.TechMembership) []projectTechDoc {
	if len(tech) == 0 {
		return nil
	}
	result := make([]projectTechDoc, 0, len(tech))
	nextID := uint64(1)
	for _, membership := range tech {
		if membership.Tech.ID == 0 {
			continue
		}
		id := membership.MembershipID
		if id == 0 {
			id = nextID
			nextID++
		}
		result = append(result, projectTechDoc{
			ID:        id,
			TechID:    membership.Tech.ID,
			Context:   string(membership.Context),
			Note:      strings.TrimSpace(membership.Note),
			SortOrder: membership.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func mapProjectDocument(doc projectDocument) model.AdminProject {
	tech := mapProjectTech(doc)
	return model.AdminProject{
		ID:          doc.ID,
		Title:       fromLocalizedDoc(doc.Title),
		Description: fromLocalizedDoc(doc.Description),
		Tech:        tech,
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
