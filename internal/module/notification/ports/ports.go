package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/cmd/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/entity"
	"github.com/jmoiron/sqlx"
)

type NotificationTemplateRepository interface {
	FindTemplateByTemplateName(ctx context.Context, templateName string) (*entity.NotificationTemplate, error)
	FindNotificationByType(ctx context.Context, notificationType string) (*entity.NotificationType, error)
	InsertNewNotification(ctx context.Context, tx *sqlx.Tx, data *entity.UserPushNotification) error
	FindListNotification(ctx context.Context, limit, offset int, search string) ([]dto.NotificationDetail, int, error)
}

type NotificationTemplateService interface {
	SendEmail(ctx context.Context, req dto.InternalNotificationRequest) error
	GetNotificationByType(ctx context.Context, notificationType string) (*dto.GetNotificationByTypeResponse, error)
	CreateNotification(ctx context.Context, req dto.CreateNotificationRequest) error
	GetListNotification(ctx context.Context, page, limit int, search string) (*dto.GetListNotificationResponse, error)
	SendFcmBatchNotification(ctx context.Context, req *dto.SendBatchFcmNotificationRequest) error
	SendFcmNotification(ctx context.Context, req *dto.SendFcmNotificationRequest) error
	SendTransactionEmail(ctx context.Context, req *dto.SendTransactionEmailRequest) error
}

type NotificationEmailAPI interface {
	SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error)
	GetNotificationByType(ctx context.Context, req *notification.GetNotificationByTypeRequest) (*notification.GetNotificationByTypeResponse, error)
	CreateNotification(ctx context.Context, req *notification.CreateNotificationRequest) (*notification.CreateNotificationResponse, error)
	GetListNotification(ctx context.Context, req *notification.GetListNotificationRequest) (*notification.GetListNotificationResponse, error)
	SendFcmBatchNotification(ctx context.Context, req *notification.SendFcmBatchNotificationRequest) (*notification.SendFcmBatchNotificationResponse, error)
	SendFcmNotification(ctx context.Context, req *notification.SendFcmNotificationRequest) (*notification.SendFcmNotificationResponse, error)
	SendTransactionEmail(ctx context.Context, req *notification.SendTransactionEmailRequest) (*notification.SendTransactionEmailResponse, error)
}
