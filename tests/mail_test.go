package tests

import (
	"testing"
	"os"

	"goAuth/config"
	"goAuth/mail"
)

func TestSendMail(t *testing.T) {
	// Skip in CI environment to avoid timeout issues
	if testing.Short() {
		t.Skip("Skipping mail test in CI environment")
	}

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		t.Fatalf("Failed to load .env: %v", err)
	}

	// Test sending mail
	err := mail.SendMail(
		os.Getenv("SMTP_EMAIL"),
		"Test Email from goAuth",
		"This is a test email to verify the mail functionality is working correctly.",
	)

	if err != nil {
		t.Fatalf("Failed to send mail: %v", err)
	}

	t.Log("Mail sent successfully!")
}