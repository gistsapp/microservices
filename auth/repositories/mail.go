package repositories

import (
	"crypto/tls"

	"github.com/gistsapp/api/auth/config"
	"github.com/gofiber/fiber/v2/log"
	gomail "gopkg.in/mail.v2"
)

type EmailService interface {
	SendVerificationEmail(email string, value string) error
}

type emailService struct {
	config config.EmailServiceConfig
}

func NewEmailService(config config.EmailServiceConfig) EmailService {
	return emailService{
		config: config,
	}
}

func (e emailService) SendVerificationEmail(email string, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.User)
	m.SetHeader("To", email)

	m.SetBody("text/html", `<html>
	<head></head>
	<body>
		<p>Hello,</p>
		<p>Please enter the following code to verify your email address:</p>
		<p>`+token+`</p>
		<p>If you didn't request this email, please ignore it.</p>
	</body>
</html>`)

	d := gomail.NewDialer(e.config.Host, e.config.Port, e.config.User, e.config.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Error(err)
		return err
	}
	log.Info("Email sent to " + email)
	return nil
}
