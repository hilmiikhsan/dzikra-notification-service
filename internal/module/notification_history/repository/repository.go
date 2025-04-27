package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification_history/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification_history/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.NotificationHistoryRepository = &notificationHistoryRepository{}

type notificationHistoryRepository struct {
	db *sqlx.DB
}

func NewNotificationHistoryRepository(db *sqlx.DB) *notificationHistoryRepository {
	return &notificationHistoryRepository{
		db: db,
	}
}

func (r *notificationHistoryRepository) InsertNotificationHistory(ctx context.Context, data *entity.NotificationHistory) error {
	_, err := r.db.DB.ExecContext(ctx, r.db.Rebind(queryInsertNotificationHistory),
		data.Recipient,
		data.TemplateID,
		data.Status,
		data.ErrorMessage,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert notification history")
		return err
	}

	return nil
}
