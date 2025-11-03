package admin

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// Service exposes administrative operations across content domains.
type Service interface {
	GetProfile(ctx context.Context) (*model.AdminProfile, error)
	UpdateProfile(ctx context.Context, input ProfileInput) (*model.AdminProfile, error)

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

	ListContactMessages(ctx context.Context) ([]model.ContactMessage, error)
	GetContactMessage(ctx context.Context, id string) (*model.ContactMessage, error)
	UpdateContactMessage(ctx context.Context, id string, input ContactUpdateInput) (*model.ContactMessage, error)
	DeleteContactMessage(ctx context.Context, id string) error

	ListBlacklist(ctx context.Context) ([]model.BlacklistEntry, error)
	AddBlacklistEntry(ctx context.Context, input BlacklistInput) (*model.BlacklistEntry, error)
	UpdateBlacklistEntry(ctx context.Context, id int64, input BlacklistInput) (*model.BlacklistEntry, error)
	RemoveBlacklistEntry(ctx context.Context, id int64) error
	IsEmailBlacklisted(ctx context.Context, email string) (bool, error)

	Summary(ctx context.Context) (*model.AdminSummary, error)
}

type service struct {
	profile   repository.AdminProfileRepository
	projects  repository.AdminProjectRepository
	research  repository.AdminResearchRepository
	contacts  repository.AdminContactRepository
	blacklist repository.BlacklistRepository
}

// NewService wires repositories into the admin service.
func NewService(
	profile repository.AdminProfileRepository,
	projects repository.AdminProjectRepository,
	research repository.AdminResearchRepository,
	contacts repository.AdminContactRepository,
	blacklist repository.BlacklistRepository,
) (Service, error) {
	if profile == nil || projects == nil || research == nil || contacts == nil || blacklist == nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "admin service: missing dependencies", nil)
	}

	return &service{
		profile:   profile,
		projects:  projects,
		research:  research,
		contacts:  contacts,
		blacklist: blacklist,
	}, nil
}

// ProfileInput captures administrator-provided profile data.
type ProfileInput struct {
	Name        model.LocalizedText
	Title       model.LocalizedText
	Affiliation model.LocalizedText
	Lab         model.LocalizedText
	Summary     model.LocalizedText
	Skills      []model.LocalizedText
	FocusAreas  []model.LocalizedText
}

func (s *service) GetProfile(ctx context.Context) (*model.AdminProfile, error) {
	return s.profile.GetAdminProfile(ctx)
}

func (s *service) UpdateProfile(ctx context.Context, input ProfileInput) (*model.AdminProfile, error) {
	if err := validateProfileInput(input); err != nil {
		return nil, err
	}

	adminProfile := model.AdminProfile{
		Name:        normalizeLocalized(input.Name),
		Title:       normalizeLocalized(input.Title),
		Affiliation: normalizeLocalized(input.Affiliation),
		Lab:         normalizeLocalized(input.Lab),
		Summary:     normalizeLocalized(input.Summary),
		Skills:      normalizeLocalizedList(input.Skills),
		FocusAreas:  normalizeLocalizedList(input.FocusAreas),
	}

	return s.profile.UpdateAdminProfile(ctx, &adminProfile)
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

// ContactUpdateInput captures moderation edits for a contact submission.
type ContactUpdateInput struct {
	Topic     string
	Message   string
	Status    model.ContactStatus
	AdminNote string
}

func (s *service) ListContactMessages(ctx context.Context) ([]model.ContactMessage, error) {
	return s.contacts.ListContactMessages(ctx)
}

func (s *service) GetContactMessage(ctx context.Context, id string) (*model.ContactMessage, error) {
	return s.contacts.GetContactMessage(ctx, strings.TrimSpace(id))
}

func (s *service) UpdateContactMessage(ctx context.Context, id string, input ContactUpdateInput) (*model.ContactMessage, error) {
	if err := validateContactUpdateInput(input); err != nil {
		return nil, err
	}

	message := &model.ContactMessage{
		ID:        strings.TrimSpace(id),
		Topic:     strings.TrimSpace(input.Topic),
		Message:   strings.TrimSpace(input.Message),
		Status:    input.Status,
		AdminNote: strings.TrimSpace(input.AdminNote),
	}
	return s.contacts.UpdateContactMessage(ctx, message)
}

func (s *service) DeleteContactMessage(ctx context.Context, id string) error {
	return s.contacts.DeleteContactMessage(ctx, strings.TrimSpace(id))
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

func (s *service) UpdateBlacklistEntry(ctx context.Context, id int64, input BlacklistInput) (*model.BlacklistEntry, error) {
	if err := validateBlacklistInput(input); err != nil {
		return nil, err
	}

	entry := model.BlacklistEntry{
		ID:     id,
		Email:  strings.ToLower(strings.TrimSpace(input.Email)),
		Reason: strings.TrimSpace(input.Reason),
	}
	updated, err := s.blacklist.UpdateBlacklistEntry(ctx, &entry)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, errs.New(errs.CodeConflict, http.StatusConflict, "email already blacklisted", err)
		}
		return nil, err
	}
	return updated, nil
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
	profile, err := s.profile.GetAdminProfile(ctx)
	if err != nil {
		return nil, err
	}

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
	for _, item := range research {
		if item.Published {
			publishedResearch++
		} else {
			draftResearch++
		}
	}

	contacts, err := s.contacts.ListContactMessages(ctx)
	if err != nil {
		return nil, err
	}
	var pendingContacts int
	for _, contact := range contacts {
		if contact.Status == model.ContactStatusPending {
			pendingContacts++
		}
	}

	blacklist, err := s.blacklist.ListBlacklistEntries(ctx)
	if err != nil {
		return nil, err
	}

	summary := &model.AdminSummary{
		PublishedProjects: publishedProjects,
		DraftProjects:     draftProjects,
		PublishedResearch: publishedResearch,
		DraftResearch:     draftResearch,
		PendingContacts:   pendingContacts,
		BlacklistEntries:  len(blacklist),
	}
	if profile != nil {
		summary.SkillCount = len(profile.Skills)
		summary.FocusAreaCount = len(profile.FocusAreas)
		summary.ProfileUpdatedAt = profile.UpdatedAt
	}
	return summary, nil
}

func validateProfileInput(input ProfileInput) error {
	nameEmpty := strings.TrimSpace(input.Name.Ja) == "" && strings.TrimSpace(input.Name.En) == ""
	summaryEmpty := strings.TrimSpace(input.Summary.Ja) == "" && strings.TrimSpace(input.Summary.En) == ""
	if nameEmpty {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "profile name is required", nil)
	}
	if summaryEmpty {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "profile summary is required", nil)
	}
	return nil
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

func validateContactUpdateInput(input ContactUpdateInput) error {
	switch input.Status {
	case model.ContactStatusPending,
		model.ContactStatusInReview,
		model.ContactStatusResolved,
		model.ContactStatusArchived:
		// ok
	case "":
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "contact status is required", nil)
	default:
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid contact status", nil)
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

func normalizeLocalizedList(items []model.LocalizedText) []model.LocalizedText {
	if len(items) == 0 {
		return nil
	}
	result := make([]model.LocalizedText, 0, len(items))
	for _, item := range items {
		normalized := normalizeLocalized(item)
		if normalized.Ja == "" && normalized.En == "" {
			continue
		}
		result = append(result, normalized)
	}
	if len(result) == 0 {
		return nil
	}
	return result
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
