package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.NotificationTemplateRepository = &notificationTemplateRepository{}

type notificationTemplateRepository struct {
	db *sqlx.DB
}

func NewNotificationTemplateRepository(db *sqlx.DB) *notificationTemplateRepository {
	return &notificationTemplateRepository{
		db: db,
	}
}

func (r *notificationTemplateRepository) FindTemplateByTemplateName(ctx context.Context, templateName string) (*entity.NotificationTemplate, error) {
	var res = new(entity.NotificationTemplate)

	err := r.db.GetContext(ctx, res, queryFindNotificationTemplateByTemplateName, templateName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find email template")
		return nil, err
	}

	return res, nil
}
