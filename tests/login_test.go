package tests

import (
	"context"
	"os"
	"testing"

	pb "goAuth/shared/proto"
)

func TestLogin(t *testing.T) {
	db, authServer := setupTest(t)

	testEmail := os.Getenv("SMTP_EMAIL")
	if testEmail == "" {
		t.Skip("SMTP_EMAIL not set, skipping test")
	}

	// Clean up any existing user first
	cleanupUser(t, db, testEmail)
	cleanupOTP(t, testEmail)

	// First signup to create user
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

	// Then test login
	loginResp, err := authServer.Login(context.Background(), &pb.LoginRequest{
		Email:    testEmail,
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if loginResp.Message != "Login successful" {
		t.Errorf("unexpected login message: %v", loginResp.Message)
	}

	if loginResp.Token == "" {
		t.Error("Expected token to be generated")
	}

	// Cleanup after test (before closing DB)
	cleanupUser(t, db, testEmail)
	cleanupOTP(t, testEmail)
	db.Close()
}