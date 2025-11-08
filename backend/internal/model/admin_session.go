package model

import "time"

// AdminSession represents a server-side persisted administrator session.
type AdminSession struct {
	ID             string
	TokenHash      string
	Subject        string
	Email          string
	Roles          []string
	UserAgent      string
	IPAddress      string
	ExpiresAt      time.Time
	LastAccessedAt time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RevokedAt      *time.Time
}
