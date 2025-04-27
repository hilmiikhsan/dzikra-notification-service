package external

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type Email struct {
	To      string
	Subject string
	Body    string
}

func (e *Email) SendEmail() error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("To", e.To)
	mailer.SetHeader("From", config.Envs.Notification.MailUser)
	mailer.SetHeader("Subject", e.Subject)
	mailer.SetBody("text/html", e.Body)

	smtpPort := config.Envs.Notification.MailPort
	intSmtpPort, _ := strconv.Atoi(smtpPort)

	dialer := gomail.NewDialer(
		config.Envs.Notification.MailHost,
		intSmtpPort,
		config.Envs.Notification.MailUser,
		config.Envs.Notification.MailPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send email")
		return err
	}

	log.Info().Msg("Email sent successfully")

	return nil
}
