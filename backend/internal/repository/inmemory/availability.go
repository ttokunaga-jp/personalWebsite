package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type availabilityRepository struct {
	busy []model.TimeWindow
}

// NewAvailabilityRepository constructs an in-memory availability repository with no busy windows.
func NewAvailabilityRepository() repository.AvailabilityRepository {
	return &availabilityRepository{
		busy: []model.TimeWindow{},
	}
}

// NewAvailabilityRepositoryWithWindows seeds the repository with predetermined busy windows.
func NewAvailabilityRepositoryWithWindows(windows []model.TimeWindow) repository.AvailabilityRepository {
	return &availabilityRepository{
		busy: append([]model.TimeWindow(nil), windows...),
	}
}

func (r *availabilityRepository) ListBusyWindows(ctx context.Context, from, to time.Time) ([]model.TimeWindow, error) {
	result := make([]model.TimeWindow, 0, len(r.busy))
	for _, window := range r.busy {
		if window.End.Before(from) || window.Start.After(to) {
			continue
		}
		result = append(result, window)
	}
	return result, nil
}
