package service

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	"github.com/takumi/personal-website/internal/service/support"
)

// AvailabilityOptions allows handlers to customise the scheduling horizon.
type AvailabilityOptions struct {
	StartDate time.Time
	Days      int
}

// AvailabilityService calculates bookable contact slots.
type AvailabilityService interface {
	GetAvailability(ctx context.Context, opts AvailabilityOptions) (*model.AvailabilityResponse, error)
}

type availabilityService struct {
	repo         repository.AvailabilityRepository
	fallbackRepo repository.AvailabilityRepository
	cfg          config.ContactConfig
}

// NewAvailabilityService wires availability logic to the repository and configuration.
func NewAvailabilityService(repo repository.AvailabilityRepository, cfg *config.AppConfig) AvailabilityService {
	return &availabilityService{
		repo:         repo,
		fallbackRepo: inmemory.NewAvailabilityRepository(),
		cfg:          cfg.Contact,
	}
}

func (s *availabilityService) GetAvailability(ctx context.Context, opts AvailabilityOptions) (*model.AvailabilityResponse, error) {
	if s.repo == nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "availability repository not configured", nil)
	}

	loc, err := time.LoadLocation(s.cfg.Timezone)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "invalid contact timezone configuration", err)
	}

	slotDuration := time.Duration(s.cfg.SlotDurationMin) * time.Minute
	if slotDuration <= 0 {
		slotDuration = 30 * time.Minute
	}

	if s.cfg.WorkdayEndHour <= s.cfg.WorkdayStartHour {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "invalid contact workday configuration", nil)
	}

	buffer := time.Duration(s.cfg.BufferMinutes) * time.Minute
	if buffer < 0 {
		buffer = 0
	}

	horizon := opts.Days
	if horizon <= 0 {
		horizon = s.cfg.HorizonDays
	}
	if horizon <= 0 {
		horizon = 14
	}

	startDate := opts.StartDate
	if startDate.IsZero() {
		startDate = time.Now().In(loc)
	} else {
		startDate = startDate.In(loc)
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	endDate := startDate.AddDate(0, 0, horizon)

	busyWindows, err := s.repo.ListBusyWindows(ctx, startDate.UTC(), endDate.UTC())
	if err != nil {
		if s.fallbackRepo != nil && support.ShouldFallback(err) {
			if fallbackWindows, fallbackErr := s.fallbackRepo.ListBusyWindows(ctx, startDate.UTC(), endDate.UTC()); fallbackErr == nil {
				busyWindows = fallbackWindows
				err = nil
			}
		}
		if err != nil {
			return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load busy windows", err)
		}
	}

	expanded := expandAndMergeWindows(busyWindows, buffer, loc)

	days := make([]model.AvailabilityDay, 0, horizon)
	for day := 0; day < horizon; day++ {
		current := startDate.AddDate(0, 0, day)
		dayStart := time.Date(current.Year(), current.Month(), current.Day(), s.cfg.WorkdayStartHour, 0, 0, 0, loc)
		dayEnd := time.Date(current.Year(), current.Month(), current.Day(), s.cfg.WorkdayEndHour, 0, 0, 0, loc)

		slots := buildSlots(dayStart, dayEnd, slotDuration, expanded)
		days = append(days, model.AvailabilityDay{
			Date:  dayStart.Format("2006-01-02"),
			Slots: slots,
		})
	}

	return &model.AvailabilityResponse{
		Timezone:    s.cfg.Timezone,
		GeneratedAt: time.Now().In(loc),
		Days:        days,
	}, nil
}

func buildSlots(dayStart, dayEnd time.Time, slotDuration time.Duration, busy []model.TimeWindow) []model.AvailabilitySlot {
	slots := make([]model.AvailabilitySlot, 0, 16)
	for cursor := dayStart; !cursor.Add(slotDuration).After(dayEnd); cursor = cursor.Add(slotDuration) {
		slotEnd := cursor.Add(slotDuration)
		status := model.AvailabilitySlotStatusAvailable
		isBookable := true
		for _, window := range busy {
			if window.Start.Before(slotEnd) && window.End.After(cursor) {
				switch window.Source {
				case model.BusyWindowSourceBlackout:
					status = model.AvailabilitySlotStatusBlackout
				default:
					status = model.AvailabilitySlotStatusReserved
				}
				isBookable = false
				if status == model.AvailabilitySlotStatusBlackout {
					break
				}
			}
		}
		slots = append(slots, model.AvailabilitySlot{
			ID:         cursor.UTC().Format(time.RFC3339),
			Start:      cursor,
			End:        slotEnd,
			Status:     status,
			IsBookable: isBookable,
		})
	}
	return slots
}

func expandAndMergeWindows(windows []model.TimeWindow, buffer time.Duration, loc *time.Location) []model.TimeWindow {
	if len(windows) == 0 {
		return nil
	}

	expanded := make([]model.TimeWindow, 0, len(windows))
	for _, window := range windows {
		start := window.Start.In(loc).Add(-buffer)
		end := window.End.In(loc).Add(buffer)
		if end.Before(start) {
			end = start
		}
		expanded = append(expanded, model.TimeWindow{
			Start:  start,
			End:    end,
			Source: window.Source,
		})
	}

	sort.Slice(expanded, func(i, j int) bool {
		return expanded[i].Start.Before(expanded[j].Start)
	})

	merged := []model.TimeWindow{expanded[0]}
	for i := 1; i < len(expanded); i++ {
		last := merged[len(merged)-1]
		current := expanded[i]
		if current.Start.Before(last.End) || current.Start.Equal(last.End) {
			if current.End.After(last.End) {
				last.End = current.End
			}
			last.Source = dominantSource(last.Source, current.Source)
			merged[len(merged)-1] = last
			continue
		}
		merged = append(merged, current)
	}

	return merged
}

func dominantSource(a, b model.BusyWindowSource) model.BusyWindowSource {
	if a == model.BusyWindowSourceBlackout || b == model.BusyWindowSourceBlackout {
		return model.BusyWindowSourceBlackout
	}
	if a == model.BusyWindowSourceReservation || b == model.BusyWindowSourceReservation {
		return model.BusyWindowSourceReservation
	}
	if a == model.BusyWindowSourceExternal || b == model.BusyWindowSourceExternal {
		return model.BusyWindowSourceExternal
	}
	return a
}
