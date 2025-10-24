package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	mailpkg "net/mail"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/calendar"
	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/infra/google"
	"github.com/takumi/personal-website/internal/mail"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// BookingService coordinates scheduling between Google Calendar, persistence, and notifications.
type BookingService interface {
	Book(ctx context.Context, req model.BookingRequest) (*model.BookingResult, error)
}

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

type bookingService struct {
	meetings     repository.MeetingRepository
	availability repository.AvailabilityRepository
	blacklist    repository.BlacklistRepository
	calendar     calendar.Client
	mailer       mail.Client
	cfg          config.BookingConfig
	contactCfg   config.ContactConfig
	calendarCB   *circuitBreaker
	mailCB       *circuitBreaker
	clock        Clock
}

func NewBookingService(
	meetings repository.MeetingRepository,
	availability repository.AvailabilityRepository,
	blacklist repository.BlacklistRepository,
	calendar calendar.Client,
	mailer mail.Client,
	cfg *config.AppConfig,
) (BookingService, error) {
	if meetings == nil || availability == nil || blacklist == nil || calendar == nil || mailer == nil || cfg == nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "booking service: missing dependencies", nil)
	}

	bookingCfg := cfg.Booking
	if bookingCfg.MaxRetries <= 0 {
		bookingCfg.MaxRetries = 3
	}
	if bookingCfg.InitialBackoff <= 0 {
		bookingCfg.InitialBackoff = 750 * time.Millisecond
	}
	if bookingCfg.BackoffMultiplier < 1.1 {
		bookingCfg.BackoffMultiplier = 2.0
	}
	if bookingCfg.RequestTimeout <= 0 {
		bookingCfg.RequestTimeout = 8 * time.Second
	}
	if bookingCfg.CircuitFailureThresh <= 0 {
		bookingCfg.CircuitFailureThresh = 3
	}
	if bookingCfg.CircuitOpenSeconds <= 0 {
		bookingCfg.CircuitOpenSeconds = 60
	}

	openDuration := time.Duration(bookingCfg.CircuitOpenSeconds) * time.Second

	return &bookingService{
		meetings:     meetings,
		availability: availability,
		blacklist:    blacklist,
		calendar:     calendar,
		mailer:       mailer,
		cfg:          bookingCfg,
		contactCfg:   cfg.Contact,
		calendarCB:   newCircuitBreaker(bookingCfg.CircuitFailureThresh, openDuration),
		mailCB:       newCircuitBreaker(bookingCfg.CircuitFailureThresh, openDuration),
		clock:        realClock{},
	}, nil
}

func (s *bookingService) Book(ctx context.Context, req model.BookingRequest) (*model.BookingResult, error) {
	if err := validateBookingRequest(req); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	agenda := strings.TrimSpace(req.Agenda)

	loc, err := time.LoadLocation(s.contactCfg.Timezone)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "booking service: invalid timezone configuration", err)
	}

	now := s.clock.Now().In(loc)
	startLocal := req.StartTime.In(loc)
	if startLocal.Before(now.Add(15 * time.Minute)) {
		return nil, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "reservation must be at least 15 minutes in the future", nil)
	}

	duration := time.Duration(req.DurationMinutes) * time.Minute
	if duration <= 0 {
		return nil, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "duration must be greater than zero", nil)
	}

	bufferMinutes := s.contactCfg.BufferMinutes
	if bufferMinutes <= 0 {
		bufferMinutes = 30
	}
	buffer := time.Duration(bufferMinutes) * time.Minute

	windowStart := startLocal.Add(-buffer)
	windowEnd := startLocal.Add(duration + buffer)

	if _, err := s.blacklist.FindBlacklistEntryByEmail(ctx, email); err == nil {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusForbidden, "email address is blocked from scheduling", nil)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to validate blacklist status", err)
	}

	busyWindows, err := s.availability.ListBusyWindows(ctx, windowStart.UTC(), windowEnd.UTC())
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load local busy windows", err)
	}

	var externalBusy []model.TimeWindow
	err = s.withRetry(ctx, s.calendarCB, "calendar availability", func(callCtx context.Context) error {
		var err error
		externalBusy, err = s.calendar.ListBusyWindows(callCtx, s.cfg.CalendarID, windowStart.UTC(), windowEnd.UTC())
		return err
	})
	if err != nil {
		if errors.Is(err, google.ErrTokenNotFound) {
			return nil, errs.New(errs.CodeInternal, http.StatusServiceUnavailable, "calendar integration requires administrator authorization", err)
		}
		return nil, err
	}

	conflicts := detectConflicts(startLocal, startLocal.Add(duration), append(busyWindows, externalBusy...), buffer, loc)
	if conflicts {
		return nil, errs.New(errs.CodeConflict, http.StatusConflict, "requested slot conflicts with existing reservations", nil)
	}

	var calendarEvent *calendar.Event
	err = s.withRetry(ctx, s.calendarCB, "calendar booking", func(callCtx context.Context) error {
		input := calendar.EventInput{
			Summary:     s.reservationSummary(name),
			Description: buildEventDescription(name, email, agenda),
			Start:       startLocal,
			End:         startLocal.Add(duration),
			Attendees:   []string{email},
		}

		var err error
		calendarEvent, err = s.calendar.CreateEvent(callCtx, s.cfg.CalendarID, input)
		return err
	})
	if err != nil {
		if errors.Is(err, google.ErrTokenNotFound) {
			return nil, errs.New(errs.CodeInternal, http.StatusServiceUnavailable, "calendar integration requires administrator authorization", err)
		}
		return nil, err
	}

	meetURL := calendarEvent.HangoutLink
	if meetURL == "" {
		meetURL = calendarEvent.HTMLLink
	}

	newMeeting := model.Meeting{
		Name:            name,
		Email:           email,
		Datetime:        startLocal.UTC(),
		DurationMinutes: req.DurationMinutes,
		MeetURL:         strings.TrimSpace(meetURL),
		CalendarEventID: calendarEvent.ID,
		Status:          model.MeetingStatusPending,
		Notes:           agenda,
	}

	stored, err := s.meetings.CreateMeeting(ctx, &newMeeting)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to persist reservation", err)
	}

	mailErr := s.withRetry(ctx, s.mailCB, "notification email", func(callCtx context.Context) error {
		message := mail.Message{
			From:    s.cfg.NotificationSender,
			To:      []string{email},
			CC:      buildNotificationCC(s.cfg.NotificationReceiver),
			Subject: fmt.Sprintf("Meeting request confirmed: %s", stored.Datetime.In(loc).Format(time.RFC1123)),
			Body:    buildConfirmationBody(name, stored.Datetime.In(loc), duration, agenda, meetURL),
		}
		return s.mailer.Send(callCtx, message)
	})
	if mailErr != nil {
		if errors.Is(mailErr, google.ErrTokenNotFound) {
			return nil, errs.New(errs.CodeInternal, http.StatusServiceUnavailable, "email integration requires administrator authorization", mailErr)
		}
		return nil, mailErr
	}

	return &model.BookingResult{
		Meeting:         *stored,
		CalendarEventID: calendarEvent.ID,
	}, nil
}

func (s *bookingService) reservationSummary(name string) string {
	if template := strings.TrimSpace(s.cfg.MeetTemplate); template != "" {
		return fmt.Sprintf("%s - %s", template, name)
	}
	return fmt.Sprintf("Consultation with %s", name)
}

func (s *bookingService) withRetry(ctx context.Context, breaker *circuitBreaker, operation string, call func(ctx context.Context) error) error {
	attempts := s.cfg.MaxRetries
	if attempts < 1 {
		attempts = 1
	}
	backoff := s.cfg.InitialBackoff
	if backoff <= 0 {
		backoff = 750 * time.Millisecond
	}
	timeout := s.cfg.RequestTimeout
	if timeout <= 0 {
		timeout = 8 * time.Second
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if !breaker.Allow(s.clock.Now()) {
			return errs.New(errs.CodeInternal, http.StatusServiceUnavailable, fmt.Sprintf("%s temporarily unavailable", operation), lastErr)
		}

		callCtx, cancel := context.WithTimeout(ctx, timeout)
		err := call(callCtx)
		cancel()
		if err == nil {
			breaker.Success()
			return nil
		}
		if errors.Is(err, google.ErrTokenNotFound) {
			if breaker != nil {
				breaker.Success()
			}
			return err
		}

		lastErr = err
		breaker.Failure(s.clock.Now())
		if !isRetryable(err) || attempt == attempts {
			return errs.New(errs.CodeInternal, http.StatusBadGateway, fmt.Sprintf("failed to execute %s", operation), err)
		}

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return errs.New(errs.CodeInternal, http.StatusGatewayTimeout, fmt.Sprintf("%s aborted due to context cancellation", operation), ctx.Err())
		}
		backoff = time.Duration(float64(backoff) * s.cfg.BackoffMultiplier)
	}

	if lastErr == nil {
		lastErr = errors.New("unknown error")
	}
	return errs.New(errs.CodeInternal, http.StatusBadGateway, fmt.Sprintf("failed to execute %s", operation), lastErr)
}

func detectConflicts(start, end time.Time, busy []model.TimeWindow, buffer time.Duration, loc *time.Location) bool {
	if len(busy) == 0 {
		return false
	}

	windows := make([]model.TimeWindow, 0, len(busy))
	for _, window := range busy {
		localised := model.TimeWindow{
			Start: window.Start.In(loc).Add(-buffer),
			End:   window.End.In(loc).Add(buffer),
		}
		if localised.End.Before(localised.Start) {
			localised.End = localised.Start
		}
		windows = append(windows, localised)
	}

	sort.Slice(windows, func(i, j int) bool {
		return windows[i].Start.Before(windows[j].Start)
	})

	for _, window := range windows {
		if window.Start.Before(end) && window.End.After(start) {
			return true
		}
	}
	return false
}

func validateBookingRequest(req model.BookingRequest) *errs.AppError {
	if strings.TrimSpace(req.Name) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "name is required", nil)
	}
	if strings.TrimSpace(req.Email) == "" {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "email is required", nil)
	}
	if _, err := mailpkg.ParseAddress(strings.TrimSpace(req.Email)); err != nil {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "email format is invalid", err)
	}
	if req.StartTime.IsZero() {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "start time is required", nil)
	}
	if req.DurationMinutes <= 0 {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "durationMinutes must be greater than zero", nil)
	}
	if req.DurationMinutes > 240 {
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "durationMinutes exceeds maximum allowed duration", nil)
	}
	return nil
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	return true
}

type breakerState int

const (
	breakerClosed breakerState = iota
	breakerOpen
	breakerHalfOpen
)

type circuitBreaker struct {
	mu        sync.Mutex
	state     breakerState
	failures  int
	threshold int
	openFor   time.Duration
	openedAt  time.Time
}

func newCircuitBreaker(threshold int, openFor time.Duration) *circuitBreaker {
	if threshold <= 0 {
		threshold = 3
	}
	if openFor <= 0 {
		openFor = time.Minute
	}
	return &circuitBreaker{state: breakerClosed, threshold: threshold, openFor: openFor}
}

func (cb *circuitBreaker) Allow(now time.Time) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case breakerClosed:
		return true
	case breakerOpen:
		if now.Sub(cb.openedAt) >= cb.openFor {
			cb.state = breakerHalfOpen
			return true
		}
		return false
	case breakerHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *circuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = breakerClosed
	cb.failures = 0
}

func (cb *circuitBreaker) Failure(now time.Time) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	if cb.state == breakerHalfOpen || cb.failures >= cb.threshold {
		cb.state = breakerOpen
		cb.openedAt = now
		cb.failures = 0
	}
}

func buildNotificationCC(receiver string) []string {
	receiver = strings.TrimSpace(receiver)
	if receiver == "" {
		return nil
	}
	return []string{receiver}
}

func buildConfirmationBody(name string, start time.Time, duration time.Duration, agenda, meetURL string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Hi %s,\n\n", name))
	builder.WriteString(fmt.Sprintf("Your meeting has been scheduled for %s (duration: %.0f minutes).\n", start.Format(time.RFC1123), duration.Minutes()))
	if agenda != "" {
		builder.WriteString("\nAgenda:\n")
		builder.WriteString(agenda)
		builder.WriteString("\n")
	}
	if meetURL != "" {
		builder.WriteString(fmt.Sprintf("\nJoin via Google Meet: %s\n", meetURL))
	}
	builder.WriteString("\nThank you,\nPortfolio Site\n")
	return builder.String()
}

func buildEventDescription(name, email, agenda string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Meeting with %s (%s)\n", name, email))
	if agenda != "" {
		builder.WriteString("\nAgenda:\n")
		builder.WriteString(agenda)
		builder.WriteString("\n")
	}
	return builder.String()
}
