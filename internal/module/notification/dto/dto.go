package dto

type InternalNotificationRequest struct {
	TemplateName string `validate:"required"`
	Recipient    string `validate:"required"`
	Placeholder  map[string]string
}
