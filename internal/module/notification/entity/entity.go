package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationTemplate struct {
	ID           uuid.UUID `db:"id"`
	TemplateName string    `db:"template_name"`
	Subject      string    `db:"subject"`
	Body         string    `db:"body"`
}

type NotificationType struct {
	ID     uuid.UUID `db:"id"`
	Type   string    `db:"type"`
	Name   string    `db:"name"`
	UserID uuid.UUID `db:"user_id"`
}

type UserPushNotification struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Detail    string    `db:"detail"`
	Url       string    `db:"url"`
	UserID    uuid.UUID `db:"user_id"`
	NTypeID   string    `db:"n_type_id"`
	CreatedAt time.Time `db:"created_at"`
}
