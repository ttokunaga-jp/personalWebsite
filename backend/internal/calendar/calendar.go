package calendar

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
)

// EventInput describes an event to insert into Google Calendar.
type EventInput struct {
	Summary     string
	Description string
	Start       time.Time
	End         time.Time
	Attendees   []string
}

// Event captures the important identifiers of a Google Calendar entry.
type Event struct {
	ID          string
	HTMLLink    string
	HangoutLink string
}

// Client abstracts the subset of Google Calendar operations required for booking.
type Client interface {
	ListBusyWindows(ctx context.Context, calendarID string, from, to time.Time) ([]model.TimeWindow, error)
	CreateEvent(ctx context.Context, calendarID string, input EventInput) (*Event, error)
}
