package entity

import "github.com/google/uuid"

type NotificationHistory struct {
	ID           uuid.UUID `db:"id"`
	Recipient    string    `db:"recipient"`
	TemplateID   uuid.UUID `db:"template_id"`
	Status       string    `db:"status"`
	ErrorMessage string    `db:"error_message"`
}
