package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type availabilityRepository struct {
	db *sqlx.DB
}

// NewAvailabilityRepository creates a repository for busy window queries.
func NewAvailabilityRepository(db *sqlx.DB) repository.AvailabilityRepository {
	return &availabilityRepository{db: db}
}

const (
	reservationsBusyQuery = `
SELECT
	start_at AS start_time,
	end_at AS end_time,
	'reservation' AS source
FROM meeting_reservations
WHERE status IN ('pending', 'confirmed')
  AND end_at > ?
  AND start_at < ?
ORDER BY start_at`

	blackoutBusyQuery = `
SELECT
	start_time,
	end_time,
	'blackout' AS source
FROM schedule_blackouts
WHERE end_time > ?
  AND start_time < ?
ORDER BY start_time`
)

type busyRow struct {
	Start  time.Time `db:"start_time"`
	End    time.Time `db:"end_time"`
	Source string    `db:"source"`
}

func (r *availabilityRepository) ListBusyWindows(ctx context.Context, from, to time.Time) ([]model.TimeWindow, error) {
	windows := make([]model.TimeWindow, 0, 16)

	var meetingRows []busyRow
	if err := r.db.SelectContext(ctx, &meetingRows, reservationsBusyQuery, from, to); err != nil {
		return nil, fmt.Errorf("select reservation busy windows: %w", err)
	}

	for _, row := range meetingRows {
		windows = append(windows, model.TimeWindow{
			Start:  row.Start.UTC(),
			End:    row.End.UTC(),
			Source: model.BusyWindowSource(row.Source),
		})
	}

	var blackoutRows []busyRow
	if err := r.db.SelectContext(ctx, &blackoutRows, blackoutBusyQuery, from, to); err != nil {
		return nil, fmt.Errorf("select blackout windows: %w", err)
	}

	for _, row := range blackoutRows {
		windows = append(windows, model.TimeWindow{
			Start:  row.Start.UTC(),
			End:    row.End.UTC(),
			Source: model.BusyWindowSource(row.Source),
		})
	}

	return windows, nil
}
