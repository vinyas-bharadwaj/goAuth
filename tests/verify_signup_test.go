package tests

import (
	"context"
	"os"
	"testing"

	"goAuth/shared/config"
	pb "goAuth/shared/proto"
)

func TestVerifySignup(t *testing.T) {
	db, authServer := setupTest(t)

	testEmail := os.Getenv("SMTP_EMAIL")
	if testEmail == "" {
		t.Skip("SMTP_EMAIL not set, skipping test")
	}

	// Clean up any existing user and OTP first
	cleanupUser(t, db, testEmail)
	cleanupOTP(t, testEmail)

	// First signup to create user and OTP
	signupResp, err := authServer.Signup(context.Background(), &pb.SignupRequest{
		Email:    testEmail,
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}

	if signupResp.Message != "An OTP has been shared to your mail" {
		t.Errorf("unexpected signup message: %v", signupResp.Message)
	}

	// Get the OTP from Redis for testing
	storedOTP, err := config.Rdb.Get(context.Background(), "otp:"+testEmail).Result()
	if err != nil {
		t.Fatalf("Failed to retrieve OTP from Redis: %v", err)
	}

	// Test VerifySignup with the actual OTP
	verifySignupResp, err := authServer.VerifySignup(context.Background(), &pb.VerifySignupRequest{
		Email: testEmail,
		Otp:   storedOTP,
	})

	if err != nil {
		t.Fatalf("Failed to verify OTP: %v", err)
	}

	if verifySignupResp.Message != "Signup successful" {
		t.Errorf("unexpected verify message: %v", verifySignupResp.Message)
	}

	if verifySignupResp.Token == "" {
		t.Error("Expected token to be generated")
	}

	// Cleanup after test (before closing DB)
	cleanupUser(t, db, testEmail)
	cleanupOTP(t, testEmail)
	db.Close()
}