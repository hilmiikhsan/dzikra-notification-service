package entity

import "github.com/google/uuid"

type NotificationTemplate struct {
	ID           uuid.UUID `db:"id"`
	TemplateName string    `db:"template_name"`
	Subject      string    `db:"subject"`
	Body         string    `db:"body"`
}
