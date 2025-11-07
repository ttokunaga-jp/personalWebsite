package service

import (
	"context"
	"testing"
	"time"

	mysqlerr "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/model"
)

func TestAvailabilityService_RespectsBuffer(t *testing.T) {
	t.Parallel()

	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	busyStart := time.Date(2024, time.March, 10, 13, 0, 0, 0, loc)
	busyEnd := busyStart.Add(time.Hour)

	repo := &stubAvailabilityRepo{
		windows: []model.TimeWindow{
			{Start: busyStart, End: busyEnd},
		},
	}

	svc := NewAvailabilityService(repo, &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:         "Asia/Tokyo",
			SlotDurationMin:  30,
			WorkdayStartHour: 9,
			WorkdayEndHour:   18,
			HorizonDays:      1,
			BufferMinutes:    30,
		},
	})

	resp, err := svc.GetAvailability(context.Background(), AvailabilityOptions{StartDate: busyStart})
	require.NoError(t, err)
	require.Len(t, resp.Days, 1)

	slots := resp.Days[0].Slots
	require.NotZero(t, len(slots))

	forbidden := map[string]struct{}{
		"12:30": {},
		"13:00": {},
		"13:30": {},
		"14:00": {},
	}

	for _, slot := range slots {
		key := slot.Start.In(loc).Format("15:04")
		_, blocked := forbidden[key]
		require.Falsef(t, blocked, "slot %s should be blocked due to buffer", key)
	}
}

func TestAvailabilityService_HonoursDaysOverride(t *testing.T) {
	t.Parallel()

	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	repo := &stubAvailabilityRepo{}

	svc := NewAvailabilityService(repo, &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:         "Asia/Tokyo",
			SlotDurationMin:  60,
			WorkdayStartHour: 10,
			WorkdayEndHour:   15,
			HorizonDays:      1,
			BufferMinutes:    0,
		},
	})

	resp, err := svc.GetAvailability(context.Background(), AvailabilityOptions{
		StartDate: time.Date(2024, time.November, 1, 0, 0, 0, 0, loc),
		Days:      3,
	})
	require.NoError(t, err)
	require.Len(t, resp.Days, 3)

	firstDay := resp.Days[0]
	require.Equal(t, "2024-11-01", firstDay.Date)
	require.Len(t, firstDay.Slots, 5)
}

type stubAvailabilityRepo struct {
	windows []model.TimeWindow
	err     error
}

func (s *stubAvailabilityRepo) ListBusyWindows(context.Context, time.Time, time.Time) ([]model.TimeWindow, error) {
	if s.err != nil {
		return nil, s.err
	}
	return append([]model.TimeWindow(nil), s.windows...), nil
}

func TestAvailabilityService_FallbackWhenRepositoryUnavailable(t *testing.T) {
	t.Parallel()

	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	repo := &stubAvailabilityRepo{
		err: &mysqlerr.MySQLError{
			Number:  1146,
			Message: "Table 'personal_website.meetings' doesn't exist",
		},
	}

	svc := NewAvailabilityService(repo, &config.AppConfig{
		Contact: config.ContactConfig{
			Timezone:         "Asia/Tokyo",
			SlotDurationMin:  30,
			WorkdayStartHour: 9,
			WorkdayEndHour:   12,
			HorizonDays:      1,
			BufferMinutes:    0,
		},
	})

	resp, reqErr := svc.GetAvailability(context.Background(), AvailabilityOptions{
		StartDate: time.Date(2024, time.June, 1, 0, 0, 0, 0, loc),
		Days:      1,
	})

	require.NoError(t, reqErr)
	require.NotNil(t, resp)
	require.Len(t, resp.Days, 1)
	require.NotEmpty(t, resp.Days[0].Slots)
}
