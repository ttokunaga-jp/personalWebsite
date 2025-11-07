package inmemory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type meetingReservationRepository struct {
	mu           sync.RWMutex
	seq          uint64
	reservations []model.MeetingReservation
}

// NewMeetingReservationRepository constructs an in-memory reservation repository.
func NewMeetingReservationRepository() repository.MeetingReservationRepository {
	repo := &meetingReservationRepository{
		reservations: make([]model.MeetingReservation, len(defaultMeetingReservations)),
	}
	copy(repo.reservations, defaultMeetingReservations)
	for _, reservation := range repo.reservations {
		if reservation.ID > repo.seq {
			repo.seq = reservation.ID
		}
	}
	return repo
}

func (r *meetingReservationRepository) CreateReservation(ctx context.Context, reservation *model.MeetingReservation) (*model.MeetingReservation, error) {
	if reservation == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	reservationCopy := copyReservation(*reservation)
	reservationCopy.ID = r.seq
	now := time.Now().UTC()
	if reservationCopy.CreatedAt.IsZero() {
		reservationCopy.CreatedAt = now
	}
	if reservationCopy.UpdatedAt.IsZero() {
		reservationCopy.UpdatedAt = now
	}
	r.reservations = append(r.reservations, reservationCopy)
	return copyReservationPtr(reservationCopy), nil
}

func (r *meetingReservationRepository) FindReservationByLookupHash(ctx context.Context, lookupHash string) (*model.MeetingReservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, entry := range r.reservations {
		if entry.LookupHash == lookupHash {
			return copyReservationPtr(entry), nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *meetingReservationRepository) ListConflictingReservations(ctx context.Context, start, end time.Time) ([]model.MeetingReservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var conflicts []model.MeetingReservation
	for _, entry := range r.reservations {
		if entry.Status == model.MeetingReservationStatusCancelled {
			continue
		}
		if entry.StartAt.Before(end) && entry.EndAt.After(start) {
			conflicts = append(conflicts, copyReservation(entry))
		}
	}
	return conflicts, nil
}

func (r *meetingReservationRepository) MarkConfirmationSent(ctx context.Context, id uint64, sentAt time.Time) (*model.MeetingReservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, entry := range r.reservations {
		if entry.ID != id {
			continue
		}
		if entry.Status != model.MeetingReservationStatusCancelled {
			entry.Status = model.MeetingReservationStatusConfirmed
		}
		sent := sentAt.UTC()
		entry.ConfirmationSentAt = &sent
		entry.LastNotificationSentAt = &sent
		entry.UpdatedAt = time.Now().UTC()
		r.reservations[index] = entry
		return copyReservationPtr(entry), nil
	}
	return nil, repository.ErrNotFound
}

func (r *meetingReservationRepository) CancelReservation(ctx context.Context, id uint64, reason string) (*model.MeetingReservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, entry := range r.reservations {
		if entry.ID != id {
			continue
		}
		entry.Status = model.MeetingReservationStatusCancelled
		entry.GoogleCalendarStatus = "cancelled"
		entry.CancellationReason = strings.TrimSpace(reason)
		entry.UpdatedAt = time.Now().UTC()
		r.reservations[index] = entry
		return copyReservationPtr(entry), nil
	}
	return nil, repository.ErrNotFound
}

func copyReservation(reservation model.MeetingReservation) model.MeetingReservation {
	result := reservation
	if reservation.ConfirmationSentAt != nil {
		timestamp := reservation.ConfirmationSentAt.UTC()
		result.ConfirmationSentAt = &timestamp
	}
	if reservation.LastNotificationSentAt != nil {
		timestamp := reservation.LastNotificationSentAt.UTC()
		result.LastNotificationSentAt = &timestamp
	}
	return result
}

func copyReservationPtr(reservation model.MeetingReservation) *model.MeetingReservation {
	copied := copyReservation(reservation)
	return &copied
}

type meetingNotificationRepository struct {
	mu            sync.RWMutex
	seq           uint64
	notifications map[uint64][]model.MeetingNotification
}

// NewMeetingNotificationRepository constructs an in-memory notification repository.
func NewMeetingNotificationRepository() repository.MeetingNotificationRepository {
	return &meetingNotificationRepository{
		notifications: make(map[uint64][]model.MeetingNotification),
	}
}

func (r *meetingNotificationRepository) RecordNotification(ctx context.Context, notification *model.MeetingNotification) (*model.MeetingNotification, error) {
	if notification == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	now := time.Now().UTC()
	entry := model.MeetingNotification{
		ID:            r.seq,
		ReservationID: notification.ReservationID,
		Type:          notification.Type,
		Status:        notification.Status,
		ErrorMessage:  notification.ErrorMessage,
		CreatedAt:     now,
	}

	r.notifications[notification.ReservationID] = append(r.notifications[notification.ReservationID], entry)
	return copyNotification(entry), nil
}

func (r *meetingNotificationRepository) ListNotifications(ctx context.Context, reservationID uint64) ([]model.MeetingNotification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries, ok := r.notifications[reservationID]
	if !ok {
		return []model.MeetingNotification{}, nil
	}

	result := make([]model.MeetingNotification, len(entries))
	for i, entry := range entries {
		result[i] = entry
	}
	return result, nil
}

func copyNotification(notification model.MeetingNotification) *model.MeetingNotification {
	entry := notification
	return &entry
}

var (
	_ repository.MeetingReservationRepository  = (*meetingReservationRepository)(nil)
	_ repository.MeetingNotificationRepository = (*meetingNotificationRepository)(nil)
)
