package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type meetingRepository struct {
	db *sqlx.DB
}

// NewMeetingRepository returns a MySQL-backed meeting repository.
func NewMeetingRepository(db *sqlx.DB) repository.MeetingRepository {
	return &meetingRepository{db: db}
}

const listMeetingsQuery = `
SELECT
	m.id,
	m.name,
	m.email,
	m.meeting_at,
	m.duration_minutes,
	m.meet_url,
	m.calendar_event_id,
	m.status,
	m.notes,
	m.created_at,
	m.updated_at
FROM meetings m
ORDER BY m.meeting_at DESC, m.id DESC`

const getMeetingQuery = `
SELECT
	m.id,
	m.name,
	m.email,
	m.meeting_at,
	m.duration_minutes,
	m.meet_url,
	m.calendar_event_id,
	m.status,
	m.notes,
	m.created_at,
	m.updated_at
FROM meetings m
WHERE m.id = ?`

const insertMeetingQuery = `
INSERT INTO meetings (
	name,
	email,
	meeting_at,
	duration_minutes,
	meet_url,
	calendar_event_id,
	status,
	notes,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

const updateMeetingQuery = `
UPDATE meetings
SET
	name = ?,
	email = ?,
	meeting_at = ?,
	duration_minutes = ?,
	meet_url = ?,
	calendar_event_id = ?,
	status = ?,
	notes = ?,
	updated_at = NOW()
WHERE id = ?`

const deleteMeetingQuery = `DELETE FROM meetings WHERE id = ?`

type meetingRow struct {
	ID              int64          `db:"id"`
	Name            sql.NullString `db:"name"`
	Email           sql.NullString `db:"email"`
	MeetingAt       sql.NullTime   `db:"meeting_at"`
	DurationMinutes sql.NullInt64  `db:"duration_minutes"`
	MeetURL         sql.NullString `db:"meet_url"`
	CalendarEventID sql.NullString `db:"calendar_event_id"`
	Status          sql.NullString `db:"status"`
	Notes           sql.NullString `db:"notes"`
	CreatedAt       sql.NullTime   `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
}

func (r *meetingRepository) ListMeetings(ctx context.Context) ([]model.Meeting, error) {
	var rows []meetingRow
	if err := r.db.SelectContext(ctx, &rows, listMeetingsQuery); err != nil {
		return nil, fmt.Errorf("select meetings: %w", err)
	}

	meetings := make([]model.Meeting, 0, len(rows))
	for _, row := range rows {
		meetings = append(meetings, mapMeetingRow(row))
	}
	return meetings, nil
}

func (r *meetingRepository) GetMeeting(ctx context.Context, id int64) (*model.Meeting, error) {
	var row meetingRow
	if err := r.db.GetContext(ctx, &row, getMeetingQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get meeting %d: %w", id, err)
	}

	meeting := mapMeetingRow(row)
	return &meeting, nil
}

func (r *meetingRepository) CreateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error) {
	if meeting == nil {
		return nil, repository.ErrInvalidInput
	}

	res, err := r.db.ExecContext(ctx, insertMeetingQuery,
		meeting.Name,
		meeting.Email,
		meeting.Datetime,
		meeting.DurationMinutes,
		meeting.MeetURL,
		meeting.CalendarEventID,
		meeting.Status,
		meeting.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("insert meeting: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("meeting last insert id: %w", err)
	}

	created, err := r.GetMeeting(ctx, id)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *meetingRepository) UpdateMeeting(ctx context.Context, meeting *model.Meeting) (*model.Meeting, error) {
	if meeting == nil {
		return nil, repository.ErrInvalidInput
	}

	res, err := r.db.ExecContext(ctx, updateMeetingQuery,
		meeting.Name,
		meeting.Email,
		meeting.Datetime,
		meeting.DurationMinutes,
		meeting.MeetURL,
		meeting.CalendarEventID,
		meeting.Status,
		meeting.Notes,
		meeting.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update meeting %d: %w", meeting.ID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected meeting %d: %w", meeting.ID, err)
	}
	if affected == 0 {
		return nil, repository.ErrNotFound
	}

	updated, err := r.GetMeeting(ctx, meeting.ID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *meetingRepository) DeleteMeeting(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, deleteMeetingQuery, id)
	if err != nil {
		return fmt.Errorf("delete meeting %d: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected delete meeting %d: %w", id, err)
	}
	if affected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func mapMeetingRow(row meetingRow) model.Meeting {
	duration := 0
	if row.DurationMinutes.Valid {
		duration = int(row.DurationMinutes.Int64)
	}
	meetingTime := row.MeetingAt.Time
	if !row.MeetingAt.Valid {
		meetingTime = timeNowUTC()
	}
	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = meetingTime
	}
	updatedAt := row.UpdatedAt.Time
	if !row.UpdatedAt.Valid {
		updatedAt = createdAt
	}

	return model.Meeting{
		ID:              row.ID,
		Name:            row.Name.String,
		Email:           row.Email.String,
		Datetime:        meetingTime,
		DurationMinutes: duration,
		MeetURL:         row.MeetURL.String,
		CalendarEventID: row.CalendarEventID.String,
		Status:          model.MeetingStatus(row.Status.String),
		Notes:           row.Notes.String,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func timeNowUTC() time.Time {
	return time.Now().UTC()
}

var _ repository.MeetingRepository = (*meetingRepository)(nil)
