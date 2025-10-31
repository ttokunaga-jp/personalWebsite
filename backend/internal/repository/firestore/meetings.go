package firestore

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type meetingRepository struct {
	base baseRepository
}

const meetingsCollection = "meetings"

type meetingDocument struct {
	ID              int64     `firestore:"id"`
	Name            string    `firestore:"name"`
	Email           string    `firestore:"email"`
	MeetingAt       time.Time `firestore:"meetingAt"`
	DurationMinutes int       `firestore:"durationMinutes"`
	MeetURL         string    `firestore:"meetUrl"`
	CalendarEventID string    `firestore:"calendarEventId"`
	Status          string    `firestore:"status"`
	Notes           string    `firestore:"notes"`
	CreatedAt       time.Time `firestore:"createdAt"`
	UpdatedAt       time.Time `firestore:"updatedAt"`
}

func NewMeetingRepository(client *firestore.Client, prefix string) repository.MeetingRepository {
	return &meetingRepository{base: newBaseRepository(client, prefix)}
}

func (r *meetingRepository) ListMeetings(ctx context.Context) ([]model.Meeting, error) {
	docs, err := r.base.collection(meetingsCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore meetings: list: %w", err)
	}

	entries, err := r.decodeMeetings(docs)
	if err != nil {
		return nil, err
	}

	result := make([]model.Meeting, 0, len(entries))
	for _, entry := range entries {
		result = append(result, mapMeetingDocument(entry))
	}
	return result, nil
}

func (r *meetingRepository) GetMeeting(ctx context.Context, id int64) (*model.Meeting, error) {
	doc, err := r.base.doc(meetingsCollection, strconv.FormatInt(id, 10)).Get(ctx)
	if err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore meetings: get %d: %w", id, err)
	}

	var payload meetingDocument
	if err := doc.DataTo(&payload); err != nil {
		return nil, fmt.Errorf("firestore meetings: decode %d: %w", id, err)
	}
	result := mapMeetingDocument(payload)
	return &result, nil
}

func (r *meetingRepository) CreateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error) {
	if meeting == nil {
		return nil, repository.ErrInvalidInput
	}

	id, err := nextID(ctx, r.base.client, r.base.prefix, meetingsCollection)
	if err != nil {
		return nil, fmt.Errorf("firestore meetings: next id: %w", err)
	}

	now := time.Now().UTC()
	entry := meetingDocument{
		ID:              id,
		Name:            meeting.Name,
		Email:           meeting.Email,
		MeetingAt:       meeting.Datetime.UTC(),
		DurationMinutes: meeting.DurationMinutes,
		MeetURL:         meeting.MeetURL,
		CalendarEventID: meeting.CalendarEventID,
		Status:          string(meeting.Status),
		Notes:           meeting.Notes,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	docRef := r.base.doc(meetingsCollection, strconv.FormatInt(id, 10))
	if _, err := docRef.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("firestore meetings: create %d: %w", id, err)
	}

	return r.GetMeeting(ctx, id)
}

func (r *meetingRepository) UpdateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error) {
	if meeting == nil {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(meetingsCollection, strconv.FormatInt(meeting.ID, 10))
	now := time.Now().UTC()
	updates := []firestore.Update{
		{Path: "name", Value: meeting.Name},
		{Path: "email", Value: meeting.Email},
		{Path: "meetingAt", Value: meeting.Datetime.UTC()},
		{Path: "durationMinutes", Value: meeting.DurationMinutes},
		{Path: "meetUrl", Value: meeting.MeetURL},
		{Path: "calendarEventId", Value: meeting.CalendarEventID},
		{Path: "status", Value: string(meeting.Status)},
		{Path: "notes", Value: meeting.Notes},
		{Path: "updatedAt", Value: now},
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore meetings: update %d: %w", meeting.ID, err)
	}

	return r.GetMeeting(ctx, meeting.ID)
}

func (r *meetingRepository) DeleteMeeting(ctx context.Context, id int64) error {
	docRef := r.base.doc(meetingsCollection, strconv.FormatInt(id, 10))
	if _, err := docRef.Delete(ctx); err != nil {
		if notFound(err) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("firestore meetings: delete %d: %w", id, err)
	}
	return nil
}

func (r *meetingRepository) decodeMeetings(docs []*firestore.DocumentSnapshot) ([]meetingDocument, error) {
	items := make([]meetingDocument, 0, len(docs))
	for _, doc := range docs {
		var entry meetingDocument
		if err := doc.DataTo(&entry); err != nil {
			return nil, fmt.Errorf("firestore meetings: decode %s: %w", doc.Ref.ID, err)
		}
		if entry.ID == 0 {
			if id, err := strconv.ParseInt(doc.Ref.ID, 10, 64); err == nil {
				entry.ID = id
			}
		}
		items = append(items, entry)
	}

	sort.Slice(items, func(i, j int) bool {
		left, right := items[i], items[j]
		if !left.MeetingAt.Equal(right.MeetingAt) {
			return left.MeetingAt.After(right.MeetingAt)
		}
		return left.ID > right.ID
	})

	return items, nil
}

func mapMeetingDocument(doc meetingDocument) model.Meeting {
	return model.Meeting{
		ID:              doc.ID,
		Name:            doc.Name,
		Email:           doc.Email,
		Datetime:        doc.MeetingAt,
		DurationMinutes: doc.DurationMinutes,
		MeetURL:         doc.MeetURL,
		CalendarEventID: doc.CalendarEventID,
		Status:          model.MeetingStatus(doc.Status),
		Notes:           doc.Notes,
		CreatedAt:       doc.CreatedAt,
		UpdatedAt:       doc.UpdatedAt,
	}
}

var _ repository.MeetingRepository = (*meetingRepository)(nil)
