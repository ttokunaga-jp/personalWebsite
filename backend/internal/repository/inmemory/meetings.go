package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type meetingRepository struct {
	mu      sync.RWMutex
	entries []model.Meeting
	nextID  int64
}

func NewMeetingRepository() repository.MeetingRepository {
	entries := make([]model.Meeting, len(defaultMeetings))
	copy(entries, defaultMeetings)
	var maxID int64
	for _, m := range entries {
		if m.ID > maxID {
			maxID = m.ID
		}
	}
	return &meetingRepository{
		entries: entries,
		nextID:  maxID + 1,
	}
}

func (r *meetingRepository) ListMeetings(ctx context.Context) ([]model.Meeting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Meeting, len(r.entries))
	for i, entry := range r.entries {
		result[i] = copyMeeting(entry)
	}
	return result, nil
}

func (r *meetingRepository) GetMeeting(ctx context.Context, id int64) (*model.Meeting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, entry := range r.entries {
		if entry.ID == id {
			copied := copyMeeting(entry)
			return &copied, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *meetingRepository) CreateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error) {
	if meeting == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	meeting.ID = r.nextID
	r.nextID++
	meeting.CreatedAt = now
	meeting.UpdatedAt = now

	r.entries = append(r.entries, copyMeeting(*meeting))
	created := copyMeeting(*meeting)
	return &created, nil
}

func (r *meetingRepository) UpdateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error) {
	if meeting == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, existing := range r.entries {
		if existing.ID == meeting.ID {
			meeting.CreatedAt = existing.CreatedAt
			meeting.UpdatedAt = time.Now().UTC()
			r.entries[idx] = copyMeeting(*meeting)
			updated := copyMeeting(*meeting)
			return &updated, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *meetingRepository) DeleteMeeting(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, existing := range r.entries {
		if existing.ID == id {
			r.entries = append(r.entries[:idx], r.entries[idx+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func copyMeeting(src model.Meeting) model.Meeting {
	return src
}

var _ repository.MeetingRepository = (*meetingRepository)(nil)
