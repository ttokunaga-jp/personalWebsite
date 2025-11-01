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
	meetings := newStubMeetingRepository()
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
			Timezone:      "UTC",
			BufferMinutes: 30,
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

	svc, err := NewBookingService(meetings, availability, blacklist, calendar, mailer, cfg)
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
	require.Len(t, meetings.created, 1)
	require.Len(t, mailer.sent, 1)
	require.Equal(t, "test@example.com", mailer.sent[0].To[0])
}

func TestBookingService_Blacklist(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	meetings := newStubMeetingRepository()
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

	svc, err := NewBookingService(meetings, availability, blacklist, calendar, mailer, cfg)
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
	require.Empty(t, meetings.created)
	require.Empty(t, mailer.sent)
}

func TestBookingService_BufferConflict(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	meetings := newStubMeetingRepository()
	conflictingStart := now.Add(2 * time.Hour)
	availability := &stubAvailabilityRepository{
		windows: []model.TimeWindow{
			{
				Start: conflictingStart.Add(-15 * time.Minute),
				End:   conflictingStart.Add(15 * time.Minute),
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

	svc, err := NewBookingService(meetings, availability, blacklist, calendar, mailer, cfg)
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
	require.Empty(t, meetings.created)
}

func TestBookingService_CalendarFailure(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	meetings := newStubMeetingRepository()
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

	svc, err := NewBookingService(meetings, availability, blacklist, calendar, mailer, cfg)
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
	require.Empty(t, meetings.created)
	require.Empty(t, mailer.sent)
}

func TestBookingService_CalendarAuthRequired(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 5, 1, 9, 0, 0, 0, time.UTC)
	meetings := newStubMeetingRepository()
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

	svc, err := NewBookingService(meetings, availability, blacklist, calendar, mailer, cfg)
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
	return append([]model.TimeWindow(nil), s.busy...), nil
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

type stubMeetingRepository struct {
	created []*model.Meeting
}

func newStubMeetingRepository() *stubMeetingRepository {
	return &stubMeetingRepository{}
}

func (s *stubMeetingRepository) ListMeetings(context.Context) ([]model.Meeting, error) {
	return nil, nil
}

func (s *stubMeetingRepository) GetMeeting(context.Context, int64) (*model.Meeting, error) {
	return nil, repository.ErrNotFound
}

func (s *stubMeetingRepository) CreateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error) {
	copied := *meeting
	s.created = append(s.created, &copied)
	return &copied, nil
}

func (s *stubMeetingRepository) UpdateMeeting(context.Context, *model.Meeting) (*model.Meeting, error) {
	return nil, nil
}

func (s *stubMeetingRepository) DeleteMeeting(context.Context, int64) error {
	return nil
}

type fixedClock struct {
	now time.Time
}

func (f fixedClock) Now() time.Time {
	return f.now
}
