package repository

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
)

// AdminProjectRepository manages project CRUD operations for the admin surface.
type AdminProjectRepository interface {
	ListAdminProjects(ctx context.Context) ([]model.AdminProject, error)
	GetAdminProject(ctx context.Context, id int64) (*model.AdminProject, error)
	CreateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error)
	UpdateAdminProject(ctx context.Context, project *model.AdminProject) (*model.AdminProject, error)
	DeleteAdminProject(ctx context.Context, id int64) error
}

// AdminSessionRepository persists and retrieves administrator sessions.
type AdminSessionRepository interface {
	CreateSession(ctx context.Context, session *model.AdminSession) (*model.AdminSession, error)
	FindSessionByHash(ctx context.Context, hash string) (*model.AdminSession, error)
	UpdateSessionActivity(ctx context.Context, hash string, lastAccessed time.Time, expiresAt time.Time) (*model.AdminSession, error)
	RevokeSession(ctx context.Context, hash string) error
}

// AdminProfileRepository manages author profile metadata.
type AdminProfileRepository interface {
	GetAdminProfile(ctx context.Context) (*model.AdminProfile, error)
	UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error)
}

// AdminResearchRepository manages research CRUD operations for the admin surface.
type AdminResearchRepository interface {
	ListAdminResearch(ctx context.Context) ([]model.AdminResearch, error)
	GetAdminResearch(ctx context.Context, id uint64) (*model.AdminResearch, error)
	CreateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error)
	UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error)
	DeleteAdminResearch(ctx context.Context, id uint64) error
}

// BlogRepository manages administrator blog CRUD.
type BlogRepository interface {
	ListBlogPosts(ctx context.Context) ([]model.BlogPost, error)
	GetBlogPost(ctx context.Context, id int64) (*model.BlogPost, error)
	CreateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error)
	UpdateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error)
	DeleteBlogPost(ctx context.Context, id int64) error
}

// MeetingReservationRepository manages reservations backed by meeting_reservations.
type MeetingReservationRepository interface {
	CreateReservation(ctx context.Context, reservation *model.MeetingReservation) (*model.MeetingReservation, error)
	FindReservationByLookupHash(ctx context.Context, lookupHash string) (*model.MeetingReservation, error)
	ListConflictingReservations(ctx context.Context, start, end time.Time) ([]model.MeetingReservation, error)
	MarkConfirmationSent(ctx context.Context, id uint64, sentAt time.Time) (*model.MeetingReservation, error)
	CancelReservation(ctx context.Context, id uint64, reason string) (*model.MeetingReservation, error)
}

// MeetingNotificationRepository records outgoing notifications linked to reservations.
type MeetingNotificationRepository interface {
	RecordNotification(ctx context.Context, notification *model.MeetingNotification) (*model.MeetingNotification, error)
	ListNotifications(ctx context.Context, reservationID uint64) ([]model.MeetingNotification, error)
}

// AdminContactRepository exposes management capabilities for contact submissions.
type AdminContactRepository interface {
	ListContactMessages(ctx context.Context) ([]model.ContactMessage, error)
	GetContactMessage(ctx context.Context, id string) (*model.ContactMessage, error)
	UpdateContactMessage(ctx context.Context, message *model.ContactMessage) (*model.ContactMessage, error)
	DeleteContactMessage(ctx context.Context, id string) error
}

// BlacklistRepository persists blacklisted emails for booking exclusion.
type BlacklistRepository interface {
	ListBlacklistEntries(ctx context.Context) ([]model.BlacklistEntry, error)
	AddBlacklistEntry(ctx context.Context, entry *model.BlacklistEntry) (*model.BlacklistEntry, error)
	UpdateBlacklistEntry(ctx context.Context, entry *model.BlacklistEntry) (*model.BlacklistEntry, error)
	RemoveBlacklistEntry(ctx context.Context, id int64) error
	FindBlacklistEntryByEmail(ctx context.Context, email string) (*model.BlacklistEntry, error)
}
