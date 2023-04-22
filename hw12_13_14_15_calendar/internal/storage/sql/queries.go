package sqlstorage

import "fmt"

var (
	eventAttributes       = `title, starts_at, ends_at, description, user_id, notify_before`
	eventAttributesWithID = `id,` + eventAttributes
	SelectEventByID       = fmt.Sprintf(`SELECT %s FROM events WHERE id=?`, eventAttributesWithID)
	InsertEvent           = fmt.Sprintf(`INSERT INTO events(%s) VALUES(?,?,?,?,?,?)`, eventAttributes)
	UpdateEvent           = `UPDATE events SET title=? starts_at=? ends_at=? description=? user_id=? notify_before=? WHERE id=?`
	DeleteEvent           = `DELETE FROM events WHERE id=?`
	SelectEventsForPeriod = fmt.Sprintf(`SELECT %s FROM events WHERE start_date >= ? AND start_date < ?`, eventAttributesWithID)
)
