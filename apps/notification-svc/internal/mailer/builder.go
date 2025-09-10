package mailer

type MailBuilder interface {
	BuildWelcomeEmail(to string, username string) (subject string, body string, err error)
	// BuildResetPasswordEmail(to string, token string) (subject string, body string, err error)
}
