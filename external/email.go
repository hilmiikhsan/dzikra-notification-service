package external

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/infrastructure/config"
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
	mailer.SetHeader("From", config.Envs.Notification.SMTPuthEmail)
	mailer.SetHeader("Subject", e.Subject)
	mailer.SetBody("text/html", e.Body)

	smtpPort := config.Envs.Notification.SMTPPort
	intSmtpPort, _ := strconv.Atoi(smtpPort)

	dialer := gomail.NewDialer(
		config.Envs.Notification.SMTPHost,
		intSmtpPort,
		config.Envs.Notification.SMTPuthEmail,
		config.Envs.Notification.SMTPAuthPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send email")
		return err
	}

	log.Info().Msg("Email sent successfully")

	return nil
}
