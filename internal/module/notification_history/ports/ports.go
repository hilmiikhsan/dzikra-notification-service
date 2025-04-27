package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification_history/entity"
)

type NotificationHistoryRepository interface {
	InsertNotificationHistory(ctx context.Context, data *entity.NotificationHistory) error
}
