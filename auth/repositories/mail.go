package repositories

import (
	"crypto/tls"

	"github.com/gistsapp/api/auth/config"
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
		<p>Please click the following link to verify your email address:</p>
		<p><a href="http://localhost:3000/verify?token=`+token+`">Verify Email</a></p>
		<p>If you didn't request this email, please ignore it.</p>
	</body>
</html>`)

	d := gomail.NewDialer(e.config.Host, e.config.Port, e.config.User, e.config.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
