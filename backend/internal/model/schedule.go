package model

import "time"

// BusyWindowSource identifies the origin of a busy window.
type BusyWindowSource string

const (
	BusyWindowSourceReservation BusyWindowSource = "reservation"
	BusyWindowSourceExternal    BusyWindowSource = "external"
	BusyWindowSourceBlackout    BusyWindowSource = "blackout"
)

// TimeWindow represents an inclusive-exclusive time span.
type TimeWindow struct {
	Start  time.Time
	End    time.Time
	Source BusyWindowSource
}

// AvailabilitySlotStatus identifies the classification of a rendered slot.
type AvailabilitySlotStatus string

const (
	AvailabilitySlotStatusAvailable AvailabilitySlotStatus = "available"
	AvailabilitySlotStatusReserved  AvailabilitySlotStatus = "reserved"
	AvailabilitySlotStatusBlackout  AvailabilitySlotStatus = "blackout"
)

// AvailabilitySlot describes a single bookable slot.
type AvailabilitySlot struct {
	ID         string                 `json:"id"`
	Start      time.Time              `json:"start"`
	End        time.Time              `json:"end"`
	Status     AvailabilitySlotStatus `json:"status"`
	IsBookable bool                   `json:"isBookable"`
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
