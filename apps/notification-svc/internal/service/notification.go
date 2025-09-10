package service

import "github.com/hassiimykyta/life-rpg/apps/notification-svc/internal/mailer"

type NotificationService struct {
	Builder mailer.MailBuilder
	Sender  *mailer.MailSender
}

func NewNotificationService(b mailer.MailBuilder, s *mailer.MailSender) *NotificationService {
	return &NotificationService{Builder: b, Sender: s}
}

func (n *NotificationService) SendWelcome(to, username string) error {
	subject, body, err := n.Builder.BuildWelcomeEmail(to, username)
	if err != nil {
		return err
	}
	return n.Sender.Send(to, subject, body)
}
