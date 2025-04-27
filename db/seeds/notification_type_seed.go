package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// notificationType seeds the `notification_types` table.
func (s *Seed) notificationType() {
	notificationTypeMaps := []map[string]any{
		{
			"type": "INFO",
			"name": "info",
		},
		{
			"type": "PROMO",
			"name": "promo",
		},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Check if notification type already exist
	var count int
	err = tx.Get(&count, `SELECT COUNT(id) FROM notification_types`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking notification_types table")
		return
	}
	if count > 0 {
		log.Info().Msg("notification type  already seeded")
		return
	}

	insertNotificationTypeQuery := `INSERT INTO notification_types (id, type, name) VALUES (:id, :type, :name)`
	for _, notificationType := range notificationTypeMaps {
		uuid, _ := utils.GenerateUUIDv7String()
		notificationType["id"] = uuid
		_, err = tx.NamedExec(insertNotificationTypeQuery, notificationType)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting notification type ")
			return
		}
	}

	log.Info().Msg("Notification Type table seeded successfully")
}
