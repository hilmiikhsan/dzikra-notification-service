package service

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"firebase.google.com/go/messaging"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/external"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/dto"
	notification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification_history/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (s *notificationTemplateService) SendEmail(ctx context.Context, req dto.InternalNotificationRequest) error {
	emailTemplate, err := s.notificationTemplateRepository.FindTemplateByTemplateName(ctx, req.TemplateName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find email template")
		return errors.Wrap(err, "failed to find email template")
	}

	tmpl, err := template.New("emailTemplate").Parse(emailTemplate.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse email template")
		return errors.Wrap(err, "failed to parse email template")
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, map[string]string{
		"full_name":      req.Placeholder["full_name"],
		"otp_number":     req.Placeholder["otp_number"],
		"url_reset_link": req.Placeholder["url_reset_link"],
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute email template with placeholder")
		return errors.Wrap(err, "failed to execute email template with placeholder")
	}

	email := external.Email{
		To:      req.Recipient,
		Subject: emailTemplate.Subject,
		Body:    tpl.String(),
	}

	err = email.SendEmail()
	if err != nil {
		log.Error().Err(err).Msg("Failed to send email")

		notifHistory := &entity.NotificationHistory{
			Recipient:    req.Recipient,
			TemplateID:   emailTemplate.ID,
			Status:       constants.FailedSendNotification,
			ErrorMessage: err.Error(),
		}

		err = s.notificationHistoryPorts.InsertNotificationHistory(ctx, notifHistory)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert notification history")
			return errors.Wrap(err, "failed to insert notification history")
		}

		return errors.Wrap(err, "failed to send email")
	}

	notifHistory := &entity.NotificationHistory{
		Recipient:  req.Recipient,
		TemplateID: emailTemplate.ID,
		Status:     constants.SuccessSendNotification,
	}

	err = s.notificationHistoryPorts.InsertNotificationHistory(ctx, notifHistory)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert notification history")
		return errors.Wrap(err, "failed to insert notification history")
	}

	return nil
}

func (s *notificationTemplateService) GetNotificationByType(ctx context.Context, notificationType string) (*dto.GetNotificationByTypeResponse, error) {
	res, err := s.notificationTemplateRepository.FindNotificationByType(ctx, notificationType)
	if err != nil {
		if err.Error() == constants.ErrNotificationTypeNotFound {
			log.Error().Err(err).Msg("No notification type found")
			return nil, errors.New(constants.ErrNotificationTypeNotFound)
		}

		log.Error().Err(err).Msg("Failed to find notification type")
		return nil, errors.Wrap(err, "failed to find notification type")
	}

	return &dto.GetNotificationByTypeResponse{
		ID:   res.ID,
		Type: res.Type,
		Name: res.Name,
	}, nil
}

func (s *notificationTemplateService) CreateNotification(ctx context.Context, req dto.CreateNotificationRequest) error {
	_, err := s.notificationTemplateRepository.FindNotificationByType(ctx, strings.ToUpper(req.NTypeID))
	if err != nil {
		if err.Error() == constants.ErrNotificationTypeNotFound {
			log.Error().Err(err).Msg("No notification type found")
			return errors.New(constants.ErrNotificationTypeNotFound)
		}

		log.Error().Err(err).Msg("Failed to find notification type")
		return errors.Wrap(err, "failed to find notification type")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user ID")
		return errors.Wrap(err, "failed to parse user ID")
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateNotification - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateNotification - Failed to rollback transaction")
			}
		}
	}()

	err = s.notificationTemplateRepository.InsertNewNotification(ctx, tx, &notification.UserPushNotification{
		Title:   req.Title,
		Detail:  req.Detail,
		Url:     req.Url,
		UserID:  userID,
		NTypeID: strings.ToUpper(req.NTypeID),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create notification")
		return errors.Wrap(err, "failed to create notification")
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateNotification - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (s *notificationTemplateService) GetListNotification(ctx context.Context, page, limit int, search string) (*dto.GetListNotificationResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list notification
	notifications, total, err := s.notificationTemplateRepository.FindListNotification(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListNotification - error getting list notification")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if notifications is nil
	if notifications == nil {
		notifications = []dto.NotificationDetail{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListNotificationResponse{
		Notification: notifications,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PageSize:     perPage,
		TotalData:    total,
	}

	// return response
	return &response, nil
}

func (s *notificationTemplateService) SendFcmBatchNotification(ctx context.Context, req *dto.SendBatchFcmNotificationRequest) error {
	msg := &messaging.MulticastMessage{
		Tokens: req.FcmToken,
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
	}

	resp, err := s.fcmClient.SendMulticast(ctx, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send FCM batch notification")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Info().Msgf("firebase: sent %d/%d messages", resp.SuccessCount, len(req.FcmToken))

	return nil
}

func (s *notificationTemplateService) SendFcmNotification(ctx context.Context, req *dto.SendFcmNotificationRequest) error {
	go func(token, title, body string) {
		if req.FcmToken != "" {
			bg := context.Background()
			msg := &messaging.Message{
				Notification: &messaging.Notification{
					Title: title,
					Body:  body,
				},
				Token: token,
			}

			id, err := s.fcmClient.Send(bg, msg)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send FCM notification")
				return
			}

			log.Info().Msgf("fcm sent messageID=%s", id)

			return
		}

		email := external.Email{
			FullName:        req.FullName,
			Email:           req.Email,
			IsStatusChanged: req.IsStatusChanged,
		}

		err := email.SenFcmNotificationEmail()
		if err != nil {
			log.Error().Err(err).Msg("Failed to send email notification")
			return
		}
	}(req.FcmToken, req.Title, req.Body)

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::SendFcmNotification - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::SendFcmNotification - Failed to rollback transaction")
			}
		}
	}()

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("service::SendFcmNotification - Failed to parse user ID")
		return errors.Wrap(err, "failed to parse user ID")
	}

	err = s.notificationTemplateRepository.InsertNewNotification(ctx, tx, &notification.UserPushNotification{
		Title:   req.Title,
		Detail:  req.Body,
		Url:     "",
		UserID:  userID,
		NTypeID: "INFO",
	})
	if err != nil {
		log.Error().Err(err).Msg("serve::SendFcmNotification - Failed to create notification")
		return errors.Wrap(err, "failed to create notification")
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::SendFcmNotification - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (s *notificationTemplateService) SendTransactionEmail(ctx context.Context, req *dto.SendTransactionEmailRequest) error {
	email := external.Email{
		FullName:        req.TOName,
		Email:           req.TOEmail,
		IsStatusChanged: req.IsStatusChanged,
	}

	err := email.SendTransactionEmail(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send transaction email")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
