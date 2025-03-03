package service

import (
	"bytes"
	"context"
	"text/template"

	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/external"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/notification/dto"
	notificationTemplatePorts "github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/notification/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/notification_history/entity"
	notificationHistoryPorts "github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/notification_history/ports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var _ notificationTemplatePorts.NotificationTemplateService = &notificationTemplateService{}

type notificationTemplateService struct {
	notificationTemplateRepository notificationTemplatePorts.NotificationTemplateRepository
	notificationHistoryPorts       notificationHistoryPorts.NotificationHistoryRepository
}

func NewNotificationTemplateService(notificationTemplateRepository notificationTemplatePorts.NotificationTemplateRepository, notificationHistoryPorts notificationHistoryPorts.NotificationHistoryRepository) *notificationTemplateService {
	return &notificationTemplateService{
		notificationTemplateRepository: notificationTemplateRepository,
		notificationHistoryPorts:       notificationHistoryPorts,
	}
}

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
