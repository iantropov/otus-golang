package sqlstorage

import "time"

type Event struct {
	ID           int64         `db:"id"`
	Title        string        `db:"title"`
	StartsAt     time.Time     `db:"starts_at"`
	EndsAt       time.Time     `db:"ends_at"`
	Description  string        `db:"description"`
	UserID       string        `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}
