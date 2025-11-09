package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type meetingReservationRepository struct {
	db *sqlx.DB
}

// NewMeetingReservationRepository returns a MySQL-backed reservation repository.
func NewMeetingReservationRepository(db *sqlx.DB) repository.MeetingReservationRepository {
	return &meetingReservationRepository{db: db}
}

type reservationRow struct {
	ID                     uint64         `db:"id"`
	Name                   sql.NullString `db:"name"`
	Email                  sql.NullString `db:"email"`
	Topic                  sql.NullString `db:"topic"`
	Message                sql.NullString `db:"message"`
	StartAt                time.Time      `db:"start_at"`
	EndAt                  time.Time      `db:"end_at"`
	DurationMinutes        int            `db:"duration_minutes"`
	GoogleEventID          sql.NullString `db:"google_event_id"`
	GoogleCalendarStatus   sql.NullString `db:"google_calendar_status"`
	Status                 string         `db:"status"`
	ConfirmationSentAt     sql.NullTime   `db:"confirmation_sent_at"`
	LastNotificationSentAt sql.NullTime   `db:"last_notification_sent_at"`
	LookupHash             string         `db:"lookup_hash"`
	CancellationReason     sql.NullString `db:"cancellation_reason"`
	CreatedAt              time.Time      `db:"created_at"`
	UpdatedAt              time.Time      `db:"updated_at"`
}

type notificationRow struct {
	ID            uint64         `db:"id"`
	ReservationID uint64         `db:"reservation_id"`
	Type          string         `db:"notification_type"`
	Status        string         `db:"status"`
	ErrorMessage  sql.NullString `db:"error_message"`
	CreatedAt     time.Time      `db:"created_at"`
}

const insertReservationQuery = `
INSERT INTO meeting_reservations (
	name,
	email,
	topic,
	message,
	start_at,
	end_at,
	duration_minutes,
	google_event_id,
	google_calendar_status,
	status,
	confirmation_sent_at,
	last_notification_sent_at,
	lookup_hash,
	cancellation_reason,
	created_at,
	updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`

const selectByLookupQuery = `
SELECT
	id,
	name,
	email,
	topic,
	message,
	start_at,
	end_at,
	duration_minutes,
	google_event_id,
	google_calendar_status,
	status,
	confirmation_sent_at,
	last_notification_sent_at,
	lookup_hash,
	cancellation_reason,
	created_at,
	updated_at
FROM meeting_reservations
WHERE lookup_hash = ?
LIMIT 1`

const conflictReservationsQuery = `
SELECT
	id,
	name,
	email,
	topic,
	message,
	start_at,
	end_at,
	duration_minutes,
	google_event_id,
	google_calendar_status,
	status,
	confirmation_sent_at,
	last_notification_sent_at,
	lookup_hash,
	cancellation_reason,
	created_at,
	updated_at
FROM meeting_reservations
WHERE status IN ('pending','confirmed')
  AND start_at < ?
  AND end_at > ?
ORDER BY start_at ASC`

const listReservationsBaseQuery = `
SELECT
	id,
	name,
	email,
	topic,
	message,
	start_at,
	end_at,
	duration_minutes,
	google_event_id,
	google_calendar_status,
	status,
	confirmation_sent_at,
	last_notification_sent_at,
	lookup_hash,
	cancellation_reason,
	created_at,
	updated_at
FROM meeting_reservations`

const markConfirmationQuery = `
UPDATE meeting_reservations
SET
	status = CASE WHEN status = 'cancelled' THEN status ELSE 'confirmed' END,
	confirmation_sent_at = ?,
	last_notification_sent_at = ?,
	updated_at = NOW(3)
WHERE id = ?`

const cancelReservationQuery = `
UPDATE meeting_reservations
SET
	status = 'cancelled',
	cancellation_reason = ?,
	google_calendar_status = 'cancelled',
	updated_at = NOW(3)
WHERE id = ?`

const insertNotificationQuery = `
INSERT INTO meeting_notifications (
	reservation_id,
	notification_type,
	status,
	error_message,
	created_at
) VALUES (?, ?, ?, ?, NOW(3))`

const listNotificationsQuery = `
SELECT
	id,
	reservation_id,
	notification_type,
	status,
	error_message,
	created_at
FROM meeting_notifications
WHERE reservation_id = ?
ORDER BY created_at ASC, id ASC`

func (r *meetingReservationRepository) CreateReservation(ctx context.Context, reservation *model.MeetingReservation) (*model.MeetingReservation, error) {
	if reservation == nil {
		return nil, repository.ErrInvalidInput
	}

	_, err := r.db.ExecContext(ctx, insertReservationQuery,
		strings.TrimSpace(reservation.Name),
		strings.ToLower(strings.TrimSpace(reservation.Email)),
		strings.TrimSpace(reservation.Topic),
		strings.TrimSpace(reservation.Message),
		reservation.StartAt.UTC(),
		reservation.EndAt.UTC(),
		reservation.DurationMinutes,
		strings.TrimSpace(reservation.GoogleEventID),
		strings.TrimSpace(reservation.GoogleCalendarStatus),
		string(reservation.Status),
		sql.NullTime{Time: timePtrValue(reservation.ConfirmationSentAt), Valid: reservation.ConfirmationSentAt != nil},
		sql.NullTime{Time: timePtrValue(reservation.LastNotificationSentAt), Valid: reservation.LastNotificationSentAt != nil},
		strings.TrimSpace(reservation.LookupHash),
		strings.TrimSpace(reservation.CancellationReason),
	)
	if err != nil {
		return nil, fmt.Errorf("insert meeting_reservations: %w", err)
	}

	return r.FindReservationByLookupHash(ctx, reservation.LookupHash)
}

func (r *meetingReservationRepository) FindReservationByLookupHash(ctx context.Context, lookupHash string) (*model.MeetingReservation, error) {
	if strings.TrimSpace(lookupHash) == "" {
		return nil, repository.ErrInvalidInput
	}

	var row reservationRow
	if err := r.db.GetContext(ctx, &row, selectByLookupQuery, lookupHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("select meeting_reservations lookup=%s: %w", lookupHash, err)
	}

	reservation := mapReservationRow(row)
	return &reservation, nil
}

func (r *meetingReservationRepository) FindReservationByID(ctx context.Context, id uint64) (*model.MeetingReservation, error) {
	if id == 0 {
		return nil, repository.ErrInvalidInput
	}
	return r.findByID(ctx, id)
}

func (r *meetingReservationRepository) ListReservations(ctx context.Context, filter repository.MeetingReservationListFilter) ([]model.MeetingReservation, error) {
	query := listReservationsBaseQuery
	var conditions []string
	var args []any

	if len(filter.Status) > 0 {
		placeholders := make([]string, 0, len(filter.Status))
		for _, status := range filter.Status {
			value := strings.TrimSpace(string(status))
			if value == "" {
				continue
			}
			placeholders = append(placeholders, "?")
			args = append(args, value)
		}
		if len(placeholders) > 0 {
			conditions = append(conditions, "status IN ("+strings.Join(placeholders, ",")+")")
		}
	}

	if email := strings.ToLower(strings.TrimSpace(filter.Email)); email != "" {
		conditions = append(conditions, "LOWER(email) = ?")
		args = append(args, email)
	}

	if filter.Date != nil {
		day := filter.Date.UTC()
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
		end := start.Add(24 * time.Hour)
		conditions = append(conditions, "start_at >= ? AND start_at < ?")
		args = append(args, start, end)
	}

	if len(conditions) > 0 {
		query = query + "\nWHERE " + strings.Join(conditions, " AND ")
	}
	query += "\nORDER BY start_at DESC, id DESC"

	var rows []reservationRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("list meeting_reservations: %w", err)
	}

	results := make([]model.MeetingReservation, 0, len(rows))
	for _, row := range rows {
		results = append(results, mapReservationRow(row))
	}
	return results, nil
}

func (r *meetingReservationRepository) ListConflictingReservations(ctx context.Context, start, end time.Time) ([]model.MeetingReservation, error) {
	var rows []reservationRow
	if err := r.db.SelectContext(ctx, &rows, conflictReservationsQuery, end.UTC(), start.UTC()); err != nil {
		return nil, fmt.Errorf("select conflicting reservations: %w", err)
	}

	results := make([]model.MeetingReservation, 0, len(rows))
	for _, row := range rows {
		results = append(results, mapReservationRow(row))
	}
	return results, nil
}

func (r *meetingReservationRepository) MarkConfirmationSent(ctx context.Context, id uint64, sentAt time.Time) (*model.MeetingReservation, error) {
	if id == 0 {
		return nil, repository.ErrInvalidInput
	}

	_, err := r.db.ExecContext(ctx, markConfirmationQuery, sentAt.UTC(), sentAt.UTC(), id)
	if err != nil {
		return nil, fmt.Errorf("update meeting_reservations confirmation id=%d: %w", id, err)
	}

	return r.findByID(ctx, id)
}

func (r *meetingReservationRepository) CancelReservation(ctx context.Context, id uint64, reason string) (*model.MeetingReservation, error) {
	if id == 0 {
		return nil, repository.ErrInvalidInput
	}

	_, err := r.db.ExecContext(ctx, cancelReservationQuery, strings.TrimSpace(reason), id)
	if err != nil {
		return nil, fmt.Errorf("cancel meeting_reservations id=%d: %w", id, err)
	}

	return r.findByID(ctx, id)
}

func (r *meetingReservationRepository) UpdateReservationStatus(ctx context.Context, id uint64, status model.MeetingReservationStatus, reason string) (*model.MeetingReservation, error) {
	if id == 0 {
		return nil, repository.ErrInvalidInput
	}

	statusValue := strings.TrimSpace(string(status))
	if statusValue == "" {
		return nil, repository.ErrInvalidInput
	}

	var (
		query string
		args  []any
	)

	switch status {
	case model.MeetingReservationStatusCancelled:
		query = `
UPDATE meeting_reservations
SET
	status = 'cancelled',
	cancellation_reason = ?,
	google_calendar_status = 'cancelled',
	updated_at = NOW(3)
WHERE id = ?`
		args = []any{strings.TrimSpace(reason), id}
	case model.MeetingReservationStatusConfirmed:
		query = `
UPDATE meeting_reservations
SET
	status = 'confirmed',
	cancellation_reason = NULL,
	google_calendar_status = CASE
		WHEN google_calendar_status IS NULL OR google_calendar_status = '' THEN 'confirmed'
		ELSE google_calendar_status
	END,
	updated_at = NOW(3)
WHERE id = ?`
		args = []any{id}
	case model.MeetingReservationStatusPending:
		query = `
UPDATE meeting_reservations
SET
	status = 'pending',
	cancellation_reason = NULL,
	updated_at = NOW(3)
WHERE id = ?`
		args = []any{id}
	default:
		return nil, repository.ErrInvalidInput
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("update meeting_reservations status id=%d: %w", id, err)
	}

	return r.findByID(ctx, id)
}

func (r *meetingReservationRepository) findByID(ctx context.Context, id uint64) (*model.MeetingReservation, error) {
	const query = `
SELECT
	id,
	name,
	email,
	topic,
	message,
	start_at,
	end_at,
	duration_minutes,
	google_event_id,
	google_calendar_status,
	status,
	confirmation_sent_at,
	last_notification_sent_at,
	lookup_hash,
	cancellation_reason,
	created_at,
	updated_at
FROM meeting_reservations
WHERE id = ?
LIMIT 1`

	var row reservationRow
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("select meeting_reservations id=%d: %w", id, err)
	}

	reservation := mapReservationRow(row)
	return &reservation, nil
}

func mapReservationRow(row reservationRow) model.MeetingReservation {
	var confirmationSentAt *time.Time
	if row.ConfirmationSentAt.Valid {
		ts := row.ConfirmationSentAt.Time.UTC()
		confirmationSentAt = &ts
	}

	var lastNotificationSentAt *time.Time
	if row.LastNotificationSentAt.Valid {
		ts := row.LastNotificationSentAt.Time.UTC()
		lastNotificationSentAt = &ts
	}

	return model.MeetingReservation{
		ID:                     row.ID,
		LookupHash:             strings.TrimSpace(row.LookupHash),
		Name:                   strings.TrimSpace(row.Name.String),
		Email:                  strings.TrimSpace(row.Email.String),
		Topic:                  strings.TrimSpace(row.Topic.String),
		Message:                strings.TrimSpace(row.Message.String),
		StartAt:                row.StartAt.UTC(),
		EndAt:                  row.EndAt.UTC(),
		DurationMinutes:        row.DurationMinutes,
		GoogleEventID:          strings.TrimSpace(row.GoogleEventID.String),
		GoogleCalendarStatus:   strings.TrimSpace(row.GoogleCalendarStatus.String),
		Status:                 model.MeetingReservationStatus(strings.TrimSpace(row.Status)),
		ConfirmationSentAt:     confirmationSentAt,
		LastNotificationSentAt: lastNotificationSentAt,
		CancellationReason:     strings.TrimSpace(row.CancellationReason.String),
		CreatedAt:              row.CreatedAt.UTC(),
		UpdatedAt:              row.UpdatedAt.UTC(),
	}
}

func timePtrValue(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.UTC()
}

type meetingNotificationRepository struct {
	db *sqlx.DB
}

// NewMeetingNotificationRepository returns a MySQL-backed notification repository.
func NewMeetingNotificationRepository(db *sqlx.DB) repository.MeetingNotificationRepository {
	return &meetingNotificationRepository{db: db}
}

func (r *meetingNotificationRepository) RecordNotification(ctx context.Context, notification *model.MeetingNotification) (*model.MeetingNotification, error) {
	if notification == nil {
		return nil, repository.ErrInvalidInput
	}

	_, err := r.db.ExecContext(ctx, insertNotificationQuery,
		notification.ReservationID,
		strings.TrimSpace(notification.Type),
		strings.TrimSpace(notification.Status),
		nullIfEmpty(notification.ErrorMessage),
	)
	if err != nil {
		return nil, fmt.Errorf("insert meeting_notifications: %w", err)
	}

	rows, err := r.ListNotifications(ctx, notification.ReservationID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, repository.ErrNotFound
	}
	return &rows[len(rows)-1], nil
}

func (r *meetingNotificationRepository) ListNotifications(ctx context.Context, reservationID uint64) ([]model.MeetingNotification, error) {
	var rows []notificationRow
	if err := r.db.SelectContext(ctx, &rows, listNotificationsQuery, reservationID); err != nil {
		return nil, fmt.Errorf("select meeting_notifications reservation_id=%d: %w", reservationID, err)
	}

	result := make([]model.MeetingNotification, 0, len(rows))
	for _, row := range rows {
		result = append(result, model.MeetingNotification{
			ID:            row.ID,
			ReservationID: row.ReservationID,
			Type:          strings.TrimSpace(row.Type),
			Status:        strings.TrimSpace(row.Status),
			ErrorMessage:  strings.TrimSpace(row.ErrorMessage.String),
			CreatedAt:     row.CreatedAt.UTC(),
		})
	}
	return result, nil
}

func nullIfEmpty(value string) sql.NullString {
	if strings.TrimSpace(value) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.TrimSpace(value), Valid: true}
}

var (
	_ repository.MeetingReservationRepository  = (*meetingReservationRepository)(nil)
	_ repository.MeetingNotificationRepository = (*meetingNotificationRepository)(nil)
)
