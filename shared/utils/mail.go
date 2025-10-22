package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	
	if from == "" || password == "" {
		return fmt.Errorf("missing email credentials")
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	
	return nil
}