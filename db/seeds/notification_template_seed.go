package seeds

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// NotificationTemplatesSeed seeds the `notification_templates` table.
func (s *Seed) notificationTemplatesSeed() {
	templates := []map[string]any{
		{
			"template_name": "register",
			"subject":       "Welcome to Dzikra App",
			"body":          "Hi welcome to our service! We're glad to have you.",
		},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Check if `notification_templates` table already contains data
	var count int
	err = tx.Get(&count, `SELECT COUNT(*) FROM notification_templates`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking notification_templates table")
		return
	}
	if count > 0 {
		log.Info().Msg("Notification templates already seeded")
		return
	}

	// Insert data into `notification_templates`
	insertTemplateQuery := `
		INSERT INTO notification_templates 
		(template_name, subject, body) 
		VALUES (:template_name, :subject, :body)`
	for _, template := range templates {
		template["id"] = uuid.New().String()
		_, err = tx.NamedExec(insertTemplateQuery, template)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting notification template")
			return
		}
	}

	log.Info().Msg("notification_templates table seeded successfully")
}
