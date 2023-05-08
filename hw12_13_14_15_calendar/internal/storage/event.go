package storage

import "time"

type EventID string

type Event struct {
	ID           EventID
	Title        string
	StartsAt     time.Time
	EndsAt       time.Time
	CreatedAt    time.Time
	Description  string
	UserID       string
	NotifyBefore time.Duration
}

// Событие - основная сущность, содержит в себе поля:
// * ID - уникальный идентификатор события (можно воспользоваться UUID);
// * Заголовок - короткий текст;
// * Дата и время события;
// * Длительность события (или дата и время окончания);
// * Описание события - длинный текст, опционально;
// * ID пользователя, владельца события;
// * За сколько времени высылать уведомление, опционально.
