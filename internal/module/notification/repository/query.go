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

	queryFindNotificationByType = `
		SELECT
			id,
			type,
			name
		FROM notification_types
		WHERE type = ? AND deleted_at IS NULL
	`

	queryInsertNewNotification = `
		INSERT INTO user_push_notifications 
		(
			title, 
			detail,
			url,
			user_id,
			n_type_id
		) VALUES (?, ?, ?, ?, ?)
	`

	queryFindListNotification = `
		SELECT
			id,
			title,
			detail,
			url,
			user_id,
			n_type_id,
			created_at
		FROM user_push_notifications
		WHERE deleted_at IS NULL
		AND concat(
				title,
				' ',
				coalesce(detail, ''),
				' ',
				coalesce(url, '')
			) ILIKE '%' || ? || '%'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListNotification = `
	SELECT
		COUNT(*) AS total
	FROM user_push_notifications
	WHERE deleted_at IS NULL
	AND concat(
			title,
			' ',
			coalesce(detail, ''),
			' ',
			coalesce(url, '')
		) ILIKE '%' || ? || '%'
	`
)
