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

	// Try to load environment variables, but don't fail if .env doesn't exist
	if err := config.LoadEnv(); err != nil {
		t.Logf("Warning: Could not load .env file: %v", err)
		// In CI, environment variables should be set directly
	}

	// Check if required environment variables are set
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	
	if smtpEmail == "" || smtpPassword == "" {
		t.Skip("SMTP_EMAIL or SMTP_PASSWORD not set, skipping mail test")
	}

	// Test sending mail
	err := mail.SendMail(
		smtpEmail,
		"Test Email from goAuth",
		"This is a test email to verify the mail functionality is working correctly.",
	)

	if err != nil {
		t.Fatalf("Failed to send mail: %v", err)
	}

	t.Log("Mail sent successfully!")
}