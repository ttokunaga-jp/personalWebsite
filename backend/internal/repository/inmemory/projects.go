package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type projectRepository struct {
	mu       sync.RWMutex
	projects []model.AdminProject
	nextID   int64
}

func NewProjectRepository() repository.ProjectRepository {
	return newProjectRepository()
}

func newProjectRepository() *projectRepository {
	maxID := int64(0)
	projects := make([]model.AdminProject, len(defaultAdminProjects))
	copy(projects, defaultAdminProjects)
	for _, p := range projects {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	return &projectRepository{
		projects: projects,
		nextID:   maxID + 1,
	}
}

func (r *projectRepository) ListProjects(ctx context.Context) ([]model.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Project
	for _, p := range r.projects {
		if !p.Published {
			continue
		}
		result = append(result, model.Project{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Tech:        cloneProjectTechMemberships(p.Tech),
			TechStack:   techDisplayNamesFromMemberships(p.Tech),
			LinkURL:     p.LinkURL,
			Year:        p.Year,
		})
	}
	return result, nil
}

func (r *projectRepository) ListAdminProjects(ctx context.Context) ([]model.AdminProject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projects := make([]model.AdminProject, len(r.projects))
	for i, p := range r.projects {
		projects[i] = copyAdminProject(p)
	}
	return projects, nil
}

func (r *projectRepository) GetAdminProject(ctx context.Context, id int64) (*model.AdminProject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.projects {
		if p.ID == id {
			proj := copyAdminProject(p)
			return &proj, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *projectRepository) CreateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error) {
	if project == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	project.ID = r.nextID
	r.nextID++
	project.CreatedAt = now
	project.UpdatedAt = now

	r.projects = append(r.projects, copyAdminProject(*project))
	created := copyAdminProject(*project)
	return &created, nil
}

func (r *projectRepository) UpdateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error) {
	if project == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, existing := range r.projects {
		if existing.ID == project.ID {
			project.CreatedAt = existing.CreatedAt
			project.UpdatedAt = time.Now().UTC()
			r.projects[idx] = copyAdminProject(*project)
			updated := copyAdminProject(*project)
			return &updated, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *projectRepository) DeleteAdminProject(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, p := range r.projects {
		if p.ID == id {
			r.projects = append(r.projects[:idx], r.projects[idx+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func copyAdminProject(src model.AdminProject) model.AdminProject {
	dst := src
	if src.SortOrder != nil {
		value := *src.SortOrder
		dst.SortOrder = &value
	}
	dst.Tech = cloneProjectTechMemberships(src.Tech)
	return dst
}

func cloneProjectTechMemberships(src []model.TechMembership) []model.TechMembership {
	if len(src) == 0 {
		return nil
	}
	result := make([]model.TechMembership, len(src))
	for i, membership := range src {
		result[i] = model.TechMembership{
			MembershipID: membership.MembershipID,
			EntityType:   membership.EntityType,
			EntityID:     membership.EntityID,
			Tech: model.TechCatalogEntry{
				ID:          membership.Tech.ID,
				Slug:        membership.Tech.Slug,
				DisplayName: membership.Tech.DisplayName,
				Category:    membership.Tech.Category,
				Level:       membership.Tech.Level,
				Icon:        membership.Tech.Icon,
				SortOrder:   membership.Tech.SortOrder,
				Active:      membership.Tech.Active,
				CreatedAt:   membership.Tech.CreatedAt,
				UpdatedAt:   membership.Tech.UpdatedAt,
			},
			Context:   membership.Context,
			Note:      membership.Note,
			SortOrder: membership.SortOrder,
		}
	}
	return result
}

func techDisplayNamesFromMemberships(tech []model.TechMembership) []string {
	if len(tech) == 0 {
		return nil
	}
	names := make([]string, 0, len(tech))
	for _, membership := range tech {
		if membership.Tech.DisplayName == "" {
			continue
		}
		names = append(names, membership.Tech.DisplayName)
	}
	if len(names) == 0 {
		return nil
	}
	return names
}

var _ repository.ProjectRepository = (*projectRepository)(nil)
var _ repository.AdminProjectRepository = (*projectRepository)(nil)
