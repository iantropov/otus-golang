package sqlstorage

import "fmt"

var (
	eventAttributes       = `title, starts_at, ends_at, created_at, description, user_id, notify_before`
	eventAttributesWithID = `id,` + eventAttributes
	SelectEventByID       = fmt.Sprintf(
		`SELECT %s FROM events WHERE id=$1`,
		eventAttributesWithID,
	)
	InsertEvent = fmt.Sprintf(
		`INSERT INTO events(%s) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		eventAttributesWithID,
	)
	UpdateEvent = `
		UPDATE events
		SET title=$1 starts_at=$2 ends_at=$3 description=$4 user_id=$5 notify_before=$6
		WHERE id=$7
	`
	DeleteEvent           = `DELETE FROM events WHERE id=$1`
	SelectEventsForPeriod = fmt.Sprintf(
		`SELECT %s FROM events WHERE start_at >= $1 AND start_at < $2`,
		eventAttributesWithID,
	)
	SelectEventsBeforeTime = fmt.Sprintf(
		`SELECT %s FROM events WHERE ends_at < $1`,
		eventAttributesWithID,
	)
	SelectEventsCreatedAfter = fmt.Sprintf(
		`SELECT %s FROM events WHERE created_at > $1`,
		eventAttributesWithID,
	)
)
