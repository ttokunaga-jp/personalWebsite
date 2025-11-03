package repository

import (
	"context"

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

// AdminProfileRepository manages author profile metadata.
type AdminProfileRepository interface {
	GetAdminProfile(ctx context.Context) (*model.AdminProfile, error)
	UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error)
}

// AdminResearchRepository manages research CRUD operations for the admin surface.
type AdminResearchRepository interface {
	ListAdminResearch(ctx context.Context) ([]model.AdminResearch, error)
	GetAdminResearch(ctx context.Context, id int64) (*model.AdminResearch, error)
	CreateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error)
	UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error)
	DeleteAdminResearch(ctx context.Context, id int64) error
}

// BlogRepository manages administrator blog CRUD.
type BlogRepository interface {
	ListBlogPosts(ctx context.Context) ([]model.BlogPost, error)
	GetBlogPost(ctx context.Context, id int64) (*model.BlogPost, error)
	CreateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error)
	UpdateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error)
	DeleteBlogPost(ctx context.Context, id int64) error
}

// MeetingRepository manages reservations.
type MeetingRepository interface {
	ListMeetings(ctx context.Context) ([]model.Meeting, error)
	GetMeeting(ctx context.Context, id int64) (*model.Meeting, error)
	CreateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error)
	UpdateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error)
	DeleteMeeting(ctx context.Context, id int64) error
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
