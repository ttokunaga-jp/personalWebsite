package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/calendar"
	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/infra/google"
	"github.com/takumi/personal-website/internal/mail"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

func TestBookingService_Success(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	reservations := newStubReservationRepository()
	notifications := newStubNotificationRepository()
	availability := &stubAvailabilityRepository{}
	blacklist := &stubBlacklistRepository{}
	calendar := &stubCalendarClient{
		event: &calendar.Event{
			ID:          "evt-123",
			HangoutLink: "https://meet.example.com/evt-123",
		},
	}
	mailer := &stubMailClient{}

	cfg := &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:         "UTC",
			BufferMinutes:    30,
			SupportEmail:     "support@example.com",
			CalendarTimezone: "UTC",
		},
		Booking: config.BookingConfig{
			CalendarID:           "primary",
			MeetTemplate:         "Portfolio Intro Session",
			NotificationSender:   "noreply@example.com",
			NotificationReceiver: "owner@example.com",
			MaxRetries:           2,
			RequestTimeout:       200 * time.Millisecond,
			InitialBackoff:       10 * time.Millisecond,
			BackoffMultiplier:    1.0,
		},
	}

	svc, err := NewBookingService(reservations, notifications, availability, blacklist, calendar, mailer, cfg)
	require.NoError(t, err)
	svc.(*bookingService).clock = fixedClock{now: now}

	result, err := svc.Book(context.Background(), model.BookingRequest{
		Name:            "Test User",
		Email:           "test@example.com",
		StartTime:       now.Add(2 * time.Hour),
		DurationMinutes: 45,
		Agenda:          "Discuss portfolio improvements",
		RecaptchaToken:  "test-token",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "evt-123", result.CalendarEventID)
	require.NotEmpty(t, result.Reservation.LookupHash)
	require.Equal(t, "support@example.com", result.SupportEmail)
	require.Equal(t, "UTC", result.CalendarTimezone)
	require.Len(t, reservations.created, 1)
	require.Len(t, notifications.recorded, 1)
	require.Len(t, mailer.sent, 1)
	require.Equal(t, "test@example.com", mailer.sent[0].To[0])
}

func TestBookingService_LookupReservation(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	reservations := newStubReservationRepository()
	notifications := newStubNotificationRepository()
	reservations.seq = 1
	reservations.entries[1] = &model.MeetingReservation{
		ID:              1,
		LookupHash:      "lookup-hash",
		Name:            "Existing User",
		Email:           "existing@example.com",
		StartAt:         now,
		EndAt:           now.Add(30 * time.Minute),
		DurationMinutes: 30,
		GoogleEventID:   "evt-existing",
		Status:          model.MeetingReservationStatusConfirmed,
	}
	reservations.index["lookup-hash"] = 1

	cfg := &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:         "UTC",
			CalendarTimezone: "UTC",
			SupportEmail:     "support@example.com",
		},
		Booking: config.BookingConfig{
			CalendarID:         "primary",
			NotificationSender: "noreply@example.com",
		},
	}

	svc, err := NewBookingService(reservations, notifications, &stubAvailabilityRepository{}, &stubBlacklistRepository{}, &stubCalendarClient{}, &stubMailClient{}, cfg)
	require.NoError(t, err)

	result, err := svc.LookupReservation(context.Background(), "lookup-hash")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "lookup-hash", result.Reservation.LookupHash)
	require.Equal(t, "evt-existing", result.CalendarEventID)
	require.Equal(t, "support@example.com", result.SupportEmail)
}

func TestBookingService_Blacklist(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	reservations := newStubReservationRepository()
	notifications := newStubNotificationRepository()
	availability := &stubAvailabilityRepository{}
	blacklist := &stubBlacklistRepository{
		blocked: map[string]bool{
			"blocked@example.com": true,
		},
	}
	calendar := &stubCalendarClient{
		event: &calendar.Event{ID: "evt-123"},
	}
	mailer := &stubMailClient{}

	cfg := &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:      "UTC",
			BufferMinutes: 30,
		},
		Booking: config.BookingConfig{
			CalendarID:         "primary",
			MaxRetries:         1,
			RequestTimeout:     100 * time.Millisecond,
			InitialBackoff:     10 * time.Millisecond,
			BackoffMultiplier:  1,
			MeetTemplate:       "Portfolio Intro Session",
			NotificationSender: "noreply@example.com",
		},
	}

	svc, err := NewBookingService(reservations, notifications, availability, blacklist, calendar, mailer, cfg)
	require.NoError(t, err)
	svc.(*bookingService).clock = fixedClock{now: now}

	_, err = svc.Book(context.Background(), model.BookingRequest{
		Name:            "Blocked User",
		Email:           "blocked@example.com",
		StartTime:       now.Add(2 * time.Hour),
		DurationMinutes: 30,
		RecaptchaToken:  "test-token",
	})
	require.Error(t, err)
	appErr := errs.From(err)
	require.Equal(t, http.StatusForbidden, appErr.Status)
	require.Empty(t, reservations.created)
	require.Empty(t, mailer.sent)
}

func TestBookingService_BufferConflict(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	reservations := newStubReservationRepository()
	notifications := newStubNotificationRepository()
	conflictingStart := now.Add(2 * time.Hour)
	availability := &stubAvailabilityRepository{
		windows: []model.TimeWindow{
			{
				Start:  conflictingStart.Add(-15 * time.Minute),
				End:    conflictingStart.Add(15 * time.Minute),
				Source: model.BusyWindowSourceReservation,
			},
		},
	}
	blacklist := &stubBlacklistRepository{}
	calendar := &stubCalendarClient{
		event: &calendar.Event{ID: "evt-456"},
	}
	mailer := &stubMailClient{}

	cfg := &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:      "UTC",
			BufferMinutes: 30,
		},
		Booking: config.BookingConfig{
			CalendarID:         "primary",
			MaxRetries:         1,
			RequestTimeout:     100 * time.Millisecond,
			InitialBackoff:     10 * time.Millisecond,
			BackoffMultiplier:  1,
			MeetTemplate:       "Portfolio Intro Session",
			NotificationSender: "noreply@example.com",
		},
	}

	svc, err := NewBookingService(reservations, notifications, availability, blacklist, calendar, mailer, cfg)
	require.NoError(t, err)
	svc.(*bookingService).clock = fixedClock{now: now}

	_, err = svc.Book(context.Background(), model.BookingRequest{
		Name:            "Conflict User",
		Email:           "conflict@example.com",
		StartTime:       conflictingStart,
		DurationMinutes: 30,
		RecaptchaToken:  "test-token",
	})
	require.Error(t, err)
	appErr := errs.From(err)
	require.Equal(t, http.StatusConflict, appErr.Status)
	require.Empty(t, reservations.created)
}

func TestBookingService_CalendarFailure(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	reservations := newStubReservationRepository()
	notifications := newStubNotificationRepository()
	availability := &stubAvailabilityRepository{}
	blacklist := &stubBlacklistRepository{}
	calendar := &stubCalendarClient{
		event:     &calendar.Event{ID: "evt-789"},
		createErr: errors.New("calendar unavailable"),
	}
	mailer := &stubMailClient{}

	cfg := &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:      "UTC",
			BufferMinutes: 30,
		},
		Booking: config.BookingConfig{
			CalendarID:         "primary",
			MaxRetries:         2,
			RequestTimeout:     20 * time.Millisecond,
			InitialBackoff:     5 * time.Millisecond,
			BackoffMultiplier:  1.0,
			MeetTemplate:       "Portfolio Intro Session",
			NotificationSender: "noreply@example.com",
		},
	}

	svc, err := NewBookingService(reservations, notifications, availability, blacklist, calendar, mailer, cfg)
	require.NoError(t, err)
	svc.(*bookingService).clock = fixedClock{now: now}

	_, err = svc.Book(context.Background(), model.BookingRequest{
		Name:            "Retry User",
		Email:           "retry@example.com",
		StartTime:       now.Add(2 * time.Hour),
		DurationMinutes: 30,
		RecaptchaToken:  "test-token",
	})
	require.Error(t, err)
	appErr := errs.From(err)
	require.Equal(t, http.StatusBadGateway, appErr.Status)
	require.Equal(t, 2, calendar.createCalls)
	require.Empty(t, reservations.created)
	require.Empty(t, mailer.sent)
}

func TestBookingService_CalendarAuthRequired(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	reservations := newStubReservationRepository()
	notifications := newStubNotificationRepository()
	availability := &stubAvailabilityRepository{}
	blacklist := &stubBlacklistRepository{}
	calendar := &stubCalendarClient{
		listErr: google.ErrTokenNotFound,
	}
	mailer := &stubMailClient{}

	cfg := &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:      "UTC",
			BufferMinutes: 30,
		},
		Booking: config.BookingConfig{
			CalendarID:         "primary",
			MaxRetries:         1,
			RequestTimeout:     20 * time.Millisecond,
			InitialBackoff:     5 * time.Millisecond,
			BackoffMultiplier:  1.0,
			MeetTemplate:       "Portfolio Intro Session",
			NotificationSender: "noreply@example.com",
		},
	}

	svc, err := NewBookingService(reservations, notifications, availability, blacklist, calendar, mailer, cfg)
	require.NoError(t, err)
	svc.(*bookingService).clock = fixedClock{now: now}

	_, err = svc.Book(context.Background(), model.BookingRequest{
		Name:            "NeedsAuth",
		Email:           "auth@example.com",
		StartTime:       now.Add(2 * time.Hour),
		DurationMinutes: 30,
		RecaptchaToken:  "test-token",
	})
	require.Error(t, err)
	appErr := errs.From(err)
	require.Equal(t, http.StatusServiceUnavailable, appErr.Status)
	require.Contains(t, strings.ToLower(appErr.Message), "authorization")
}

type stubAvailabilityRepository struct {
	windows []model.TimeWindow
}

func (s *stubAvailabilityRepository) ListBusyWindows(context.Context, time.Time, time.Time) ([]model.TimeWindow, error) {
	return append([]model.TimeWindow(nil), s.windows...), nil
}

type stubBlacklistRepository struct {
	blocked map[string]bool
}

func (s *stubBlacklistRepository) ListBlacklistEntries(context.Context) ([]model.BlacklistEntry, error) {
	return nil, nil
}

func (s *stubBlacklistRepository) AddBlacklistEntry(context.Context, *model.BlacklistEntry) (*model.BlacklistEntry, error) {
	return nil, nil
}

func (s *stubBlacklistRepository) UpdateBlacklistEntry(context.Context, *model.BlacklistEntry) (*model.BlacklistEntry, error) {
	return nil, nil
}

func (s *stubBlacklistRepository) RemoveBlacklistEntry(context.Context, int64) error {
	return nil
}

func (s *stubBlacklistRepository) FindBlacklistEntryByEmail(ctx context.Context, email string) (*model.BlacklistEntry, error) {
	if s.blocked != nil && s.blocked[email] {
		return &model.BlacklistEntry{Email: email}, nil
	}
	return nil, repository.ErrNotFound
}

type stubCalendarClient struct {
	busy        []model.TimeWindow
	event       *calendar.Event
	listErr     error
	createErr   error
	listCalls   int
	createCalls int
}

func (s *stubCalendarClient) ListBusyWindows(context.Context, string, time.Time, time.Time) ([]model.TimeWindow, error) {
	s.listCalls++
	if s.listErr != nil {
		return nil, s.listErr
	}
	result := make([]model.TimeWindow, len(s.busy))
	for i, window := range s.busy {
		if window.Source == "" {
			window.Source = model.BusyWindowSourceExternal
		}
		result[i] = window
	}
	return result, nil
}

func (s *stubCalendarClient) CreateEvent(context.Context, string, calendar.EventInput) (*calendar.Event, error) {
	s.createCalls++
	if s.createErr != nil {
		return nil, s.createErr
	}
	return s.event, nil
}

type stubMailClient struct {
	sent []mail.Message
	err  error
}

func (s *stubMailClient) Send(ctx context.Context, message mail.Message) error {
	if s.err != nil {
		return s.err
	}
	s.sent = append(s.sent, message)
	return nil
}

type stubReservationRepository struct {
	created   []*model.MeetingReservation
	conflicts []model.MeetingReservation
	entries   map[uint64]*model.MeetingReservation
	index     map[string]uint64
	seq       uint64
	markErr   error
	cancelErr error
}

func newStubReservationRepository() *stubReservationRepository {
	return &stubReservationRepository{
		entries: make(map[uint64]*model.MeetingReservation),
		index:   make(map[string]uint64),
	}
}

func (s *stubReservationRepository) CreateReservation(ctx context.Context, reservation *model.MeetingReservation) (*model.MeetingReservation, error) {
	s.seq++
	stored := *reservation
	stored.ID = s.seq
	s.entries[stored.ID] = &stored
	s.index[stored.LookupHash] = stored.ID

	createdCopy := stored
	s.created = append(s.created, &createdCopy)
	return cloneReservation(&stored), nil
}

func (s *stubReservationRepository) FindReservationByLookupHash(ctx context.Context, lookupHash string) (*model.MeetingReservation, error) {
	id, ok := s.index[lookupHash]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneReservation(s.entries[id]), nil
}

func (s *stubReservationRepository) ListConflictingReservations(ctx context.Context, start, end time.Time) ([]model.MeetingReservation, error) {
	conflicts := make([]model.MeetingReservation, len(s.conflicts))
	copy(conflicts, s.conflicts)
	return conflicts, nil
}

func (s *stubReservationRepository) MarkConfirmationSent(ctx context.Context, id uint64, sentAt time.Time) (*model.MeetingReservation, error) {
	if s.markErr != nil {
		return nil, s.markErr
	}
	entry, ok := s.entries[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	ts := sentAt.UTC()
	entry.Status = model.MeetingReservationStatusConfirmed
	entry.ConfirmationSentAt = &ts
	entry.LastNotificationSentAt = &ts
	entry.UpdatedAt = ts
	s.entries[id] = entry
	return cloneReservation(entry), nil
}

func (s *stubReservationRepository) CancelReservation(ctx context.Context, id uint64, reason string) (*model.MeetingReservation, error) {
	if s.cancelErr != nil {
		return nil, s.cancelErr
	}
	entry, ok := s.entries[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	entry.Status = model.MeetingReservationStatusCancelled
	entry.CancellationReason = reason
	entry.GoogleCalendarStatus = "cancelled"
	entry.UpdatedAt = time.Now().UTC()
	s.entries[id] = entry
	return cloneReservation(entry), nil
}

func cloneReservation(reservation *model.MeetingReservation) *model.MeetingReservation {
	if reservation == nil {
		return nil
	}
	copy := *reservation
	if reservation.ConfirmationSentAt != nil {
		ts := reservation.ConfirmationSentAt.UTC()
		copy.ConfirmationSentAt = &ts
	}
	if reservation.LastNotificationSentAt != nil {
		ts := reservation.LastNotificationSentAt.UTC()
		copy.LastNotificationSentAt = &ts
	}
	return &copy
}

type stubNotificationRepository struct {
	recorded []model.MeetingNotification
	err      error
}

func newStubNotificationRepository() *stubNotificationRepository {
	return &stubNotificationRepository{}
}

func (s *stubNotificationRepository) RecordNotification(ctx context.Context, notification *model.MeetingNotification) (*model.MeetingNotification, error) {
	if s.err != nil {
		return nil, s.err
	}
	entry := *notification
	entry.ID = uint64(len(s.recorded) + 1)
	entry.CreatedAt = time.Now().UTC()
	s.recorded = append(s.recorded, entry)
	return &entry, nil
}

func (s *stubNotificationRepository) ListNotifications(ctx context.Context, reservationID uint64) ([]model.MeetingNotification, error) {
	result := make([]model.MeetingNotification, 0, len(s.recorded))
	for _, entry := range s.recorded {
		if entry.ReservationID == reservationID {
			result = append(result, entry)
		}
	}
	return result, nil
}

type fixedClock struct {
	now time.Time
}

func (f fixedClock) Now() time.Time {
	return f.now
}
