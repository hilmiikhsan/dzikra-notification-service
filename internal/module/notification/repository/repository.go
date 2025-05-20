package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/pkg/utils"
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

func (r *notificationTemplateRepository) FindNotificationByType(ctx context.Context, notificationType string) (*entity.NotificationType, error) {
	var res = new(entity.NotificationType)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindNotificationByType), notificationType)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("No notification type found")
			return nil, errors.New(constants.ErrNotificationTypeNotFound)
		}

		log.Error().Err(err).Msg("Failed to find notification type")
		return nil, err
	}

	return res, nil
}

func (r *notificationTemplateRepository) InsertNewNotification(ctx context.Context, tx *sqlx.Tx, data *entity.UserPushNotification) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewNotification),
		data.Title,
		data.Detail,
		data.Url,
		data.UserID,
		data.NTypeID,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewNotification - Failed to insert new notification")
		return err
	}

	return nil
}

func (r *notificationTemplateRepository) FindListNotification(ctx context.Context, limit, offset int, search string) ([]dto.NotificationDetail, int, error) {
	var responses []entity.UserPushNotification

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListNotification), search, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::FindListBanner - error executing query")
		return nil, 0, err
	}

	var total int

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListNotification), search); err != nil {
		log.Error().Err(err).Msg("repository::FindListVoucher - error counting banner")
		return nil, 0, err
	}

	notifications := make([]dto.NotificationDetail, 0, len(responses))
	for _, v := range responses {
		notifications = append(notifications, dto.NotificationDetail{
			ID:        v.ID,
			Title:     v.Title,
			Detail:    &v.Detail,
			Url:       v.Url,
			NTypeID:   v.NTypeID,
			UserID:    v.UserID.String(),
			CreatedAt: utils.FormatTime(v.CreatedAt),
		})
	}

	return notifications, total, nil
}
