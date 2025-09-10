package mailer

import "gopkg.in/gomail.v2"

type MailSender struct {
	dialer *gomail.Dialer
	from   string
}

func New(dialer *gomail.Dialer, from string) *MailSender {
	return &MailSender{dialer: dialer, from: from}
}

func (s *MailSender) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}
