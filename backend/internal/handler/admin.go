package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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

	input := req.toInput()
	project, err := h.svc.CreateProject(c.Request.Context(), input)
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

func (h *AdminHandler) ListBlogPosts(c *gin.Context) {
	posts, err := h.svc.ListBlogPosts(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *AdminHandler) CreateBlogPost(c *gin.Context) {
	var req blogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid blog post payload", err))
		return
	}
	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}
	post, svcErr := h.svc.CreateBlogPost(c.Request.Context(), input)
	if svcErr != nil {
		respondError(c, svcErr)
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *AdminHandler) GetBlogPost(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	post, err := h.svc.GetBlogPost(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *AdminHandler) UpdateBlogPost(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req blogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid blog post payload", err))
		return
	}
	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}
	post, svcErr := h.svc.UpdateBlogPost(c.Request.Context(), id, input)
	if svcErr != nil {
		respondError(c, svcErr)
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *AdminHandler) DeleteBlogPost(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteBlogPost(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AdminHandler) ListMeetings(c *gin.Context) {
	meetings, err := h.svc.ListMeetings(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, meetings)
}

func (h *AdminHandler) CreateMeeting(c *gin.Context) {
	var req meetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid meeting payload", err))
		return
	}
	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}
	meeting, svcErr := h.svc.CreateMeeting(c.Request.Context(), input)
	if svcErr != nil {
		respondError(c, svcErr)
		return
	}
	c.JSON(http.StatusCreated, meeting)
}

func (h *AdminHandler) GetMeeting(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	meeting, err := h.svc.GetMeeting(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, meeting)
}

func (h *AdminHandler) UpdateMeeting(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req meetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid meeting payload", err))
		return
	}
	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}
	meeting, svcErr := h.svc.UpdateMeeting(c.Request.Context(), id, input)
	if svcErr != nil {
		respondError(c, svcErr)
		return
	}
	c.JSON(http.StatusOK, meeting)
}

func (h *AdminHandler) DeleteMeeting(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteMeeting(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

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

func parseIDParam(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid id", err))
		return 0, false
	}
	return id, true
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

type blogPostRequest struct {
	Title       model.LocalizedText `json:"title"`
	Summary     model.LocalizedText `json:"summary"`
	ContentMD   model.LocalizedText `json:"contentMd"`
	Tags        []string            `json:"tags"`
	Published   bool                `json:"published"`
	PublishedAt *string             `json:"publishedAt"`
}

func (r blogPostRequest) toInput() (adminsvc.BlogPostInput, error) {
	var publishedAt *time.Time
	if r.PublishedAt != nil && strings.TrimSpace(*r.PublishedAt) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*r.PublishedAt))
		if err != nil {
			return adminsvc.BlogPostInput{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid publishedAt format", err)
		}
		publishedAt = &parsed
	}
	return adminsvc.BlogPostInput{
		Title:       r.Title,
		Summary:     r.Summary,
		ContentMD:   r.ContentMD,
		Tags:        r.Tags,
		Published:   r.Published,
		PublishedAt: publishedAt,
	}, nil
}

type meetingRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Datetime        string `json:"datetime"`
	DurationMinutes int    `json:"durationMinutes"`
	MeetURL         string `json:"meetUrl"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
}

func (r meetingRequest) toInput() (adminsvc.MeetingInput, error) {
	parsedTime, err := time.Parse(time.RFC3339, strings.TrimSpace(r.Datetime))
	if err != nil {
		return adminsvc.MeetingInput{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid datetime", err)
	}
	status := model.MeetingStatus(strings.TrimSpace(r.Status))
	return adminsvc.MeetingInput{
		Name:            r.Name,
		Email:           r.Email,
		Datetime:        parsedTime,
		DurationMinutes: r.DurationMinutes,
		MeetURL:         r.MeetURL,
		Status:          status,
		Notes:           r.Notes,
	}, nil
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
