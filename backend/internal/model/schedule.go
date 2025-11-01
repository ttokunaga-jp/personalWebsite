package model

import "time"

// TimeWindow represents an inclusive-exclusive time span.
type TimeWindow struct {
	Start time.Time
	End   time.Time
}

// AvailabilitySlot describes a single bookable slot.
type AvailabilitySlot struct {
	ID         string    `json:"id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	IsBookable bool      `json:"isBookable"`
}

// AvailabilityDay aggregates slots for a calendar date.
type AvailabilityDay struct {
	Date  string             `json:"date"`
	Slots []AvailabilitySlot `json:"slots"`
}

// AvailabilityResponse is returned by the contact availability endpoint.
type AvailabilityResponse struct {
	Timezone    string            `json:"timezone"`
	GeneratedAt time.Time         `json:"generatedAt"`
	Days        []AvailabilityDay `json:"days"`
}
