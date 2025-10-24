package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/takumi/personal-website/internal/calendar"
	"github.com/takumi/personal-website/internal/model"
)

const (
	calendarBaseURL   = "https://www.googleapis.com/calendar/v3"
	freeBusyEndpoint  = "/freeBusy"
	calendarEventsFmt = "/calendars/%s/events"
)

// CalendarAPIClient implements Google Calendar REST calls.
type CalendarAPIClient struct {
	client        *http.Client
	tokenProvider TokenProvider
	timezone      string
}

// NewCalendarAPIClient constructs a Google Calendar client.
func NewCalendarAPIClient(httpClient *http.Client, tokenProvider TokenProvider, timezone string) *CalendarAPIClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &CalendarAPIClient{
		client:        httpClient,
		tokenProvider: tokenProvider,
		timezone:      timezone,
	}
}

func (c *CalendarAPIClient) ListBusyWindows(ctx context.Context, calendarID string, from, to time.Time) ([]model.TimeWindow, error) {
	token, err := c.tokenProvider.AccessToken(ctx)
	// If token provider is not configured, allow graceful degradation.
	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"timeMin": from.Format(time.RFC3339),
		"timeMax": to.Format(time.RFC3339),
		"items": []map[string]string{
			{"id": calendarID},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("calendar freebusy marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, calendarBaseURL+freeBusyEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("calendar freebusy request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calendar freebusy call: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("calendar freebusy error: status=%d body=%s", resp.StatusCode, string(payload))
	}

	var decoded struct {
		Calendars map[string]struct {
			Busy []struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"busy"`
		} `json:"calendars"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("calendar freebusy decode: %w", err)
	}

	calendarEntry, ok := decoded.Calendars[calendarID]
	if !ok {
		return nil, nil
	}

	windows := make([]model.TimeWindow, 0, len(calendarEntry.Busy))
	for _, busy := range calendarEntry.Busy {
		start, err := time.Parse(time.RFC3339, busy.Start)
		if err != nil {
			continue
		}
		end, err := time.Parse(time.RFC3339, busy.End)
		if err != nil {
			continue
		}
		windows = append(windows, model.TimeWindow{Start: start.UTC(), End: end.UTC()})
	}

	return windows, nil
}

func (c *CalendarAPIClient) CreateEvent(ctx context.Context, calendarID string, input calendar.EventInput) (*calendar.Event, error) {
	token, err := c.tokenProvider.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	event := map[string]any{
		"summary":     input.Summary,
		"description": input.Description,
		"start": map[string]string{
			"dateTime": input.Start.Format(time.RFC3339),
			"timeZone": c.timezone,
		},
		"end": map[string]string{
			"dateTime": input.End.Format(time.RFC3339),
			"timeZone": c.timezone,
		},
		"attendees": buildAttendees(input.Attendees),
		"conferenceData": map[string]any{
			"createRequest": map[string]any{
				"requestId": fmt.Sprintf("booking-%d", time.Now().UnixNano()),
				"conferenceSolutionKey": map[string]string{
					"type": "hangoutsMeet",
				},
			},
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("calendar insert marshal: %w", err)
	}

	path := fmt.Sprintf(calendarEventsFmt, url.PathEscape(calendarID))
	endpoint := fmt.Sprintf("%s%s?conferenceDataVersion=1", calendarBaseURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("calendar insert request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calendar insert call: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("calendar insert error: status=%d body=%s", resp.StatusCode, string(payload))
	}

	var decoded struct {
		ID          string `json:"id"`
		HTMLLink    string `json:"htmlLink"`
		HangoutLink string `json:"hangoutLink"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("calendar insert decode: %w", err)
	}

	return &calendar.Event{
		ID:          decoded.ID,
		HTMLLink:    decoded.HTMLLink,
		HangoutLink: decoded.HangoutLink,
	}, nil
}

func buildAttendees(addresses []string) []map[string]string {
	list := make([]map[string]string, 0, len(addresses))
	for _, addr := range addresses {
		if addr == "" {
			continue
		}
		list = append(list, map[string]string{"email": addr})
	}
	return list
}
