package admin

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// Service exposes administrative operations across content domains.
type Service interface {
	ListProjects(ctx context.Context) ([]model.AdminProject, error)
	GetProject(ctx context.Context, id int64) (*model.AdminProject, error)
	CreateProject(ctx context.Context, input ProjectInput) (*model.AdminProject, error)
	UpdateProject(ctx context.Context, id int64, input ProjectInput) (*model.AdminProject, error)
	DeleteProject(ctx context.Context, id int64) error

	ListResearch(ctx context.Context) ([]model.AdminResearch, error)
	GetResearch(ctx context.Context, id int64) (*model.AdminResearch, error)
	CreateResearch(ctx context.Context, input ResearchInput) (*model.AdminResearch, error)
	UpdateResearch(ctx context.Context, id int64, input ResearchInput) (*model.AdminResearch, error)
	DeleteResearch(ctx context.Context, id int64) error

	ListBlogPosts(ctx context.Context) ([]model.BlogPost, error)
	GetBlogPost(ctx context.Context, id int64) (*model.BlogPost, error)
	CreateBlogPost(ctx context.Context, input BlogPostInput) (*model.BlogPost, error)
	UpdateBlogPost(ctx context.Context, id int64, input BlogPostInput) (*model.BlogPost, error)
	DeleteBlogPost(ctx context.Context, id int64) error

	ListMeetings(ctx context.Context) ([]model.Meeting, error)
	GetMeeting(ctx context.Context, id int64) (*model.Meeting, error)
	CreateMeeting(ctx context.Context, input MeetingInput) (*model.Meeting, error)
	UpdateMeeting(ctx context.Context, id int64, input MeetingInput) (*model.Meeting, error)
	DeleteMeeting(ctx context.Context, id int64) error

	ListBlacklist(ctx context.Context) ([]model.BlacklistEntry, error)
	AddBlacklistEntry(ctx context.Context, input BlacklistInput) (*model.BlacklistEntry, error)
	RemoveBlacklistEntry(ctx context.Context, id int64) error
	IsEmailBlacklisted(ctx context.Context, email string) (bool, error)

	Summary(ctx context.Context) (*model.AdminSummary, error)
}

type service struct {
	projects  repository.AdminProjectRepository
	research  repository.AdminResearchRepository
	blogs     repository.BlogRepository
	meetings  repository.MeetingRepository
	blacklist repository.BlacklistRepository
}

// NewService wires repositories into the admin service.
func NewService(
	projects repository.AdminProjectRepository,
	research repository.AdminResearchRepository,
	blogs repository.BlogRepository,
	meetings repository.MeetingRepository,
	blacklist repository.BlacklistRepository,
) (Service, error) {
	if projects == nil || research == nil || blogs == nil || meetings == nil || blacklist == nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "admin service: missing dependencies", nil)
	}

	return &service{
		projects:  projects,
		research:  research,
		blogs:     blogs,
		meetings:  meetings,
		blacklist: blacklist,
	}, nil
}

// ProjectInput captures administrator-provided project data.
type ProjectInput struct {
	Title       model.LocalizedText
	Description model.LocalizedText
	TechStack   []string
	LinkURL     string
	Year        int
	Published   bool
	SortOrder   *int
}

func (s *service) ListProjects(ctx context.Context) ([]model.AdminProject, error) {
	return s.projects.ListAdminProjects(ctx)
}

func (s *service) GetProject(ctx context.Context, id int64) (*model.AdminProject, error) {
	return s.projects.GetAdminProject(ctx, id)
}

func (s *service) CreateProject(ctx context.Context, input ProjectInput) (*model.AdminProject, error) {
	if err := validateProjectInput(input); err != nil {
		return nil, err
	}

	project := model.AdminProject{
		Title:       normalizeLocalized(input.Title),
		Description: normalizeLocalized(input.Description),
		TechStack:   copyStrings(input.TechStack),
		LinkURL:     strings.TrimSpace(input.LinkURL),
		Year:        input.Year,
		Published:   input.Published,
		SortOrder:   copyIntPointer(input.SortOrder),
	}
	return s.projects.CreateAdminProject(ctx, &project)
}

func (s *service) UpdateProject(ctx context.Context, id int64, input ProjectInput) (*model.AdminProject, error) {
	if err := validateProjectInput(input); err != nil {
		return nil, err
	}

	project := model.AdminProject{
		ID:          id,
		Title:       normalizeLocalized(input.Title),
		Description: normalizeLocalized(input.Description),
		TechStack:   copyStrings(input.TechStack),
		LinkURL:     strings.TrimSpace(input.LinkURL),
		Year:        input.Year,
		Published:   input.Published,
		SortOrder:   copyIntPointer(input.SortOrder),
	}
	return s.projects.UpdateAdminProject(ctx, &project)
}

func (s *service) DeleteProject(ctx context.Context, id int64) error {
	return s.projects.DeleteAdminProject(ctx, id)
}

// ResearchInput captures administrator-provided research data.
type ResearchInput struct {
	Title     model.LocalizedText
	Summary   model.LocalizedText
	ContentMD model.LocalizedText
	Year      int
	Published bool
}

func (s *service) ListResearch(ctx context.Context) ([]model.AdminResearch, error) {
	return s.research.ListAdminResearch(ctx)
}

func (s *service) GetResearch(ctx context.Context, id int64) (*model.AdminResearch, error) {
	return s.research.GetAdminResearch(ctx, id)
}

func (s *service) CreateResearch(ctx context.Context, input ResearchInput) (*model.AdminResearch, error) {
	if err := validateResearchInput(input); err != nil {
		return nil, err
	}

	item := model.AdminResearch{
		Title:     normalizeLocalized(input.Title),
		Summary:   normalizeLocalized(input.Summary),
		ContentMD: normalizeLocalized(input.ContentMD),
		Year:      input.Year,
		Published: input.Published,
	}
	return s.research.CreateAdminResearch(ctx, &item)
}

func (s *service) UpdateResearch(ctx context.Context, id int64, input ResearchInput) (*model.AdminResearch, error) {
	if err := validateResearchInput(input); err != nil {
		return nil, err
	}

	item := model.AdminResearch{
		ID:        id,
		Title:     normalizeLocalized(input.Title),
		Summary:   normalizeLocalized(input.Summary),
		ContentMD: normalizeLocalized(input.ContentMD),
		Year:      input.Year,
		Published: input.Published,
	}
	return s.research.UpdateAdminResearch(ctx, &item)
}

func (s *service) DeleteResearch(ctx context.Context, id int64) error {
	return s.research.DeleteAdminResearch(ctx, id)
}

// BlogPostInput captures admin blog post data.
type BlogPostInput struct {
	Title       model.LocalizedText
	Summary     model.LocalizedText
	ContentMD   model.LocalizedText
	Tags        []string
	Published   bool
	PublishedAt *time.Time
}

func (s *service) ListBlogPosts(ctx context.Context) ([]model.BlogPost, error) {
	return s.blogs.ListBlogPosts(ctx)
}

func (s *service) GetBlogPost(ctx context.Context, id int64) (*model.BlogPost, error) {
	return s.blogs.GetBlogPost(ctx, id)
}

func (s *service) CreateBlogPost(ctx context.Context, input BlogPostInput) (*model.BlogPost, error) {
	if err := validateBlogPostInput(input); err != nil {
		return nil, err
	}

	post := model.BlogPost{
		Title:       normalizeLocalized(input.Title),
		Summary:     normalizeLocalized(input.Summary),
		ContentMD:   normalizeLocalized(input.ContentMD),
		Tags:        copyStrings(input.Tags),
		Published:   input.Published,
		PublishedAt: copyTimePointer(input.PublishedAt),
	}
	if post.Published && post.PublishedAt == nil {
		now := time.Now().UTC()
		post.PublishedAt = &now
	}
	return s.blogs.CreateBlogPost(ctx, &post)
}

func (s *service) UpdateBlogPost(ctx context.Context, id int64, input BlogPostInput) (*model.BlogPost, error) {
	if err := validateBlogPostInput(input); err != nil {
		return nil, err
	}

	post := model.BlogPost{
		ID:          id,
		Title:       normalizeLocalized(input.Title),
		Summary:     normalizeLocalized(input.Summary),
		ContentMD:   normalizeLocalized(input.ContentMD),
		Tags:        copyStrings(input.Tags),
		Published:   input.Published,
		PublishedAt: copyTimePointer(input.PublishedAt),
	}
	if post.Published && post.PublishedAt == nil {
		now := time.Now().UTC()
		post.PublishedAt = &now
	}
	return s.blogs.UpdateBlogPost(ctx, &post)
}

func (s *service) DeleteBlogPost(ctx context.Context, id int64) error {
	return s.blogs.DeleteBlogPost(ctx, id)
}

// MeetingInput captures reservation data edits.
type MeetingInput struct {
	Name            string
	Email           string
	Datetime        time.Time
	DurationMinutes int
	MeetURL         string
	Status          model.MeetingStatus
	Notes           string
}

func (s *service) ListMeetings(ctx context.Context) ([]model.Meeting, error) {
	return s.meetings.ListMeetings(ctx)
}

func (s *service) GetMeeting(ctx context.Context, id int64) (*model.Meeting, error) {
	return s.meetings.GetMeeting(ctx, id)
}

func (s *service) CreateMeeting(ctx context.Context, input MeetingInput) (*model.Meeting, error) {
	if err := validateMeetingInput(input); err != nil {
		return nil, err
	}

	meeting := model.Meeting{
		Name:            strings.TrimSpace(input.Name),
		Email:           strings.ToLower(strings.TrimSpace(input.Email)),
		Datetime:        input.Datetime.UTC(),
		DurationMinutes: input.DurationMinutes,
		MeetURL:         strings.TrimSpace(input.MeetURL),
		Status:          input.Status,
		Notes:           strings.TrimSpace(input.Notes),
	}
	return s.meetings.CreateMeeting(ctx, &meeting)
}

func (s *service) UpdateMeeting(ctx context.Context, id int64, input MeetingInput) (*model.Meeting, error) {
	if err := validateMeetingInput(input); err != nil {
		return nil, err
	}

	meeting := model.Meeting{
		ID:              id,
		Name:            strings.TrimSpace(input.Name),
		Email:           strings.ToLower(strings.TrimSpace(input.Email)),
		Datetime:        input.Datetime.UTC(),
		DurationMinutes: input.DurationMinutes,
		MeetURL:         strings.TrimSpace(input.MeetURL),
		Status:          input.Status,
		Notes:           strings.TrimSpace(input.Notes),
	}
	return s.meetings.UpdateMeeting(ctx, &meeting)
}

func (s *service) DeleteMeeting(ctx context.Context, id int64) error {
	return s.meetings.DeleteMeeting(ctx, id)
}

// BlacklistInput captures email blocking requests.
type BlacklistInput struct {
	Email  string
	Reason string
}

func (s *service) ListBlacklist(ctx context.Context) ([]model.BlacklistEntry, error) {
	return s.blacklist.ListBlacklistEntries(ctx)
}

func (s *service) AddBlacklistEntry(ctx context.Context, input BlacklistInput) (*model.BlacklistEntry, error) {
	if err := validateBlacklistInput(input); err != nil {
		return nil, err
	}

	entry := model.BlacklistEntry{
		Email:  strings.ToLower(strings.TrimSpace(input.Email)),
		Reason: strings.TrimSpace(input.Reason),
	}
	result, err := s.blacklist.AddBlacklistEntry(ctx, &entry)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, errs.New(errs.CodeConflict, http.StatusConflict, "email already blacklisted", err)
		}
		return nil, err
	}
	return result, nil
}

func (s *service) RemoveBlacklistEntry(ctx context.Context, id int64) error {
	return s.blacklist.RemoveBlacklistEntry(ctx, id)
}

func (s *service) IsEmailBlacklisted(ctx context.Context, email string) (bool, error) {
	_, err := s.blacklist.FindBlacklistEntryByEmail(ctx, email)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (s *service) Summary(ctx context.Context) (*model.AdminSummary, error) {
	projects, err := s.projects.ListAdminProjects(ctx)
	if err != nil {
		return nil, err
	}
	var publishedProjects, draftProjects int
	for _, p := range projects {
		if p.Published {
			publishedProjects++
		} else {
			draftProjects++
		}
	}

	research, err := s.research.ListAdminResearch(ctx)
	if err != nil {
		return nil, err
	}
	var publishedResearch, draftResearch int
	for _, rch := range research {
		if rch.Published {
			publishedResearch++
		} else {
			draftResearch++
		}
	}

	posts, err := s.blogs.ListBlogPosts(ctx)
	if err != nil {
		return nil, err
	}
	var publishedBlogs, draftBlogs int
	for _, post := range posts {
		if post.Published {
			publishedBlogs++
		} else {
			draftBlogs++
		}
	}

	meetings, err := s.meetings.ListMeetings(ctx)
	if err != nil {
		return nil, err
	}
	var pendingMeetings int
	for _, meeting := range meetings {
		if meeting.Status == model.MeetingStatusPending {
			pendingMeetings++
		}
	}

	blacklist, err := s.blacklist.ListBlacklistEntries(ctx)
	if err != nil {
		return nil, err
	}

	return &model.AdminSummary{
		PublishedProjects: publishedProjects,
		DraftProjects:     draftProjects,
		PublishedResearch: publishedResearch,
		DraftResearch:     draftResearch,
		PublishedBlogs:    publishedBlogs,
		DraftBlogs:        draftBlogs,
		PendingMeetings:   pendingMeetings,
		BlacklistEntries:  len(blacklist),
	}, nil
}

func validateProjectInput(input ProjectInput) error {
	if strings.TrimSpace(input.Title.Ja) == "" && strings.TrimSpace(input.Title.En) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "project title is required", nil)
	}
	if input.Year <= 0 {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "project year must be positive", nil)
	}
	return nil
}

func validateResearchInput(input ResearchInput) error {
	if strings.TrimSpace(input.Title.Ja) == "" && strings.TrimSpace(input.Title.En) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "research title is required", nil)
	}
	if input.Year <= 0 {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "research year must be positive", nil)
	}
	return nil
}

func validateBlogPostInput(input BlogPostInput) error {
	if strings.TrimSpace(input.Title.Ja) == "" && strings.TrimSpace(input.Title.En) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "blog post title is required", nil)
	}
	if strings.TrimSpace(input.ContentMD.Ja) == "" && strings.TrimSpace(input.ContentMD.En) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "blog post content is required", nil)
	}
	return nil
}

func validateMeetingInput(input MeetingInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "meeting name is required", nil)
	}
	if strings.TrimSpace(input.Email) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "meeting email is required", nil)
	}
	if input.DurationMinutes <= 0 {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "meeting duration must be positive", nil)
	}
	if input.Status != model.MeetingStatusPending && input.Status != model.MeetingStatusConfirmed && input.Status != model.MeetingStatusCancelled {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid meeting status", nil)
	}
	return nil
}

func validateBlacklistInput(input BlacklistInput) error {
	if strings.TrimSpace(input.Email) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "email is required", nil)
	}
	return nil
}

func normalizeLocalized(text model.LocalizedText) model.LocalizedText {
	return model.LocalizedText{
		Ja: strings.TrimSpace(text.Ja),
		En: strings.TrimSpace(text.En),
	}
}

func copyStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = strings.TrimSpace(v)
	}
	return result
}

func copyIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func copyTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	clone := value.UTC()
	return &clone
}
