package model

import "time"

// TechLevel enumerates the declared proficiency for a technology.
type TechLevel string

const (
	TechLevelBeginner     TechLevel = "beginner"
	TechLevelIntermediate TechLevel = "intermediate"
	TechLevelAdvanced     TechLevel = "advanced"
)

// TechCatalogEntry represents a single technology descriptor within the canonical catalog.
type TechCatalogEntry struct {
	ID          uint64    `json:"id"`
	Slug        string    `json:"slug"`
	DisplayName string    `json:"displayName"`
	Category    string    `json:"category,omitempty"`
	Level       TechLevel `json:"level"`
	Icon        string    `json:"icon,omitempty"`
	SortOrder   int       `json:"sortOrder"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TechContext specifies how a technology is used in a given entity.
type TechContext string

const (
	TechContextPrimary    TechContext = "primary"
	TechContextSupporting TechContext = "supporting"
)

// TechMembership binds a technology catalog entry to a particular entity with contextual metadata.
type TechMembership struct {
	MembershipID uint64           `json:"membershipId"`
	EntityType   string           `json:"entityType"`
	EntityID     uint64           `json:"entityId"`
	Tech         TechCatalogEntry `json:"tech"`
	Context      TechContext      `json:"context"`
	Note         string           `json:"note,omitempty"`
	SortOrder    int              `json:"sortOrder"`
}

// ProfileThemeMode denotes the prefered theme mode.
type ProfileThemeMode string

const (
	ProfileThemeModeLight  ProfileThemeMode = "light"
	ProfileThemeModeDark   ProfileThemeMode = "dark"
	ProfileThemeModeSystem ProfileThemeMode = "system"
)

// ProfileTheme captures theme preferences surfaced on the SPA.
type ProfileTheme struct {
	Mode        ProfileThemeMode `json:"mode"`
	AccentColor string           `json:"accentColor,omitempty"`
}

// ProfileLab describes laboratory metadata.
type ProfileLab struct {
	Name    LocalizedText `json:"name"`
	Advisor LocalizedText `json:"advisor"`
	Room    LocalizedText `json:"room"`
	URL     string        `json:"url,omitempty"`
}

// ProfileAffiliationKind distinguishes between formal affiliations and community activities.
type ProfileAffiliationKind string

const (
	ProfileAffiliationKindAffiliation ProfileAffiliationKind = "affiliation"
	ProfileAffiliationKindCommunity   ProfileAffiliationKind = "community"
)

// ProfileAffiliation models memberships such as university departments or communities.
type ProfileAffiliation struct {
	ID          uint64                 `json:"id"`
	ProfileID   uint64                 `json:"profileId"`
	Kind        ProfileAffiliationKind `json:"kind"`
	Name        string                 `json:"name"`
	URL         string                 `json:"url,omitempty"`
	Description LocalizedText          `json:"description"`
	StartedAt   time.Time              `json:"startedAt"`
	SortOrder   int                    `json:"sortOrder"`
}

// ProfileWorkExperience captures work or research history.
type ProfileWorkExperience struct {
	ID           uint64        `json:"id"`
	ProfileID    uint64        `json:"profileId"`
	Organization LocalizedText `json:"organization"`
	Role         LocalizedText `json:"role"`
	Summary      LocalizedText `json:"summary"`
	StartedAt    time.Time     `json:"startedAt"`
	EndedAt      *time.Time    `json:"endedAt,omitempty"`
	ExternalURL  string        `json:"externalUrl,omitempty"`
	SortOrder    int           `json:"sortOrder"`
}

// ProfileSocialProvider enumerates supported social link providers.
type ProfileSocialProvider string

const (
	ProfileSocialProviderGitHub   ProfileSocialProvider = "github"
	ProfileSocialProviderZenn     ProfileSocialProvider = "zenn"
	ProfileSocialProviderLinkedIn ProfileSocialProvider = "linkedin"
	ProfileSocialProviderX        ProfileSocialProvider = "x"
	ProfileSocialProviderEmail    ProfileSocialProvider = "email"
	ProfileSocialProviderWebsite  ProfileSocialProvider = "website"
	ProfileSocialProviderOther    ProfileSocialProvider = "other"
)

// ProfileSocialLink stores social link metadata.
type ProfileSocialLink struct {
	ID        uint64                `json:"id"`
	ProfileID uint64                `json:"profileId"`
	Provider  ProfileSocialProvider `json:"provider"`
	Label     LocalizedText         `json:"label"`
	URL       string                `json:"url"`
	IsFooter  bool                  `json:"isFooter"`
	SortOrder int                   `json:"sortOrder"`
}

// ProfileTechSection groups technologies for profile display.
type ProfileTechSection struct {
	ID         uint64           `json:"id"`
	ProfileID  uint64           `json:"profileId"`
	Title      LocalizedText    `json:"title"`
	Layout     string           `json:"layout"`
	Breakpoint string           `json:"breakpoint"`
	SortOrder  int              `json:"sortOrder"`
	Members    []TechMembership `json:"members"`
}

// ProfileDocument aggregates all profile metadata needed for the public and admin experiences.
type ProfileDocument struct {
	ID           uint64                  `json:"id"`
	DisplayName  string                  `json:"displayName"`
	Headline     LocalizedText           `json:"headline"`
	Summary      LocalizedText           `json:"summary"`
	AvatarURL    string                  `json:"avatarUrl"`
	Location     LocalizedText           `json:"location"`
	Theme        ProfileTheme            `json:"theme"`
	Lab          ProfileLab              `json:"lab"`
	Affiliations []ProfileAffiliation    `json:"affiliations"`
	Communities  []ProfileAffiliation    `json:"communities"`
	WorkHistory  []ProfileWorkExperience `json:"workHistory"`
	TechSections []ProfileTechSection    `json:"techSections"`
	SocialLinks  []ProfileSocialLink     `json:"socialLinks"`
	Home         *HomePageConfigDocument `json:"home,omitempty"`
	UpdatedAt    time.Time               `json:"updatedAt"`
}

// ProjectLinkType distinguishes between the kinds of related project links.
type ProjectLinkType string

const (
	ProjectLinkTypeRepo    ProjectLinkType = "repo"
	ProjectLinkTypeDemo    ProjectLinkType = "demo"
	ProjectLinkTypeArticle ProjectLinkType = "article"
	ProjectLinkTypeSlides  ProjectLinkType = "slides"
	ProjectLinkTypeOther   ProjectLinkType = "other"
)

// ProjectLink references an external resource related to the project.
type ProjectLink struct {
	ID        uint64          `json:"id"`
	ProjectID uint64          `json:"projectId"`
	Type      ProjectLinkType `json:"type"`
	Label     LocalizedText   `json:"label"`
	URL       string          `json:"url"`
	SortOrder int             `json:"sortOrder"`
}

// ProjectPeriod captures the optional project timeline.
type ProjectPeriod struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// ProjectDocument holds the project representation backed by the new schema.
type ProjectDocument struct {
	ID            uint64           `json:"id"`
	Slug          string           `json:"slug"`
	Title         LocalizedText    `json:"title"`
	Summary       LocalizedText    `json:"summary"`
	Description   LocalizedText    `json:"description"`
	CoverImageURL string           `json:"coverImageUrl"`
	PrimaryLink   string           `json:"primaryLink"`
	Links         []ProjectLink    `json:"links"`
	Period        ProjectPeriod    `json:"period"`
	Tech          []TechMembership `json:"tech"`
	Highlight     bool             `json:"highlight"`
	Published     bool             `json:"published"`
	SortOrder     int              `json:"sortOrder"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

// ResearchKind identifies whether an entry is a research highlight or legacy blog item.
type ResearchKind string

const (
	ResearchKindResearch ResearchKind = "research"
	ResearchKindBlog     ResearchKind = "blog"
)

// ResearchLinkType denotes the category of supporting link.
type ResearchLinkType string

const (
	ResearchLinkTypePaper  ResearchLinkType = "paper"
	ResearchLinkTypeSlides ResearchLinkType = "slides"
	ResearchLinkTypeVideo  ResearchLinkType = "video"
	ResearchLinkTypeCode   ResearchLinkType = "code"
	ResearchLinkTypeOther  ResearchLinkType = "external"
)

// ResearchLink references an auxiliary asset related to the entry.
type ResearchLink struct {
	ID        uint64           `json:"id"`
	EntryID   uint64           `json:"entryId"`
	Type      ResearchLinkType `json:"type"`
	Label     LocalizedText    `json:"label"`
	URL       string           `json:"url"`
	SortOrder int              `json:"sortOrder"`
}

// ResearchAsset represents media attached to an entry.
type ResearchAsset struct {
	ID        uint64        `json:"id"`
	EntryID   uint64        `json:"entryId"`
	URL       string        `json:"url"`
	Caption   LocalizedText `json:"caption"`
	SortOrder int           `json:"sortOrder"`
}

// ResearchTag associates arbitrary keywords with an entry.
type ResearchTag struct {
	ID        uint64 `json:"id"`
	EntryID   uint64 `json:"entryId"`
	Value     string `json:"value"`
	SortOrder int    `json:"sortOrder"`
}

// ResearchDocument is the primary aggregate for research/blog content.
type ResearchDocument struct {
	ID                uint64           `json:"id"`
	Slug              string           `json:"slug"`
	Kind              ResearchKind     `json:"kind"`
	Title             LocalizedText    `json:"title"`
	Overview          LocalizedText    `json:"overview"`
	Outcome           LocalizedText    `json:"outcome"`
	Outlook           LocalizedText    `json:"outlook"`
	ExternalURL       string           `json:"externalUrl"`
	PublishedAt       time.Time        `json:"publishedAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
	HighlightImageURL string           `json:"highlightImageUrl"`
	ImageAlt          LocalizedText    `json:"imageAlt"`
	IsDraft           bool             `json:"isDraft"`
	Tags              []ResearchTag    `json:"tags"`
	Links             []ResearchLink   `json:"links"`
	Assets            []ResearchAsset  `json:"assets"`
	Tech              []TechMembership `json:"tech"`
}

// ContactTopicV2 describes a selectable topic rendered on the contact form.
type ContactTopicV2 struct {
	ID          string        `json:"id"`
	Label       LocalizedText `json:"label"`
	Description LocalizedText `json:"description"`
}

// ContactFormSettingsV2 holds the configurable attributes of the contact form/public booking experience.
type ContactFormSettingsV2 struct {
	ID                 uint64           `json:"id"`
	HeroTitle          LocalizedText    `json:"heroTitle"`
	HeroDescription    LocalizedText    `json:"heroDescription"`
	Topics             []ContactTopicV2 `json:"topics"`
	ConsentText        LocalizedText    `json:"consentText"`
	MinimumLeadHours   int              `json:"minimumLeadHours"`
	RecaptchaSiteKey   string           `json:"recaptchaSiteKey"`
	SupportEmail       string           `json:"supportEmail"`
	CalendarTimezone   string           `json:"calendarTimezone"`
	GoogleCalendarID   string           `json:"googleCalendarId"`
	BookingWindowDays  int              `json:"bookingWindowDays"`
	MeetingURLTemplate string           `json:"meetingUrlTemplate"`
	CreatedAt          time.Time        `json:"createdAt"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}

// HomeQuickLink describes hero quick links on the home screen.
type HomeQuickLink struct {
	ID          uint64        `json:"id"`
	ConfigID    uint64        `json:"configId"`
	Section     string        `json:"section"`
	Label       LocalizedText `json:"label"`
	Description LocalizedText `json:"description"`
	CTA         LocalizedText `json:"cta"`
	TargetURL   string        `json:"targetUrl"`
	SortOrder   int           `json:"sortOrder"`
}

// HomeChipSource configures dynamic chip rendering on the home screen.
type HomeChipSource struct {
	ID        uint64        `json:"id"`
	ConfigID  uint64        `json:"configId"`
	Source    string        `json:"source"`
	Label     LocalizedText `json:"label"`
	Limit     int           `json:"limit"`
	SortOrder int           `json:"sortOrder"`
}

// HomePageConfigDocument holds layout preferences for the home screen.
type HomePageConfigDocument struct {
	ID           uint64           `json:"id"`
	ProfileID    uint64           `json:"profileId"`
	HeroSubtitle LocalizedText    `json:"heroSubtitle"`
	QuickLinks   []HomeQuickLink  `json:"quickLinks"`
	ChipSources  []HomeChipSource `json:"chipSources"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

// MeetingReservationV2 mirrors meeting_reservations schema.
type MeetingReservationV2 struct {
	ID                     uint64     `json:"id"`
	Name                   string     `json:"name"`
	Email                  string     `json:"email"`
	Topic                  string     `json:"topic"`
	Message                string     `json:"message"`
	StartAt                time.Time  `json:"startAt"`
	EndAt                  time.Time  `json:"endAt"`
	DurationMinutes        int        `json:"durationMinutes"`
	GoogleEventID          string     `json:"googleEventId"`
	GoogleCalendarStatus   string     `json:"googleCalendarStatus"`
	Status                 string     `json:"status"`
	ConfirmationSentAt     *time.Time `json:"confirmationSentAt,omitempty"`
	LastNotificationSentAt *time.Time `json:"lastNotificationSentAt,omitempty"`
	LookupHash             string     `json:"lookupHash"`
	CancellationReason     string     `json:"cancellationReason"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

// MeetingNotificationV2 mirrors meeting_notifications schema.
type MeetingNotificationV2 struct {
	ID            uint64    `json:"id"`
	ReservationID uint64    `json:"reservationId"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	ErrorMessage  string    `json:"errorMessage"`
	CreatedAt     time.Time `json:"createdAt"`
}
