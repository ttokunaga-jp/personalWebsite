package model

import "time"

// AdminProfile captures editable profile metadata for the administrator UI.
type AdminProfile struct {
	Name        LocalizedText   `json:"name"`
	Title       LocalizedText   `json:"title"`
	Affiliation LocalizedText   `json:"affiliation"`
	Lab         LocalizedText   `json:"lab"`
	Summary     LocalizedText   `json:"summary"`
	Skills      []LocalizedText `json:"skills"`
	FocusAreas  []LocalizedText `json:"focusAreas"`
	UpdatedAt   *time.Time      `json:"updatedAt,omitempty"`
}

// AdminProject represents a project entry including administrative metadata.
type AdminProject struct {
	ID          int64            `json:"id"`
	Title       LocalizedText    `json:"title"`
	Description LocalizedText    `json:"description"`
	Tech        []TechMembership `json:"tech"`
	LinkURL     string           `json:"linkUrl"`
	Year        int              `json:"year"`
	Published   bool             `json:"published"`
	SortOrder   *int             `json:"sortOrder,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

// AdminResearch includes full research/blog entry content with draft management fields.
type AdminResearch struct {
	ID                uint64           `json:"id"`
	Slug              string           `json:"slug"`
	Kind              ResearchKind     `json:"kind"`
	Title             LocalizedText    `json:"title"`
	Overview          LocalizedText    `json:"overview"`
	Outcome           LocalizedText    `json:"outcome"`
	Outlook           LocalizedText    `json:"outlook"`
	ExternalURL       string           `json:"externalUrl"`
	HighlightImageURL string           `json:"highlightImageUrl"`
	ImageAlt          LocalizedText    `json:"imageAlt"`
	PublishedAt       time.Time        `json:"publishedAt"`
	IsDraft           bool             `json:"isDraft"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
	Tags              []ResearchTag    `json:"tags"`
	Links             []ResearchLink   `json:"links"`
	Assets            []ResearchAsset  `json:"assets"`
	Tech              []TechMembership `json:"tech"`
}

// BlogPost models an article managed through the admin surface.
type BlogPost struct {
	ID          int64         `json:"id"`
	Title       LocalizedText `json:"title"`
	Summary     LocalizedText `json:"summary"`
	ContentMD   LocalizedText `json:"contentMd"`
	Tags        []string      `json:"tags"`
	Published   bool          `json:"published"`
	PublishedAt *time.Time    `json:"publishedAt,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// MeetingStatus captures the lifecycle state of a meeting reservation.
type MeetingStatus string

const (
	MeetingStatusPending   MeetingStatus = "pending"
	MeetingStatusConfirmed MeetingStatus = "confirmed"
	MeetingStatusCancelled MeetingStatus = "cancelled"
)

// Meeting holds reservation data.
type Meeting struct {
	ID              int64         `json:"id"`
	Name            string        `json:"name"`
	Email           string        `json:"email"`
	Datetime        time.Time     `json:"datetime"`
	DurationMinutes int           `json:"durationMinutes"`
	MeetURL         string        `json:"meetUrl"`
	CalendarEventID string        `json:"calendarEventId"`
	Status          MeetingStatus `json:"status"`
	Notes           string        `json:"notes"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

// ContactStatus expresses the review lifecycle of a contact submission.
type ContactStatus string

const (
	ContactStatusPending  ContactStatus = "pending"
	ContactStatusInReview ContactStatus = "in_review"
	ContactStatusResolved ContactStatus = "resolved"
	ContactStatusArchived ContactStatus = "archived"
)

// ContactMessage captures incoming contact submissions enriched with moderation metadata.
type ContactMessage struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Topic     string        `json:"topic"`
	Message   string        `json:"message"`
	Status    ContactStatus `json:"status"`
	AdminNote string        `json:"adminNote"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// BlacklistEntry captures blacklisted emails that should be rejected.
type BlacklistEntry struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
}

// AdminSummary aggregates key admin metrics for dashboard display.
type AdminSummary struct {
	ProfileUpdatedAt  *time.Time `json:"profileUpdatedAt,omitempty"`
	SkillCount        int        `json:"skillCount"`
	FocusAreaCount    int        `json:"focusAreaCount"`
	PublishedProjects int        `json:"publishedProjects"`
	DraftProjects     int        `json:"draftProjects"`
	PublishedResearch int        `json:"publishedResearch"`
	DraftResearch     int        `json:"draftResearch"`
	PendingContacts   int        `json:"pendingContacts"`
	BlacklistEntries  int        `json:"blacklistEntries"`
}
