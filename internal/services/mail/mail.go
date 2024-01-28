package mail

import (
	"echo-sqlc-template/internal/config"
	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func Connect() {
	dialer = gomail.NewDialer(config.Data.Smtp.Host, config.Data.Smtp.Port, config.Data.Smtp.Username, config.Data.Smtp.Password)
}

func SendCode(recipient, code, title string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Data.Smtp.From)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/html", code)

	return dialer.DialAndSend(msg)
}
