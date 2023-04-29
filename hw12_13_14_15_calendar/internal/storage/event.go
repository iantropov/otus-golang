package storage

import "time"

type EventID string

type Event struct {
	ID           EventID
	Title        string
	StartsAt     time.Time
	EndsAt       time.Time
	Description  string
	UserID       string
	NotifyBefore time.Duration
}
