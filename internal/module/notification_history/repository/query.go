package repository

const (
	queryInsertNotificationHistory = `
		INSERT INTO notification_histories 
		(
			recipient, 
			template_id, 
			status, 
			error_message
		) VALUES (?, ?, ?, ?)
	`
)
