package mailer

import "fmt"

type HTMLBuilder struct{}

func NewHTMLBuilder() *HTMLBuilder {
	return &HTMLBuilder{}
}

func (b *HTMLBuilder) BuildWelcomeEmail(to, username string) (string, string, error) {
	subject := "Welcome to Life-RPG 🎉"
	body := fmt.Sprintf("<h1>Hello, %s!</h1><p>Welcome to our platform 🚀</p>", username)
	return subject, body, nil
}
