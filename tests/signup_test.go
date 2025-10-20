package tests

import (
	"context"
	"os"
	"testing"

	pb "goAuth/proto"
)

func TestSignup(t *testing.T) {
	db, authServer := setupTest(t)
	
	testEmail := os.Getenv("SMTP_EMAIL")
	if testEmail == "" {
		t.Skip("SMTP_EMAIL not set, skipping test")
	}

	// Clean up any existing user first
	cleanupUser(t, db, testEmail)
	cleanupOTP(t, testEmail)

	// Test Signup
	signupResp, err := authServer.Signup(context.Background(), &pb.SignupRequest{
		Email:    testEmail,
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}

	if signupResp.Message != "An OTP has been shared to your mail" {
		t.Errorf("unexpected message: %v", signupResp.Message)
	}

	// Cleanup after test (before closing DB)
	cleanupUser(t, db, testEmail)
	cleanupOTP(t, testEmail)
	db.Close()
}