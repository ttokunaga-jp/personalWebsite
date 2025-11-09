package handler

import (
	"fmt"
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

// Profile management --------------------------------------------------------

// Tech catalog --------------------------------------------------------------

func (h *AdminHandler) ListTechCatalog(c *gin.Context) {
	includeInactive := c.Query("includeInactive") == "true"
	entries, err := h.svc.ListTechCatalog(c.Request.Context(), includeInactive)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entries})
}

func (h *AdminHandler) CreateTechCatalogEntry(c *gin.Context) {
	var req techCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid tech catalog payload", err))
		return
	}

	input, appErr := req.toInput()
	if appErr != nil {
		respondError(c, appErr)
		return
	}

	entry, err := h.svc.CreateTechCatalogEntry(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": entry})
}

func (h *AdminHandler) UpdateTechCatalogEntry(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req techCatalogUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid tech catalog update payload", err))
		return
	}

	input, appErr := req.toInput()
	if appErr != nil {
		respondError(c, appErr)
		return
	}

	entry, err := h.svc.UpdateTechCatalogEntry(c.Request.Context(), uint64(id), input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": entry})
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

	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}

	profile, err := h.svc.UpdateProfile(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// Home settings ------------------------------------------------------------

func (h *AdminHandler) GetHomeSettings(c *gin.Context) {
	settings, err := h.svc.GetHomeSettings(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *AdminHandler) UpdateHomeSettings(c *gin.Context) {
	var req homeSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid home settings payload", err))
		return
	}

	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}

	settings, err := h.svc.UpdateHomeSettings(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
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
	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}
	item, err := h.svc.CreateResearch(c.Request.Context(), input)
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
	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}
	item, err := h.svc.UpdateResearch(c.Request.Context(), id, input)
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

func (h *AdminHandler) GetContactSettings(c *gin.Context) {
	settings, err := h.svc.GetContactSettings(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *AdminHandler) UpdateContactSettings(c *gin.Context) {
	var req contactSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid contact settings payload", err))
		return
	}
	input, err := req.toInput()
	if err != nil {
		respondError(c, err)
		return
	}

	settings, err := h.svc.UpdateContactSettings(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

// Social links & meeting template -------------------------------------------

func (h *AdminHandler) ListSocialLinks(c *gin.Context) {
	links, err := h.svc.ListSocialLinks(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": links})
}

func (h *AdminHandler) ReplaceSocialLinks(c *gin.Context) {
	var req socialLinksReplaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid social links payload", err))
		return
	}

	inputs, appErr := req.toInputs()
	if appErr != nil {
		respondError(c, appErr)
		return
	}

	links, err := h.svc.ReplaceSocialLinks(c.Request.Context(), inputs)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": links})
}

func (h *AdminHandler) GetMeetingURLTemplate(c *gin.Context) {
	template, err := h.svc.GetMeetingURLTemplate(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (h *AdminHandler) UpdateMeetingURLTemplate(c *gin.Context) {
	var req meetingTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid meeting template payload", err))
		return
	}

	template := strings.TrimSpace(req.Template)
	updated, err := h.svc.UpdateMeetingURLTemplate(c.Request.Context(), template)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

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
	c.JSON(http.StatusOK, gin.H{"data": entries})
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
	c.JSON(http.StatusCreated, gin.H{"data": entry})
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
	c.JSON(http.StatusOK, gin.H{"data": entry})
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

// Reservations ----------------------------------------------------------------

type reservationResponse struct {
	ID                     uint64                         `json:"id"`
	Name                   string                         `json:"name"`
	Email                  string                         `json:"email"`
	Topic                  string                         `json:"topic"`
	Message                string                         `json:"message"`
	StartAt                time.Time                      `json:"startAt"`
	EndAt                  time.Time                      `json:"endAt"`
	DurationMinutes        int                            `json:"durationMinutes"`
	GoogleEventID          string                         `json:"googleEventId,omitempty"`
	GoogleCalendarStatus   string                         `json:"googleCalendarStatus,omitempty"`
	Status                 model.MeetingReservationStatus `json:"status"`
	ConfirmationSentAt     *time.Time                     `json:"confirmationSentAt,omitempty"`
	LastNotificationSentAt *time.Time                     `json:"lastNotificationSentAt,omitempty"`
	LookupHash             string                         `json:"lookupHash"`
	CancellationReason     string                         `json:"cancellationReason,omitempty"`
	CreatedAt              time.Time                      `json:"createdAt"`
	UpdatedAt              time.Time                      `json:"updatedAt"`
	Notifications          []model.MeetingNotification    `json:"notifications,omitempty"`
}

type reservationUpdateRequest struct {
	Status             string `json:"status"`
	CancellationReason string `json:"cancellationReason"`
}

func (h *AdminHandler) ListReservations(c *gin.Context) {
	filter, err := parseReservationFilter(c)
	if err != nil {
		respondError(c, err)
		return
	}

	reservations, err := h.svc.ListReservations(c.Request.Context(), filter)
	if err != nil {
		respondError(c, err)
		return
	}

	response := make([]reservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		notifications, err := h.svc.ListReservationNotifications(c.Request.Context(), reservation.ID)
		if err != nil {
			respondError(c, err)
			return
		}
		response = append(response, makeReservationResponse(reservation, notifications))
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *AdminHandler) UpdateReservationStatus(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req reservationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid reservation update payload", err))
		return
	}

	status := model.MeetingReservationStatus(strings.TrimSpace(strings.ToLower(req.Status)))
	reservation, err := h.svc.UpdateReservationStatus(c.Request.Context(), uint64(id), status, req.CancellationReason)
	if err != nil {
		respondError(c, err)
		return
	}

	notifications, err := h.svc.ListReservationNotifications(c.Request.Context(), reservation.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": makeReservationResponse(*reservation, notifications)})
}

func (h *AdminHandler) RetryReservationNotification(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	reservation, err := h.svc.RetryReservationNotification(c.Request.Context(), uint64(id))
	if err != nil {
		respondError(c, err)
		return
	}

	notifications, err := h.svc.ListReservationNotifications(c.Request.Context(), reservation.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": makeReservationResponse(*reservation, notifications)})
}

func parseReservationFilter(c *gin.Context) (adminsvc.ReservationFilter, error) {
	filter := adminsvc.ReservationFilter{}

	statusValues := c.QueryArray("status")
	if value := strings.TrimSpace(c.Query("status")); value != "" {
		statusValues = append(statusValues, value)
	}
	if len(statusValues) > 0 {
		for _, raw := range statusValues {
			for _, piece := range strings.Split(raw, ",") {
				value := strings.TrimSpace(strings.ToLower(piece))
				if value == "" {
					continue
				}
				status := model.MeetingReservationStatus(value)
				switch status {
				case model.MeetingReservationStatusPending,
					model.MeetingReservationStatusConfirmed,
					model.MeetingReservationStatusCancelled:
					filter.Status = append(filter.Status, status)
				default:
					return adminsvc.ReservationFilter{}, errs.New(
						errs.CodeInvalidInput,
						http.StatusBadRequest,
						fmt.Sprintf("unsupported reservation status %q", value),
						nil,
					)
				}
			}
		}
	}

	filter.Email = strings.TrimSpace(c.Query("email"))

	if dateString := strings.TrimSpace(c.Query("date")); dateString != "" {
		parsed, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			return adminsvc.ReservationFilter{}, errs.New(
				errs.CodeInvalidInput,
				http.StatusBadRequest,
				"invalid date format (expected YYYY-MM-DD)",
				err,
			)
		}
		filter.Date = &parsed
	}

	return filter, nil
}

func makeReservationResponse(reservation model.MeetingReservation, notifications []model.MeetingNotification) reservationResponse {
	response := reservationResponse{
		ID:                     reservation.ID,
		Name:                   strings.TrimSpace(reservation.Name),
		Email:                  strings.TrimSpace(reservation.Email),
		Topic:                  strings.TrimSpace(reservation.Topic),
		Message:                reservation.Message,
		StartAt:                reservation.StartAt,
		EndAt:                  reservation.EndAt,
		DurationMinutes:        reservation.DurationMinutes,
		GoogleEventID:          strings.TrimSpace(reservation.GoogleEventID),
		GoogleCalendarStatus:   strings.TrimSpace(reservation.GoogleCalendarStatus),
		Status:                 reservation.Status,
		ConfirmationSentAt:     reservation.ConfirmationSentAt,
		LastNotificationSentAt: reservation.LastNotificationSentAt,
		LookupHash:             reservation.LookupHash,
		CancellationReason:     strings.TrimSpace(reservation.CancellationReason),
		CreatedAt:              reservation.CreatedAt,
		UpdatedAt:              reservation.UpdatedAt,
	}
	if len(notifications) > 0 {
		response.Notifications = notifications
	}
	return response
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
	DisplayName  string               `json:"displayName"`
	Headline     model.LocalizedText  `json:"headline"`
	Summary      model.LocalizedText  `json:"summary"`
	AvatarURL    string               `json:"avatarUrl"`
	Location     model.LocalizedText  `json:"location"`
	Theme        profileThemeRequest  `json:"theme"`
	Lab          profileLabRequest    `json:"lab"`
	Affiliations []affiliationRequest `json:"affiliations"`
	Communities  []affiliationRequest `json:"communities"`
	WorkHistory  []workHistoryRequest `json:"workHistory"`
	SocialLinks  []socialLinkRequest  `json:"socialLinks"`
}

type profileThemeRequest struct {
	Mode        string `json:"mode"`
	AccentColor string `json:"accentColor"`
}

type profileLabRequest struct {
	Name    model.LocalizedText `json:"name"`
	Advisor model.LocalizedText `json:"advisor"`
	Room    model.LocalizedText `json:"room"`
	URL     string              `json:"url"`
}

type affiliationRequest struct {
	ID          uint64              `json:"id"`
	Name        string              `json:"name"`
	URL         string              `json:"url"`
	Description model.LocalizedText `json:"description"`
	StartedAt   string              `json:"startedAt"`
	SortOrder   int                 `json:"sortOrder"`
}

type workHistoryRequest struct {
	ID           uint64              `json:"id"`
	Organization model.LocalizedText `json:"organization"`
	Role         model.LocalizedText `json:"role"`
	Summary      model.LocalizedText `json:"summary"`
	StartedAt    string              `json:"startedAt"`
	EndedAt      *string             `json:"endedAt"`
	ExternalURL  string              `json:"externalUrl"`
	SortOrder    int                 `json:"sortOrder"`
}

type socialLinkRequest struct {
	ID        uint64              `json:"id"`
	Provider  string              `json:"provider"`
	Label     model.LocalizedText `json:"label"`
	URL       string              `json:"url"`
	IsFooter  bool                `json:"isFooter"`
	SortOrder int                 `json:"sortOrder"`
}

func (r profileRequest) toInput() (adminsvc.ProfileInput, error) {
	input := adminsvc.ProfileInput{
		DisplayName: strings.TrimSpace(r.DisplayName),
		Headline:    r.Headline,
		Summary:     r.Summary,
		AvatarURL:   strings.TrimSpace(r.AvatarURL),
		Location:    r.Location,
		Theme: adminsvc.ProfileThemeInput{
			Mode:        strings.TrimSpace(r.Theme.Mode),
			AccentColor: strings.TrimSpace(r.Theme.AccentColor),
		},
		Lab: adminsvc.ProfileLabInput{
			Name:    r.Lab.Name,
			Advisor: r.Lab.Advisor,
			Room:    r.Lab.Room,
			URL:     strings.TrimSpace(r.Lab.URL),
		},
	}

	for idx, affiliation := range r.Affiliations {
		startedAt, err := parseISOTime(affiliation.StartedAt, fmt.Sprintf("affiliations[%d].startedAt", idx))
		if err != nil {
			return adminsvc.ProfileInput{}, err
		}
		input.Affiliations = append(input.Affiliations, adminsvc.ProfileAffiliationInput{
			ID:          affiliation.ID,
			Name:        affiliation.Name,
			URL:         affiliation.URL,
			Description: affiliation.Description,
			StartedAt:   startedAt,
			SortOrder:   affiliation.SortOrder,
		})
	}

	for idx, community := range r.Communities {
		startedAt, err := parseISOTime(community.StartedAt, fmt.Sprintf("communities[%d].startedAt", idx))
		if err != nil {
			return adminsvc.ProfileInput{}, err
		}
		input.Communities = append(input.Communities, adminsvc.ProfileAffiliationInput{
			ID:          community.ID,
			Name:        community.Name,
			URL:         community.URL,
			Description: community.Description,
			StartedAt:   startedAt,
			SortOrder:   community.SortOrder,
		})
	}

	for idx, history := range r.WorkHistory {
		startedAt, err := parseISOTime(history.StartedAt, fmt.Sprintf("workHistory[%d].startedAt", idx))
		if err != nil {
			return adminsvc.ProfileInput{}, err
		}
		var endedAt *time.Time
		if history.EndedAt != nil && strings.TrimSpace(*history.EndedAt) != "" {
			parsed, err := parseISOTime(*history.EndedAt, fmt.Sprintf("workHistory[%d].endedAt", idx))
			if err != nil {
				return adminsvc.ProfileInput{}, err
			}
			endedAt = &parsed
		}
		input.WorkHistory = append(input.WorkHistory, adminsvc.ProfileWorkHistoryInput{
			ID:           history.ID,
			Organization: history.Organization,
			Role:         history.Role,
			Summary:      history.Summary,
			StartedAt:    startedAt,
			EndedAt:      endedAt,
			ExternalURL:  history.ExternalURL,
			SortOrder:    history.SortOrder,
		})
	}

	for idx, link := range r.SocialLinks {
		provider := model.ProfileSocialProvider(strings.ToLower(strings.TrimSpace(link.Provider)))
		switch provider {
		case model.ProfileSocialProviderGitHub,
			model.ProfileSocialProviderZenn,
			model.ProfileSocialProviderLinkedIn,
			model.ProfileSocialProviderX,
			model.ProfileSocialProviderEmail,
			model.ProfileSocialProviderOther:
			// ok
		default:
			return adminsvc.ProfileInput{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, fmt.Sprintf("socialLinks[%d].provider is invalid", idx), nil)
		}
		input.SocialLinks = append(input.SocialLinks, adminsvc.ProfileSocialLinkInput{
			ID:        link.ID,
			Provider:  provider,
			Label:     link.Label,
			URL:       link.URL,
			IsFooter:  link.IsFooter,
			SortOrder: link.SortOrder,
		})
	}

	return input, nil
}

type projectRequest struct {
	Title       model.LocalizedText  `json:"title"`
	Description model.LocalizedText  `json:"description"`
	Tech        []projectTechRequest `json:"tech"`
	LinkURL     string               `json:"linkUrl"`
	Year        int                  `json:"year"`
	Published   bool                 `json:"published"`
	SortOrder   *int                 `json:"sortOrder"`
}

func (r projectRequest) toInput() adminsvc.ProjectInput {
	input := adminsvc.ProjectInput{
		Title:       r.Title,
		Description: r.Description,
		LinkURL:     r.LinkURL,
		Year:        r.Year,
		Published:   r.Published,
		SortOrder:   r.SortOrder,
	}
	if len(r.Tech) > 0 {
		tech := make([]adminsvc.ProjectTechInput, 0, len(r.Tech))
		for _, membership := range r.Tech {
			tech = append(tech, adminsvc.ProjectTechInput{
				MembershipID: membership.MembershipID,
				TechID:       membership.TechID,
				Context:      model.TechContext(strings.TrimSpace(membership.Context)),
				Note:         strings.TrimSpace(membership.Note),
				SortOrder:    membership.SortOrder,
			})
		}
		input.Tech = tech
	}
	return input
}

type projectTechRequest struct {
	MembershipID uint64 `json:"membershipId"`
	TechID       uint64 `json:"techId"`
	Context      string `json:"context"`
	Note         string `json:"note"`
	SortOrder    int    `json:"sortOrder"`
}

type researchRequest struct {
	Slug              string                 `json:"slug"`
	Kind              string                 `json:"kind"`
	Title             model.LocalizedText    `json:"title"`
	Overview          model.LocalizedText    `json:"overview"`
	Outcome           model.LocalizedText    `json:"outcome"`
	Outlook           model.LocalizedText    `json:"outlook"`
	ExternalURL       string                 `json:"externalUrl"`
	HighlightImageURL string                 `json:"highlightImageUrl"`
	ImageAlt          model.LocalizedText    `json:"imageAlt"`
	PublishedAt       string                 `json:"publishedAt"`
	IsDraft           bool                   `json:"isDraft"`
	Tags              []researchTagRequest   `json:"tags"`
	Links             []researchLinkRequest  `json:"links"`
	Assets            []researchAssetRequest `json:"assets"`
	Tech              []researchTechRequest  `json:"tech"`
}

type researchTagRequest struct {
	ID        uint64 `json:"id"`
	Value     string `json:"value"`
	SortOrder int    `json:"sortOrder"`
}

type researchLinkRequest struct {
	ID        uint64              `json:"id"`
	Type      string              `json:"type"`
	Label     model.LocalizedText `json:"label"`
	URL       string              `json:"url"`
	SortOrder int                 `json:"sortOrder"`
}

type researchAssetRequest struct {
	ID        uint64              `json:"id"`
	URL       string              `json:"url"`
	Caption   model.LocalizedText `json:"caption"`
	SortOrder int                 `json:"sortOrder"`
}

type researchTechRequest struct {
	MembershipID uint64 `json:"membershipId"`
	TechID       uint64 `json:"techId"`
	Context      string `json:"context"`
	Note         string `json:"note"`
	SortOrder    int    `json:"sortOrder"`
}

func (r researchRequest) toInput() (adminsvc.ResearchInput, error) {
	publishedAt, err := parseRFC3339Timestamp(r.PublishedAt)
	if err != nil {
		return adminsvc.ResearchInput{}, err
	}

	input := adminsvc.ResearchInput{
		Slug:              strings.TrimSpace(r.Slug),
		Kind:              model.ResearchKind(strings.TrimSpace(r.Kind)),
		Title:             r.Title,
		Overview:          r.Overview,
		Outcome:           r.Outcome,
		Outlook:           r.Outlook,
		ExternalURL:       strings.TrimSpace(r.ExternalURL),
		HighlightImageURL: strings.TrimSpace(r.HighlightImageURL),
		ImageAlt:          r.ImageAlt,
		PublishedAt:       publishedAt,
		IsDraft:           r.IsDraft,
	}

	if len(r.Tags) > 0 {
		tags := make([]adminsvc.ResearchTagInput, 0, len(r.Tags))
		for _, tag := range r.Tags {
			tags = append(tags, adminsvc.ResearchTagInput{
				ID:        tag.ID,
				Value:     strings.TrimSpace(tag.Value),
				SortOrder: tag.SortOrder,
			})
		}
		input.Tags = tags
	}

	if len(r.Links) > 0 {
		links := make([]adminsvc.ResearchLinkInput, 0, len(r.Links))
		for _, link := range r.Links {
			links = append(links, adminsvc.ResearchLinkInput{
				ID:        link.ID,
				Type:      model.ResearchLinkType(strings.TrimSpace(link.Type)),
				Label:     link.Label,
				URL:       strings.TrimSpace(link.URL),
				SortOrder: link.SortOrder,
			})
		}
		input.Links = links
	}

	if len(r.Assets) > 0 {
		assets := make([]adminsvc.ResearchAssetInput, 0, len(r.Assets))
		for _, asset := range r.Assets {
			assets = append(assets, adminsvc.ResearchAssetInput{
				ID:        asset.ID,
				URL:       strings.TrimSpace(asset.URL),
				Caption:   asset.Caption,
				SortOrder: asset.SortOrder,
			})
		}
		input.Assets = assets
	}

	if len(r.Tech) > 0 {
		tech := make([]adminsvc.ResearchTechInput, 0, len(r.Tech))
		for _, item := range r.Tech {
			tech = append(tech, adminsvc.ResearchTechInput{
				MembershipID: item.MembershipID,
				TechID:       item.TechID,
				Context:      model.TechContext(strings.TrimSpace(item.Context)),
				Note:         strings.TrimSpace(item.Note),
				SortOrder:    item.SortOrder,
			})
		}
		input.Tech = tech
	}

	return input, nil
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

type techCatalogRequest struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"displayName"`
	Category    string `json:"category"`
	Level       string `json:"level"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sortOrder"`
	Active      *bool  `json:"active"`
}

func (r techCatalogRequest) toInput() (adminsvc.TechCatalogInput, *errs.AppError) {
	active := true
	if r.Active != nil {
		active = *r.Active
	}
	return adminsvc.TechCatalogInput{
		Slug:        r.Slug,
		DisplayName: r.DisplayName,
		Category:    r.Category,
		Level:       model.TechLevel(strings.TrimSpace(strings.ToLower(r.Level))),
		Icon:        r.Icon,
		SortOrder:   r.SortOrder,
		Active:      active,
	}, nil
}

type techCatalogUpdateRequest struct {
	Slug        *string `json:"slug"`
	DisplayName *string `json:"displayName"`
	Category    *string `json:"category"`
	Level       *string `json:"level"`
	Icon        *string `json:"icon"`
	SortOrder   *int    `json:"sortOrder"`
	Active      *bool   `json:"active"`
}

func (r techCatalogUpdateRequest) toInput() (adminsvc.TechCatalogUpdateInput, *errs.AppError) {
	input := adminsvc.TechCatalogUpdateInput{
		Slug:        r.Slug,
		DisplayName: r.DisplayName,
		Category:    r.Category,
		Icon:        r.Icon,
		SortOrder:   r.SortOrder,
		Active:      r.Active,
	}
	if r.Level != nil {
		level := model.TechLevel(strings.TrimSpace(strings.ToLower(*r.Level)))
		input.Level = &level
	}
	return input, nil
}

type socialLinksReplaceRequest struct {
	Links []replaceSocialLinkRequest `json:"links"`
}

type replaceSocialLinkRequest struct {
	Provider  string              `json:"provider"`
	Label     model.LocalizedText `json:"label"`
	URL       string              `json:"url"`
	IsFooter  bool                `json:"isFooter"`
	SortOrder int                 `json:"sortOrder"`
}

func (r socialLinksReplaceRequest) toInputs() ([]adminsvc.SocialLinkInput, *errs.AppError) {
	inputs := make([]adminsvc.SocialLinkInput, 0, len(r.Links))
	for _, link := range r.Links {
		provider := model.ProfileSocialProvider(strings.TrimSpace(strings.ToLower(link.Provider)))
		inputs = append(inputs, adminsvc.SocialLinkInput{
			Provider:  provider,
			Label:     link.Label,
			URL:       link.URL,
			IsFooter:  link.IsFooter,
			SortOrder: link.SortOrder,
		})
	}
	return inputs, nil
}

type meetingTemplateRequest struct {
	Template string `json:"template"`
}

type contactSettingsRequest struct {
	ID                 uint64                `json:"id"`
	HeroTitle          model.LocalizedText   `json:"heroTitle"`
	HeroDescription    model.LocalizedText   `json:"heroDescription"`
	Topics             []contactTopicRequest `json:"topics"`
	ConsentText        model.LocalizedText   `json:"consentText"`
	MinimumLeadHours   int                   `json:"minimumLeadHours"`
	RecaptchaSiteKey   string                `json:"recaptchaSiteKey"`
	SupportEmail       string                `json:"supportEmail"`
	CalendarTimezone   string                `json:"calendarTimezone"`
	GoogleCalendarID   string                `json:"googleCalendarId"`
	BookingWindowDays  int                   `json:"bookingWindowDays"`
	MeetingURLTemplate string                `json:"meetingUrlTemplate"`
	UpdatedAt          string                `json:"updatedAt"`
}

type contactTopicRequest struct {
	ID          string              `json:"id"`
	Label       model.LocalizedText `json:"label"`
	Description model.LocalizedText `json:"description"`
}

type homeSettingsRequest struct {
	ID           uint64                  `json:"id"`
	ProfileID    uint64                  `json:"profileId"`
	HeroSubtitle model.LocalizedText     `json:"heroSubtitle"`
	QuickLinks   []homeQuickLinkRequest  `json:"quickLinks"`
	ChipSources  []homeChipSourceRequest `json:"chipSources"`
	UpdatedAt    string                  `json:"updatedAt"`
}

type homeQuickLinkRequest struct {
	ID          uint64              `json:"id"`
	Section     string              `json:"section"`
	Label       model.LocalizedText `json:"label"`
	Description model.LocalizedText `json:"description"`
	CTA         model.LocalizedText `json:"cta"`
	TargetURL   string              `json:"targetUrl"`
	SortOrder   int                 `json:"sortOrder"`
}

type homeChipSourceRequest struct {
	ID        uint64              `json:"id"`
	Source    string              `json:"source"`
	Label     model.LocalizedText `json:"label"`
	Limit     int                 `json:"limit"`
	SortOrder int                 `json:"sortOrder"`
}

func (r contactSettingsRequest) toInput() (adminsvc.ContactSettingsInput, error) {
	updatedAt := strings.TrimSpace(r.UpdatedAt)
	if updatedAt == "" {
		return adminsvc.ContactSettingsInput{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "updatedAt is required", nil)
	}
	parsed, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return adminsvc.ContactSettingsInput{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "updatedAt must be RFC3339 timestamp", err)
	}

	topics := make([]adminsvc.ContactTopicInput, 0, len(r.Topics))
	for _, topic := range r.Topics {
		topics = append(topics, adminsvc.ContactTopicInput{
			ID:          topic.ID,
			Label:       topic.Label,
			Description: topic.Description,
		})
	}

	return adminsvc.ContactSettingsInput{
		ID:                 r.ID,
		HeroTitle:          r.HeroTitle,
		HeroDescription:    r.HeroDescription,
		Topics:             topics,
		ConsentText:        r.ConsentText,
		MinimumLeadHours:   r.MinimumLeadHours,
		RecaptchaSiteKey:   r.RecaptchaSiteKey,
		SupportEmail:       r.SupportEmail,
		CalendarTimezone:   r.CalendarTimezone,
		GoogleCalendarID:   r.GoogleCalendarID,
		BookingWindowDays:  r.BookingWindowDays,
		MeetingURLTemplate: strings.TrimSpace(r.MeetingURLTemplate),
		ExpectedUpdatedAt:  parsed,
	}, nil
}

func (r homeSettingsRequest) toInput() (adminsvc.HomeSettingsInput, error) {
	updatedAt := strings.TrimSpace(r.UpdatedAt)
	if updatedAt == "" {
		return adminsvc.HomeSettingsInput{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "updatedAt is required", nil)
	}
	parsed, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return adminsvc.HomeSettingsInput{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "updatedAt must be RFC3339 timestamp", err)
	}

	quickLinks := make([]adminsvc.HomeQuickLinkInput, 0, len(r.QuickLinks))
	for _, link := range r.QuickLinks {
		quickLinks = append(quickLinks, adminsvc.HomeQuickLinkInput{
			ID:          link.ID,
			Section:     link.Section,
			Label:       link.Label,
			Description: link.Description,
			CTA:         link.CTA,
			TargetURL:   strings.TrimSpace(link.TargetURL),
			SortOrder:   link.SortOrder,
		})
	}

	chipSources := make([]adminsvc.HomeChipSourceInput, 0, len(r.ChipSources))
	for _, chip := range r.ChipSources {
		chipSources = append(chipSources, adminsvc.HomeChipSourceInput{
			ID:        chip.ID,
			Source:    chip.Source,
			Label:     chip.Label,
			Limit:     chip.Limit,
			SortOrder: chip.SortOrder,
		})
	}

	return adminsvc.HomeSettingsInput{
		ID:                r.ID,
		ProfileID:         r.ProfileID,
		HeroSubtitle:      r.HeroSubtitle,
		QuickLinks:        quickLinks,
		ChipSources:       chipSources,
		ExpectedUpdatedAt: parsed,
	}, nil
}

func (r blacklistRequest) toInput() adminsvc.BlacklistInput {
	return adminsvc.BlacklistInput{
		Email:  r.Email,
		Reason: r.Reason,
	}
}

func parseRFC3339Timestamp(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "publishedAt is required", nil)
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339}
	var parseErr error
	for _, layout := range layouts {
		if ts, err := time.Parse(layout, trimmed); err == nil {
			return ts, nil
		} else {
			parseErr = err
		}
	}
	return time.Time{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid publishedAt", parseErr)
}

func parseISOTime(value string, field string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, fmt.Sprintf("%s is required", field), nil)
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
	}
	var parseErr error
	for _, layout := range layouts {
		if ts, err := time.Parse(layout, trimmed); err == nil {
			return ts, nil
		} else {
			parseErr = err
		}
	}
	return time.Time{}, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, fmt.Sprintf("%s must be an ISO8601 timestamp", field), parseErr)
}
