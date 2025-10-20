package tests

import (
	"testing"

	"goAuth/config"
	"goAuth/mail"
)

func TestSendMail(t *testing.T) {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		t.Fatalf("Failed to load .env: %v", err)
	}

	// Test sending mail
	err := mail.SendMail(
		"vinyasbharadwaj101@gmail.com", // Replace with your test email
		"Test Email from goAuth",
		"This is a test email to verify the mail functionality is working correctly.",
	)

	if err != nil {
		t.Fatalf("Failed to send mail: %v", err)
	}

	t.Log("Mail sent successfully!")
}