package storage

import "time"

type EventId string

type Event struct {
	Id           EventId
	Title        string
	StartsAt     time.Time
	EndsAt       time.Time
	Description  string
	UserId       string
	NotifyBefore time.Duration
}
