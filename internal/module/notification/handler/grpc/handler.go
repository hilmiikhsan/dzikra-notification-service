package grpc

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/cmd/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/integration/firebase"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/ports"
	notificationTemplateRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/service"
	notificationHistoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification_history/repository"
	"github.com/rs/zerolog/log"
)

type NotificationEmailAPI struct {
	EmailService ports.NotificationTemplateService
	notification.UnimplementedNotificationServiceServer
}

func NewNotificationEmailAPI() *NotificationEmailAPI {
	var handler = new(NotificationEmailAPI)

	// firebase messaging
	fcmClient := firebase.InitFirebaseMessaging()

	notificationTemplateRepository := notificationTemplateRepository.NewNotificationTemplateRepository(adapter.Adapters.DzikraPostgres)
	notificationHistoryRepository := notificationHistoryRepository.NewNotificationHistoryRepository(adapter.Adapters.DzikraPostgres)

	emailService := service.NewNotificationTemplateService(
		adapter.Adapters.DzikraPostgres,
		notificationTemplateRepository,
		notificationHistoryRepository,
		fcmClient,
	)

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
		log.Error().Err(err).Msg("error validating request")
		return &notification.SendNotificationResponse{
			Message: "failed to validate request",
		}, nil
	}

	err := api.EmailService.SendEmail(ctx, internalReq)
	if err != nil {
		log.Error().Err(err).Msg("error sending email")
		return &notification.SendNotificationResponse{
			Message: "failed to send email",
		}, nil
	}

	return &notification.SendNotificationResponse{
		Message: "success",
	}, nil
}

func (api *NotificationEmailAPI) GetNotificationByType(ctx context.Context, req *notification.GetNotificationByTypeRequest) (*notification.GetNotificationByTypeResponse, error) {
	res, err := api.EmailService.GetNotificationByType(ctx, req.Type)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrNotificationTypeNotFound) {
			log.Error().Err(err).Msg("notification type not found")
			return &notification.GetNotificationByTypeResponse{
				Message: constants.ErrNotificationTypeNotFound,
			}, nil
		}

		log.Error().Err(err).Msg("error getting notification by type")
		return &notification.GetNotificationByTypeResponse{
			Message: "failed to get notification by type",
		}, nil
	}

	return &notification.GetNotificationByTypeResponse{
		Id:      res.ID.String(),
		Type:    res.Type,
		Name:    res.Name,
		Message: "success",
	}, nil
}

func (api *NotificationEmailAPI) CreateNotification(ctx context.Context, req *notification.CreateNotificationRequest) (*notification.CreateNotificationResponse, error) {
	err := api.EmailService.CreateNotification(ctx, dto.CreateNotificationRequest{
		Title:   req.Title,
		Detail:  req.Detail,
		Url:     req.Url,
		UserID:  req.UserId,
		NTypeID: req.NTypeId,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrNotificationTypeNotFound) {
			log.Error().Err(err).Msg("notification type not found")
			return &notification.CreateNotificationResponse{
				Message: constants.ErrNotificationTypeNotFound,
			}, nil
		}

		log.Error().Err(err).Msg("error creating notification")
		return &notification.CreateNotificationResponse{
			Message: "failed to create notification",
		}, nil
	}

	return &notification.CreateNotificationResponse{
		Message: "success",
	}, nil
}

func (api *NotificationEmailAPI) GetListNotification(ctx context.Context, req *notification.GetListNotificationRequest) (*notification.GetListNotificationResponse, error) {
	res, err := api.EmailService.GetListNotification(ctx, int(req.Page), int(req.Limit), req.Search)
	if err != nil {
		log.Error().Err(err).Msg("error getting notification list")
		return &notification.GetListNotificationResponse{
			Message: "failed to get notification by type",
		}, nil
	}

	var notifications []*notification.NotificationDetail
	for _, v := range res.Notification {
		notifications = append(notifications, &notification.NotificationDetail{
			Id:        int64(v.ID),
			Title:     v.Title,
			Detail:    *v.Detail,
			Url:       v.Url,
			NTypeId:   v.NTypeID,
			UserId:    v.UserID,
			CreatedAt: v.CreatedAt,
		})
	}

	return &notification.GetListNotificationResponse{
		Notifications: notifications,
		TotalPage:     int32(res.TotalPages),
		CurrentPage:   int32(res.CurrentPage),
		PageSize:      int32(res.PageSize),
		TotalData:     int32(res.TotalData),
		Message:       "success",
	}, nil
}

func (api *NotificationEmailAPI) SendFcmBatchNotification(ctx context.Context, req *notification.SendFcmBatchNotificationRequest) (*notification.SendFcmBatchNotificationResponse, error) {
	err := api.EmailService.SendFcmBatchNotification(ctx, &dto.SendBatchFcmNotificationRequest{
		FcmToken: req.FcmTokens,
		Title:    req.Title,
		Body:     req.Body,
	})
	if err != nil {
		log.Error().Err(err).Msg("error sending fcm batch notification")
		return &notification.SendFcmBatchNotificationResponse{
			Message: "failed to send fcm batch notification",
		}, nil
	}

	return &notification.SendFcmBatchNotificationResponse{
		Message: "success",
	}, nil
}

func (api *NotificationEmailAPI) SendFcmNotification(ctx context.Context, req *notification.SendFcmNotificationRequest) (*notification.SendFcmNotificationResponse, error) {
	err := api.EmailService.SendFcmNotification(ctx, &dto.SendFcmNotificationRequest{
		FcmToken:        req.FcmToken,
		Title:           req.Title,
		Body:            req.Body,
		UserID:          req.UserId,
		IsStatusChanged: req.IsStatusChanged,
		FullName:        req.FullName,
		Email:           req.Email,
	})
	if err != nil {
		log.Error().Err(err).Msg("error sending fcm notification")
		return &notification.SendFcmNotificationResponse{
			Message: "failed to send fcm notification",
		}, nil
	}

	return &notification.SendFcmNotificationResponse{
		Message: "success",
	}, nil
}
