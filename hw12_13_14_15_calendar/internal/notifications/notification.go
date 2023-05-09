package notifications

import (
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Notification struct {
	ID       storage.EventID `json:"id"`
	Title    string          `json:"title"`
	StartsAt time.Time       `json:"startsAt"`
	UserID   string          `json:"userId"`
}

// #### Уведомление
// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика, содержит поля:
// * ID события;
// * Заголовок события;
// * Дата события;
// * Пользователь, которому отправлять.
