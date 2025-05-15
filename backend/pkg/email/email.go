package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
)

type Email struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func NewEmailSender(host string, port string, username string, password string, sender string) *Email {
	return &Email{
		SMTPHost:     host,
		SMTPPort:     port,
		SMTPUsername: username,
		SMTPPassword: password,
		SMTPFrom:     sender,
	}
}

func (e *Email) SendEmail(payload interface{}, pathTemplate string, to []string, subject string) error {

	absPath, err := filepath.Abs(pathTemplate)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	t, err := template.ParseFiles(absPath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	buff := new(bytes.Buffer)
	if err := t.Execute(buff, payload); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	msg := "To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		buff.String()

	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)

	if err := smtp.SendMail(e.SMTPHost+":"+e.SMTPPort, auth, e.SMTPFrom, to, []byte(msg)); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
