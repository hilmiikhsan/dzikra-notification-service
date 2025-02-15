package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/cmd/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/notification/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/notification/entity"
)

type NotificationTemplateRepository interface {
	FindTemplateByTemplateName(ctx context.Context, templateName string) (*entity.NotificationTemplate, error)
}

type NotificationTemplateService interface {
	SendEmail(ctx context.Context, req dto.InternalNotificationRequest) error
}

type NotificationEmailAPI interface {
	SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error)
}
