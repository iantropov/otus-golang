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

// Событие - основная сущность, содержит в себе поля:
// * ID - уникальный идентификатор события (можно воспользоваться UUID);
// * Заголовок - короткий текст;
// * Дата и время события;
// * Длительность события (или дата и время окончания);
// * Описание события - длинный текст, опционально;
// * ID пользователя, владельца события;
// * За сколько времени высылать уведомление, опционально.

// type Event struct {
// 	ID           EventID
// 	Title        string
// 	StartsAt     time.Time
// 	EndsAt       time.Time
// 	Description  string
// 	UserID       string
// 	NotifyBefore time.Duration
// }
