package model

import "time"

// AdminProject represents a project entry including administrative metadata.
type AdminProject struct {
	ID          int64         `json:"id"`
	Title       LocalizedText `json:"title"`
	Description LocalizedText `json:"description"`
	TechStack   []string      `json:"techStack"`
	LinkURL     string        `json:"linkUrl"`
	Year        int           `json:"year"`
	Published   bool          `json:"published"`
	SortOrder   *int          `json:"sortOrder,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// AdminResearch includes research content with draft management fields.
type AdminResearch struct {
	ID        int64         `json:"id"`
	Title     LocalizedText `json:"title"`
	Summary   LocalizedText `json:"summary"`
	ContentMD LocalizedText `json:"contentMd"`
	Year      int           `json:"year"`
	Published bool          `json:"published"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
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
	Status          MeetingStatus `json:"status"`
	Notes           string        `json:"notes"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
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
	PublishedProjects int `json:"publishedProjects"`
	DraftProjects     int `json:"draftProjects"`
	PublishedResearch int `json:"publishedResearch"`
	DraftResearch     int `json:"draftResearch"`
	PublishedBlogs    int `json:"publishedBlogs"`
	DraftBlogs        int `json:"draftBlogs"`
	PendingMeetings   int `json:"pendingMeetings"`
	BlacklistEntries  int `json:"blacklistEntries"`
}
