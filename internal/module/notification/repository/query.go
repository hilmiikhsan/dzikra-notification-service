package repository

const (
	queryFindNotificationTemplateByTemplateName = `
		SELECT
			id,
			template_name,
			subject,
			body
		FROM notification_templates
		WHERE template_name = $1
	`
)
