package external

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type Email struct {
	To              string
	Subject         string
	Body            string
	FullName        string
	Email           string
	IsStatusChanged bool
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

func (e *Email) SenFcmNotificationEmail() error {
	intro := "Your order has been processed successfully."
	if e.IsStatusChanged {
		intro = "Your order status has been changed."
	}

	const tpl = `
        <p>Hi {{ .Name }},</p>
        <p>{{ .Intro }}</p>
        <p><a href="https://ecommerce.example.com/login" style="background:#d4ba38;padding:8px 12px;color:#fff;text-decoration:none;border-radius:4px;">Go to Dashboard</a></p>
        <p>Thank you for your purchase.</p>
    `

	data := struct {
		Name  string
		Intro string
	}{e.FullName, intro}

	var bodyBuf bytes.Buffer
	t := template.Must(template.New("email").Parse(tpl))
	if err := t.Execute(&bodyBuf, data); err != nil {
		log.Error().Err(err).Msg("Failed to execute template")
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", e.FullName)
	mailer.SetHeader("To", e.Email)
	mailer.SetHeader("Subject", "Order Status")
	mailer.SetBody("text/html", bodyBuf.String())

	smtpPort := config.Envs.Notification.MailPort
	intSmtpPort, _ := strconv.Atoi(smtpPort)

	dialer := gomail.NewDialer(
		config.Envs.Notification.MailHost,
		intSmtpPort,
		config.Envs.Notification.MailUser,
		config.Envs.Notification.MailPassword,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Error().Err(err).Msg("Failed to send email")
		return err
	}

	log.Info().Msgf("Email sent to %s", e.Email)

	return nil
}
