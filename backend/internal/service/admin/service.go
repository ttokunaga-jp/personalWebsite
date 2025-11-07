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

	ListTechCatalog(ctx context.Context, includeInactive bool) ([]model.TechCatalogEntry, error)

	Summary(ctx context.Context) (*model.AdminSummary, error)
}

type service struct {
	profile     repository.AdminProfileRepository
	projects    repository.AdminProjectRepository
	research    repository.AdminResearchRepository
	contacts    repository.AdminContactRepository
	blacklist   repository.BlacklistRepository
	techCatalog repository.TechCatalogRepository
}

// NewService wires repositories into the admin service.
func NewService(
	profile repository.AdminProfileRepository,
	projects repository.AdminProjectRepository,
	research repository.AdminResearchRepository,
	contacts repository.AdminContactRepository,
	blacklist repository.BlacklistRepository,
	techCatalog repository.TechCatalogRepository,
) (Service, error) {
	if profile == nil || projects == nil || research == nil || contacts == nil || blacklist == nil || techCatalog == nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "admin service: missing dependencies", nil)
	}

	return &service{
		profile:     profile,
		projects:    projects,
		research:    research,
		contacts:    contacts,
		blacklist:   blacklist,
		techCatalog: techCatalog,
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
	Tech        []ProjectTechInput
	LinkURL     string
	Year        int
	Published   bool
	SortOrder   *int
}

// ProjectTechInput represents a technology association supplied by the administrator UI.
type ProjectTechInput struct {
	MembershipID uint64
	TechID       uint64
	Context      model.TechContext
	Note         string
	SortOrder    int
}

// ResearchInput captures administrator-provided research/blog entry data.
type ResearchInput struct {
	Slug              string
	Kind              model.ResearchKind
	Title             model.LocalizedText
	Overview          model.LocalizedText
	Outcome           model.LocalizedText
	Outlook           model.LocalizedText
	ExternalURL       string
	HighlightImageURL string
	ImageAlt          model.LocalizedText
	PublishedAt       time.Time
	IsDraft           bool
	Tags              []ResearchTagInput
	Links             []ResearchLinkInput
	Assets            []ResearchAssetInput
	Tech              []ResearchTechInput
}

// ResearchTagInput represents a single tag row supplied by the administrator UI.
type ResearchTagInput struct {
	ID        uint64
	Value     string
	SortOrder int
}

// ResearchLinkInput represents an auxiliary link row supplied by the administrator UI.
type ResearchLinkInput struct {
	ID        uint64
	Type      model.ResearchLinkType
	Label     model.LocalizedText
	URL       string
	SortOrder int
}

// ResearchAssetInput represents an image asset row supplied by the administrator UI.
type ResearchAssetInput struct {
	ID        uint64
	URL       string
	Caption   model.LocalizedText
	SortOrder int
}

// ResearchTechInput represents a technology association supplied by the administrator UI.
type ResearchTechInput struct {
	MembershipID uint64
	TechID       uint64
	Context      model.TechContext
	Note         string
	SortOrder    int
}

const (
	projectEntityType      = "project"
	researchBlogEntityType = "research_blog"
)

func buildAdminResearchFromInput(id uint64, input ResearchInput) model.AdminResearch {
	entry := model.AdminResearch{
		ID:                id,
		Slug:              strings.TrimSpace(input.Slug),
		Kind:              input.Kind,
		Title:             normalizeLocalized(input.Title),
		Overview:          normalizeLocalized(input.Overview),
		Outcome:           normalizeLocalized(input.Outcome),
		Outlook:           normalizeLocalized(input.Outlook),
		ExternalURL:       strings.TrimSpace(input.ExternalURL),
		HighlightImageURL: strings.TrimSpace(input.HighlightImageURL),
		ImageAlt:          normalizeLocalized(input.ImageAlt),
		PublishedAt:       input.PublishedAt.UTC(),
		IsDraft:           input.IsDraft,
	}
	entry.Tags = normalizeResearchTags(id, input.Tags)
	entry.Links = normalizeResearchLinks(id, input.Links)
	entry.Assets = normalizeResearchAssets(id, input.Assets)
	entry.Tech = normalizeResearchTech(id, input.Tech)
	return entry
}

func normalizeResearchTags(entryID uint64, inputs []ResearchTagInput) []model.ResearchTag {
	if len(inputs) == 0 {
		return nil
	}
	result := make([]model.ResearchTag, 0, len(inputs))
	for _, item := range inputs {
		value := strings.TrimSpace(item.Value)
		if value == "" {
			continue
		}
		result = append(result, model.ResearchTag{
			ID:        item.ID,
			EntryID:   entryID,
			Value:     value,
			SortOrder: item.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func normalizeResearchLinks(entryID uint64, inputs []ResearchLinkInput) []model.ResearchLink {
	if len(inputs) == 0 {
		return nil
	}
	result := make([]model.ResearchLink, 0, len(inputs))
	for _, item := range inputs {
		url := strings.TrimSpace(item.URL)
		if url == "" {
			continue
		}
		result = append(result, model.ResearchLink{
			ID:        item.ID,
			EntryID:   entryID,
			Type:      item.Type,
			Label:     normalizeLocalized(item.Label),
			URL:       url,
			SortOrder: item.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func normalizeResearchAssets(entryID uint64, inputs []ResearchAssetInput) []model.ResearchAsset {
	if len(inputs) == 0 {
		return nil
	}
	result := make([]model.ResearchAsset, 0, len(inputs))
	for _, item := range inputs {
		url := strings.TrimSpace(item.URL)
		if url == "" {
			continue
		}
		result = append(result, model.ResearchAsset{
			ID:        item.ID,
			EntryID:   entryID,
			URL:       url,
			Caption:   normalizeLocalized(item.Caption),
			SortOrder: item.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func normalizeProjectTech(projectID int64, inputs []ProjectTechInput) []model.TechMembership {
	if len(inputs) == 0 {
		return nil
	}
	result := make([]model.TechMembership, 0, len(inputs))
	entityID := uint64(0)
	if projectID > 0 {
		entityID = uint64(projectID)
	}
	for _, item := range inputs {
		if item.TechID == 0 {
			continue
		}
		result = append(result, model.TechMembership{
			MembershipID: item.MembershipID,
			EntityType:   projectEntityType,
			EntityID:     entityID,
			Tech: model.TechCatalogEntry{
				ID: item.TechID,
			},
			Context:   item.Context,
			Note:      strings.TrimSpace(item.Note),
			SortOrder: item.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func normalizeResearchTech(entryID uint64, inputs []ResearchTechInput) []model.TechMembership {
	if len(inputs) == 0 {
		return nil
	}
	result := make([]model.TechMembership, 0, len(inputs))
	for _, item := range inputs {
		if item.TechID == 0 {
			continue
		}
		result = append(result, model.TechMembership{
			MembershipID: item.MembershipID,
			EntityType:   researchBlogEntityType,
			EntityID:     entryID,
			Tech: model.TechCatalogEntry{
				ID: item.TechID,
			},
			Context:   item.Context,
			Note:      strings.TrimSpace(item.Note),
			SortOrder: item.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func isValidResearchKind(kind model.ResearchKind) bool {
	switch kind {
	case model.ResearchKindResearch, model.ResearchKindBlog:
		return true
	default:
		return false
	}
}

func isValidResearchLinkType(linkType model.ResearchLinkType) bool {
	switch linkType {
	case model.ResearchLinkTypePaper,
		model.ResearchLinkTypeSlides,
		model.ResearchLinkTypeVideo,
		model.ResearchLinkTypeCode,
		model.ResearchLinkTypeOther:
		return true
	default:
		return false
	}
}

func isValidTechContext(context model.TechContext) bool {
	switch context {
	case model.TechContextPrimary, model.TechContextSupporting:
		return true
	default:
		return false
	}
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
		Tech:        normalizeProjectTech(0, input.Tech),
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
		Tech:        normalizeProjectTech(id, input.Tech),
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

func (s *service) ListResearch(ctx context.Context) ([]model.AdminResearch, error) {
	return s.research.ListAdminResearch(ctx)
}

func (s *service) GetResearch(ctx context.Context, id int64) (*model.AdminResearch, error) {
	if id <= 0 {
		return nil, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid research id", nil)
	}
	return s.research.GetAdminResearch(ctx, uint64(id))
}

func (s *service) CreateResearch(ctx context.Context, input ResearchInput) (*model.AdminResearch, error) {
	if err := validateResearchInput(input); err != nil {
		return nil, err
	}

	entry := buildAdminResearchFromInput(0, input)
	return s.research.CreateAdminResearch(ctx, &entry)
}

func (s *service) UpdateResearch(ctx context.Context, id int64, input ResearchInput) (*model.AdminResearch, error) {
	if err := validateResearchInput(input); err != nil {
		return nil, err
	}

	if id <= 0 {
		return nil, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid research id", nil)
	}

	entry := buildAdminResearchFromInput(uint64(id), input)
	return s.research.UpdateAdminResearch(ctx, &entry)
}

func (s *service) DeleteResearch(ctx context.Context, id int64) error {
	if id <= 0 {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid research id", nil)
	}
	return s.research.DeleteAdminResearch(ctx, uint64(id))
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

func (s *service) ListTechCatalog(ctx context.Context, includeInactive bool) ([]model.TechCatalogEntry, error) {
	entries, err := s.techCatalog.ListTechCatalog(ctx, includeInactive)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load tech catalog", err)
	}
	return entries, nil
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
		if item.IsDraft {
			draftResearch++
			continue
		}
		publishedResearch++
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
	for _, tech := range input.Tech {
		if tech.TechID == 0 {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "techId is required", nil)
		}
		if !isValidTechContext(tech.Context) {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid tech context", nil)
		}
	}
	return nil
}

func validateResearchInput(input ResearchInput) error {
	slug := strings.TrimSpace(input.Slug)
	if slug == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "research slug is required", nil)
	}
	if len(slug) > 128 {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "research slug must be 128 characters or less", nil)
	}
	if strings.Contains(slug, " ") {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "research slug cannot contain spaces", nil)
	}
	if !isValidResearchKind(input.Kind) {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid research kind", nil)
	}
	if strings.TrimSpace(input.Title.Ja) == "" && strings.TrimSpace(input.Title.En) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "research title is required", nil)
	}
	if input.PublishedAt.IsZero() {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "publishedAt is required", nil)
	}
	if strings.TrimSpace(input.ExternalURL) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "externalUrl is required", nil)
	}
	for _, tag := range input.Tags {
		if strings.TrimSpace(tag.Value) == "" {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "tag value is required", nil)
		}
	}
	for _, link := range input.Links {
		if !isValidResearchLinkType(link.Type) {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid research link type", nil)
		}
		if strings.TrimSpace(link.URL) == "" {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "link url is required", nil)
		}
	}
	for _, asset := range input.Assets {
		if strings.TrimSpace(asset.URL) == "" {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "asset url is required", nil)
		}
	}
	for _, tech := range input.Tech {
		if tech.TechID == 0 {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "techId is required", nil)
		}
		if !isValidTechContext(tech.Context) {
			return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid tech context", nil)
		}
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

func copyIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}
