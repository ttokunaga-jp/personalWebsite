package model

import "time"

// MeetingReservationStatus captures the lifecycle of a reservation.
type MeetingReservationStatus string

const (
	MeetingReservationStatusPending   MeetingReservationStatus = "pending"
	MeetingReservationStatusConfirmed MeetingReservationStatus = "confirmed"
	MeetingReservationStatusCancelled MeetingReservationStatus = "cancelled"
)

// MeetingReservation represents a booking persisted in meeting_reservations.
type MeetingReservation struct {
	ID                     uint64                   `json:"id"`
	LookupHash             string                   `json:"lookupHash"`
	Name                   string                   `json:"name"`
	Email                  string                   `json:"email"`
	Topic                  string                   `json:"topic"`
	Message                string                   `json:"message"`
	StartAt                time.Time                `json:"startAt"`
	EndAt                  time.Time                `json:"endAt"`
	DurationMinutes        int                      `json:"durationMinutes"`
	GoogleEventID          string                   `json:"googleEventId"`
	GoogleCalendarStatus   string                   `json:"googleCalendarStatus"`
	Status                 MeetingReservationStatus `json:"status"`
	ConfirmationSentAt     *time.Time               `json:"confirmationSentAt,omitempty"`
	LastNotificationSentAt *time.Time               `json:"lastNotificationSentAt,omitempty"`
	CancellationReason     string                   `json:"cancellationReason,omitempty"`
	CreatedAt              time.Time                `json:"createdAt"`
	UpdatedAt              time.Time                `json:"updatedAt"`
}

// MeetingNotification captures entries written to meeting_notifications.
type MeetingNotification struct {
	ID            uint64    `json:"id"`
	ReservationID uint64    `json:"reservationId"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	ErrorMessage  string    `json:"errorMessage,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}
