package model

import "time"

// BookingRequest represents an inbound reservation submitted from the contact form.
type BookingRequest struct {
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	StartTime       time.Time `json:"startTime"`
	DurationMinutes int       `json:"durationMinutes"`
	Agenda          string    `json:"agenda"`
	Topic           string    `json:"topic"`
	RecaptchaToken  string    `json:"recaptchaToken"`
}

// BookingResult summarises a booked meeting reservation and associated metadata.
type BookingResult struct {
	Reservation      MeetingReservation `json:"reservation"`
	CalendarEventID  string             `json:"calendarEventId"`
	SupportEmail     string             `json:"supportEmail"`
	CalendarTimezone string             `json:"calendarTimezone"`
}
