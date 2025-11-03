package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	adminsvc "github.com/takumi/personal-website/internal/service/admin"
)

// AdminHandler exposes administrative API endpoints.
type AdminHandler struct {
	svc adminsvc.Service
}

// NewAdminHandler constructs the admin handler.
func NewAdminHandler(svc adminsvc.Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// Summary returns aggregated dashboard metrics.
func (h *AdminHandler) Summary(c *gin.Context) {
	summary, err := h.svc.Summary(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}

// Profile management --------------------------------------------------------

// GetProfile returns the current editable profile.
func (h *AdminHandler) GetProfile(c *gin.Context) {
	profile, err := h.svc.GetProfile(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateProfile persists profile updates.
func (h *AdminHandler) UpdateProfile(c *gin.Context) {
	var req profileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid profile payload", err))
		return
	}

	profile, err := h.svc.UpdateProfile(c.Request.Context(), req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// Project management -------------------------------------------------------

func (h *AdminHandler) ListProjects(c *gin.Context) {
	projects, err := h.svc.ListProjects(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *AdminHandler) CreateProject(c *gin.Context) {
	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid project payload", err))
		return
	}

	project, err := h.svc.CreateProject(c.Request.Context(), req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, project)
}

func (h *AdminHandler) GetProject(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	project, err := h.svc.GetProject(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *AdminHandler) UpdateProject(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid project payload", err))
		return
	}
	project, err := h.svc.UpdateProject(c.Request.Context(), id, req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *AdminHandler) DeleteProject(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteProject(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Research management ------------------------------------------------------

func (h *AdminHandler) ListResearch(c *gin.Context) {
	research, err := h.svc.ListResearch(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, research)
}

func (h *AdminHandler) CreateResearch(c *gin.Context) {
	var req researchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid research payload", err))
		return
	}
	item, err := h.svc.CreateResearch(c.Request.Context(), req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *AdminHandler) GetResearch(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	item, err := h.svc.GetResearch(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *AdminHandler) UpdateResearch(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req researchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid research payload", err))
		return
	}
	item, err := h.svc.UpdateResearch(c.Request.Context(), id, req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *AdminHandler) DeleteResearch(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteResearch(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Contact management -------------------------------------------------------

func (h *AdminHandler) ListContacts(c *gin.Context) {
	messages, err := h.svc.ListContactMessages(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *AdminHandler) GetContact(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid id", nil))
		return
	}

	message, err := h.svc.GetContactMessage(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, message)
}

func (h *AdminHandler) UpdateContact(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid id", nil))
		return
	}

	var req contactUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid contact payload", err))
		return
	}

	message, err := h.svc.UpdateContactMessage(c.Request.Context(), id, req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, message)
}

func (h *AdminHandler) DeleteContact(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid id", nil))
		return
	}
	if err := h.svc.DeleteContactMessage(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Blacklist management ------------------------------------------------------

func (h *AdminHandler) ListBlacklist(c *gin.Context) {
	entries, err := h.svc.ListBlacklist(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *AdminHandler) CreateBlacklist(c *gin.Context) {
	var req blacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid blacklist payload", err))
		return
	}
	entry, err := h.svc.AddBlacklistEntry(c.Request.Context(), req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *AdminHandler) UpdateBlacklist(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req blacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid blacklist payload", err))
		return
	}
	entry, err := h.svc.UpdateBlacklistEntry(c.Request.Context(), id, req.toInput())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *AdminHandler) DeleteBlacklist(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.RemoveBlacklistEntry(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Helpers ------------------------------------------------------------------

func parseIDParam(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid id", err))
		return 0, false
	}
	return id, true
}

// Request payloads ---------------------------------------------------------

type profileRequest struct {
	Name        model.LocalizedText   `json:"name"`
	Title       model.LocalizedText   `json:"title"`
	Affiliation model.LocalizedText   `json:"affiliation"`
	Lab         model.LocalizedText   `json:"lab"`
	Summary     model.LocalizedText   `json:"summary"`
	Skills      []model.LocalizedText `json:"skills"`
	FocusAreas  []model.LocalizedText `json:"focusAreas"`
}

func (r profileRequest) toInput() adminsvc.ProfileInput {
	return adminsvc.ProfileInput{
		Name:        r.Name,
		Title:       r.Title,
		Affiliation: r.Affiliation,
		Lab:         r.Lab,
		Summary:     r.Summary,
		Skills:      r.Skills,
		FocusAreas:  r.FocusAreas,
	}
}

type projectRequest struct {
	Title       model.LocalizedText `json:"title"`
	Description model.LocalizedText `json:"description"`
	TechStack   []string            `json:"techStack"`
	LinkURL     string              `json:"linkUrl"`
	Year        int                 `json:"year"`
	Published   bool                `json:"published"`
	SortOrder   *int                `json:"sortOrder"`
}

func (r projectRequest) toInput() adminsvc.ProjectInput {
	return adminsvc.ProjectInput{
		Title:       r.Title,
		Description: r.Description,
		TechStack:   r.TechStack,
		LinkURL:     r.LinkURL,
		Year:        r.Year,
		Published:   r.Published,
		SortOrder:   r.SortOrder,
	}
}

type researchRequest struct {
	Title     model.LocalizedText `json:"title"`
	Summary   model.LocalizedText `json:"summary"`
	ContentMD model.LocalizedText `json:"contentMd"`
	Year      int                 `json:"year"`
	Published bool                `json:"published"`
}

func (r researchRequest) toInput() adminsvc.ResearchInput {
	return adminsvc.ResearchInput{
		Title:     r.Title,
		Summary:   r.Summary,
		ContentMD: r.ContentMD,
		Year:      r.Year,
		Published: r.Published,
	}
}

type contactUpdateRequest struct {
	Topic     string `json:"topic"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	AdminNote string `json:"adminNote"`
}

func (r contactUpdateRequest) toInput() adminsvc.ContactUpdateInput {
	return adminsvc.ContactUpdateInput{
		Topic:     r.Topic,
		Message:   r.Message,
		Status:    model.ContactStatus(strings.TrimSpace(r.Status)),
		AdminNote: r.AdminNote,
	}
}

type blacklistRequest struct {
	Email  string `json:"email"`
	Reason string `json:"reason"`
}

func (r blacklistRequest) toInput() adminsvc.BlacklistInput {
	return adminsvc.BlacklistInput{
		Email:  r.Email,
		Reason: r.Reason,
	}
}
