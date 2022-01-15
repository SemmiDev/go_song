package mail

type Mailer interface {
	Send(receiver string, subject string, body string)
}
