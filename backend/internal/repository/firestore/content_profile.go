package firestore

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

const (
	profileCollectionName = "profiles"
	profileDocumentKey    = "primary"
)

type profileDocumentV2 struct {
	DisplayName  string           `firestore:"displayName"`
	Headline     localizedDoc     `firestore:"headline"`
	Summary      localizedDoc     `firestore:"summary"`
	AvatarURL    string           `firestore:"avatarUrl"`
	Location     localizedDoc     `firestore:"location"`
	Theme        profileThemeDoc  `firestore:"theme"`
	Lab          profileLabDoc    `firestore:"lab"`
	Affiliations []affiliationDoc `firestore:"affiliations"`
	Communities  []affiliationDoc `firestore:"communities"`
	WorkHistory  []workHistoryDoc `firestore:"workHistory"`
	TechSections []techSectionDoc `firestore:"techSections"`
	SocialLinks  []socialLinkDoc  `firestore:"socialLinks"`
	UpdatedAt    time.Time        `firestore:"updatedAt"`
}

type profileThemeDoc struct {
	Mode        string `firestore:"mode"`
	AccentColor string `firestore:"accentColor"`
}

type profileLabDoc struct {
	Name    localizedDoc `firestore:"name"`
	Advisor localizedDoc `firestore:"advisor"`
	Room    localizedDoc `firestore:"room"`
	URL     string       `firestore:"url"`
}

type affiliationDoc struct {
	ID          int64        `firestore:"id"`
	Kind        string       `firestore:"kind"`
	Name        string       `firestore:"name"`
	URL         string       `firestore:"url"`
	Description localizedDoc `firestore:"description"`
	StartedAt   time.Time    `firestore:"startedAt"`
	SortOrder   int          `firestore:"sortOrder"`
}

type workHistoryDoc struct {
	ID           int64        `firestore:"id"`
	Organization localizedDoc `firestore:"organization"`
	Role         localizedDoc `firestore:"role"`
	Summary      localizedDoc `firestore:"summary"`
	StartedAt    time.Time    `firestore:"startedAt"`
	EndedAt      *time.Time   `firestore:"endedAt"`
	ExternalURL  string       `firestore:"externalUrl"`
	SortOrder    int          `firestore:"sortOrder"`
}

type techSectionDoc struct {
	ID         int64               `firestore:"id"`
	Title      localizedDoc        `firestore:"title"`
	Layout     string              `firestore:"layout"`
	Breakpoint string              `firestore:"breakpoint"`
	SortOrder  int                 `firestore:"sortOrder"`
	Members    []techMembershipDoc `firestore:"members"`
}

type techMembershipDoc struct {
	MembershipID int64          `firestore:"membershipId"`
	Tech         techCatalogDoc `firestore:"tech"`
	Context      string         `firestore:"context"`
	Note         string         `firestore:"note"`
	SortOrder    int            `firestore:"sortOrder"`
}

type techCatalogDoc struct {
	ID          int64     `firestore:"id"`
	Slug        string    `firestore:"slug"`
	DisplayName string    `firestore:"displayName"`
	Category    string    `firestore:"category"`
	Level       string    `firestore:"level"`
	Icon        string    `firestore:"icon"`
	SortOrder   int       `firestore:"sortOrder"`
	Active      bool      `firestore:"active"`
	CreatedAt   time.Time `firestore:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt"`
}

type socialLinkDoc struct {
	ID        int64        `firestore:"id"`
	Provider  string       `firestore:"provider"`
	Label     localizedDoc `firestore:"label"`
	URL       string       `firestore:"url"`
	IsFooter  bool         `firestore:"isFooter"`
	SortOrder int          `firestore:"sortOrder"`
}

type contentProfileRepository struct {
	base baseRepository
}

// NewContentProfileRepository creates a Firestore-backed ContentProfileRepository.
func NewContentProfileRepository(client *firestore.Client, prefix string) repository.ContentProfileRepository {
	return &contentProfileRepository{base: newBaseRepository(client, prefix)}
}

func (r *contentProfileRepository) GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error) {
	docRef := r.base.doc(profileCollectionName, profileDocumentKey)
	snapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("firestore profile document %s: %w", profileDocumentKey, err)
	}

	var doc profileDocumentV2
	if err := snapshot.DataTo(&doc); err != nil {
		return nil, fmt.Errorf("firestore decode profile document: %w", err)
	}

	return convertProfileDocument(doc), nil
}

func convertProfileDocument(doc profileDocumentV2) *model.ProfileDocument {
	affiliations := make([]model.ProfileAffiliation, 0, len(doc.Affiliations))
	for _, a := range doc.Affiliations {
		affiliations = append(affiliations, model.ProfileAffiliation{
			ID:          uint64(a.ID),
			Kind:        model.ProfileAffiliationKind(strings.TrimSpace(a.Kind)),
			Name:        strings.TrimSpace(a.Name),
			URL:         strings.TrimSpace(a.URL),
			Description: fromLocalizedDoc(a.Description),
			StartedAt:   a.StartedAt,
			SortOrder:   a.SortOrder,
		})
	}

	communities := make([]model.ProfileAffiliation, 0, len(doc.Communities))
	for _, c := range doc.Communities {
		communities = append(communities, model.ProfileAffiliation{
			ID:          uint64(c.ID),
			Kind:        model.ProfileAffiliationKind(strings.TrimSpace(c.Kind)),
			Name:        strings.TrimSpace(c.Name),
			URL:         strings.TrimSpace(c.URL),
			Description: fromLocalizedDoc(c.Description),
			StartedAt:   c.StartedAt,
			SortOrder:   c.SortOrder,
		})
	}

	history := make([]model.ProfileWorkExperience, 0, len(doc.WorkHistory))
	for _, entry := range doc.WorkHistory {
		history = append(history, model.ProfileWorkExperience{
			ID:           uint64(entry.ID),
			Organization: fromLocalizedDoc(entry.Organization),
			Role:         fromLocalizedDoc(entry.Role),
			Summary:      fromLocalizedDoc(entry.Summary),
			StartedAt:    entry.StartedAt,
			EndedAt:      entry.EndedAt,
			ExternalURL:  strings.TrimSpace(entry.ExternalURL),
			SortOrder:    entry.SortOrder,
		})
	}

	sections := make([]model.ProfileTechSection, 0, len(doc.TechSections))
	for _, section := range doc.TechSections {
		members := make([]model.TechMembership, 0, len(section.Members))
		for _, member := range section.Members {
			members = append(members, model.TechMembership{
				MembershipID: uint64(member.MembershipID),
				EntityType:   "profile_section",
				EntityID:     uint64(section.ID),
				Tech: model.TechCatalogEntry{
					ID:          uint64(member.Tech.ID),
					Slug:        strings.TrimSpace(member.Tech.Slug),
					DisplayName: strings.TrimSpace(member.Tech.DisplayName),
					Category:    strings.TrimSpace(member.Tech.Category),
					Level:       model.TechLevel(strings.TrimSpace(member.Tech.Level)),
					Icon:        strings.TrimSpace(member.Tech.Icon),
					SortOrder:   member.Tech.SortOrder,
					Active:      member.Tech.Active,
					CreatedAt:   member.Tech.CreatedAt,
					UpdatedAt:   member.Tech.UpdatedAt,
				},
				Context:   model.TechContext(strings.TrimSpace(member.Context)),
				Note:      strings.TrimSpace(member.Note),
				SortOrder: member.SortOrder,
			})
		}
		sections = append(sections, model.ProfileTechSection{
			ID:         uint64(section.ID),
			Title:      fromLocalizedDoc(section.Title),
			Layout:     strings.TrimSpace(section.Layout),
			Breakpoint: strings.TrimSpace(section.Breakpoint),
			SortOrder:  section.SortOrder,
			Members:    members,
		})
	}

	socialLinks := make([]model.ProfileSocialLink, 0, len(doc.SocialLinks))
	for _, link := range doc.SocialLinks {
		socialLinks = append(socialLinks, model.ProfileSocialLink{
			ID:        uint64(link.ID),
			Provider:  model.ProfileSocialProvider(strings.TrimSpace(link.Provider)),
			Label:     fromLocalizedDoc(link.Label),
			URL:       strings.TrimSpace(link.URL),
			IsFooter:  link.IsFooter,
			SortOrder: link.SortOrder,
		})
	}

	return &model.ProfileDocument{
		DisplayName: strings.TrimSpace(doc.DisplayName),
		Headline:    fromLocalizedDoc(doc.Headline),
		Summary:     fromLocalizedDoc(doc.Summary),
		AvatarURL:   strings.TrimSpace(doc.AvatarURL),
		Location:    fromLocalizedDoc(doc.Location),
		Theme: model.ProfileTheme{
			Mode:        model.ProfileThemeMode(strings.TrimSpace(doc.Theme.Mode)),
			AccentColor: strings.TrimSpace(doc.Theme.AccentColor),
		},
		Lab: model.ProfileLab{
			Name:    fromLocalizedDoc(doc.Lab.Name),
			Advisor: fromLocalizedDoc(doc.Lab.Advisor),
			Room:    fromLocalizedDoc(doc.Lab.Room),
			URL:     strings.TrimSpace(doc.Lab.URL),
		},
		Affiliations: affiliations,
		Communities:  communities,
		WorkHistory:  history,
		TechSections: sections,
		SocialLinks:  socialLinks,
		UpdatedAt:    doc.UpdatedAt,
	}
}

var _ repository.ContentProfileRepository = (*contentProfileRepository)(nil)
