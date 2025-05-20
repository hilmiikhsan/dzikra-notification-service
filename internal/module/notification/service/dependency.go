package service

import (
	"firebase.google.com/go/messaging"
	notificationTemplatePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/ports"
	notificationHistoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification_history/ports"
	"github.com/jmoiron/sqlx"
)

var _ notificationTemplatePorts.NotificationTemplateService = &notificationTemplateService{}

type notificationTemplateService struct {
	db                             *sqlx.DB
	notificationTemplateRepository notificationTemplatePorts.NotificationTemplateRepository
	notificationHistoryPorts       notificationHistoryPorts.NotificationHistoryRepository
	fcmClient                      *messaging.Client
}

func NewNotificationTemplateService(
	db *sqlx.DB,
	notificationTemplateRepository notificationTemplatePorts.NotificationTemplateRepository,
	notificationHistoryPorts notificationHistoryPorts.NotificationHistoryRepository,
	fcmClient *messaging.Client,
) *notificationTemplateService {
	return &notificationTemplateService{
		db:                             db,
		notificationTemplateRepository: notificationTemplateRepository,
		notificationHistoryPorts:       notificationHistoryPorts,
		fcmClient:                      fcmClient,
	}
}
