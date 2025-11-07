package firestore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type availabilityRepository struct {
	base baseRepository
}

const blackoutsCollection = "schedule_blackouts"

type blackoutDocument struct {
	StartTime time.Time `firestore:"startTime"`
	EndTime   time.Time `firestore:"endTime"`
}

func NewAvailabilityRepository(client *firestore.Client, prefix string) repository.AvailabilityRepository {
	return &availabilityRepository{base: newBaseRepository(client, prefix)}
}

func (r *availabilityRepository) ListBusyWindows(ctx context.Context, from, to time.Time) ([]model.TimeWindow, error) {
	var windows []model.TimeWindow

	meetingWindows, err := r.fetchMeetingWindows(ctx, from, to)
	if err != nil {
		return nil, err
	}
	windows = append(windows, meetingWindows...)

	blackoutWindows, err := r.fetchBlackoutWindows(ctx, from, to)
	if err != nil {
		return nil, err
	}
	windows = append(windows, blackoutWindows...)

	return windows, nil
}

func (r *availabilityRepository) fetchMeetingWindows(ctx context.Context, from, to time.Time) ([]model.TimeWindow, error) {
	// Firestore-backed meeting reservations are not implemented; return no busy windows.
	return []model.TimeWindow{}, nil
}

func (r *availabilityRepository) fetchBlackoutWindows(ctx context.Context, from, to time.Time) ([]model.TimeWindow, error) {
	query := r.base.collection(blackoutsCollection).
		Where("endTime", ">", from).
		Where("startTime", "<", to)

	snapshots, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore availability: list blackouts: %w", err)
	}

	windows := make([]model.TimeWindow, 0, len(snapshots))
	for _, snap := range snapshots {
		var entry blackoutDocument
		if err := snap.DataTo(&entry); err != nil {
			return nil, fmt.Errorf("firestore availability: decode blackout %s: %w", snap.Ref.ID, err)
		}
		windows = append(windows, model.TimeWindow{
			Start:  entry.StartTime,
			End:    entry.EndTime,
			Source: model.BusyWindowSourceBlackout,
		})
	}
	return windows, nil
}

var _ repository.AvailabilityRepository = (*availabilityRepository)(nil)
