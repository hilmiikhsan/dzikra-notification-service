package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/cmd/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/ports"
	notificationTemplateRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/service"
	notificationHistoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification_history/repository"
	"github.com/gofiber/fiber/v2/log"
)

type NotificationEmailAPI struct {
	EmailService ports.NotificationTemplateService
	notification.UnimplementedNotificationServiceServer
}

func NewNotificationEmailAPI() *NotificationEmailAPI {
	var handler = new(NotificationEmailAPI)

	notificationTemplateRepository := notificationTemplateRepository.NewNotificationTemplateRepository(adapter.Adapters.DzikraPostgres)
	notificationHistoryRepository := notificationHistoryRepository.NewNotificationHistoryRepository(adapter.Adapters.DzikraPostgres)

	emailService := service.NewNotificationTemplateService(notificationTemplateRepository, notificationHistoryRepository)

	handler.EmailService = emailService

	return handler
}

func (api *NotificationEmailAPI) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	var (
		validators = adapter.Adapters.Validator
	)

	internalReq := dto.InternalNotificationRequest{
		TemplateName: req.TemplateName,
		Recipient:    req.Recipient,
		Placeholder:  req.Placeholders,
	}

	if err := validators.Validate(internalReq); err != nil {
		log.Error("error validating request: ", err)
		return &notification.SendNotificationResponse{
			Message: "failed to validate request",
		}, nil
	}

	err := api.EmailService.SendEmail(ctx, internalReq)
	if err != nil {
		log.Error("error sending email: ", err)
		return &notification.SendNotificationResponse{
			Message: "failed to send email",
		}, nil
	}

	return &notification.SendNotificationResponse{
		Message: "success",
	}, nil
}
