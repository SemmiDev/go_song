package mail

import (
	"log"

	mail "github.com/xhit/go-simple-mail/v2"
)

type SMTP struct {
	client *mail.SMTPClient
}

func NewSTMP() *SMTP {
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = "-"
	server.Password = "-"
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	return &SMTP{smtpClient}
}

func (s *SMTP) Send(receiver string, subject string, body string) {
	email := mail.NewMSG()
	email.SetFrom("From Me <-@gmail.com>")
	email.AddTo(receiver)
	email.AddCc("-@gmail.com")
	email.SetSubject(subject)

	email.SetBody(mail.TextHTML, body)
	email.Send(s.client)
}
