package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"goAuth/shared/config"
	pb "goAuth/shared/proto"
	"goAuth/shared/utils"

	_ "github.com/mattn/go-sqlite3"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

// Copy of handler methods for testing
func (s *AuthServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	_, err := s.db.Exec("INSERT INTO users(email, password) VALUES(?, ?)", req.Email, string(hash))
	if err != nil {
		return nil, errors.New("Email already exists!")
	}

	// Generate OTP
	otp := utils.GenerateOTP()

	// Store the OTP in redis with a 5 min time limit
	err = config.Rdb.Set(ctx, "otp:"+req.Email, otp, 5*time.Minute).Err()
	if err != nil {
		return nil, errors.New("Error in saving OTP to redis")
	}

	// Skip actual email sending in CI environment
	if !testing.Short() {
		subject := "Your Signup OTP"
		body := fmt.Sprintf("<p>Your OTP is <b>%s</b>. It expires in 5 minutes.</p>", otp)
		if err := utils.SendMail(req.Email, subject, body); err != nil {
			return nil, errors.New("Error sending OTP")
		}
	}

	return &pb.SignupResponse{Message: "An OTP has been shared to your mail"}, nil
}

func (s *AuthServer) VerifySignup(ctx context.Context, req *pb.VerifySignupRequest) (*pb.AuthResponse, error) {
	storedOTP, err := config.Rdb.Get(ctx, "otp:"+req.Email).Result()
	if err != nil {
		return nil, errors.New("Error retrieving OTP from redis")
	}

	if storedOTP != req.Otp {
		return nil, errors.New("The entered OTP does not match the one sent to your mail")
	}

	// Delete OTP after verification
	config.Rdb.Del(ctx, "otp:"+req.Email)

	// For testing, we'll use a simple JWT generation
	token := "test_jwt_token_" + req.Email

	return &pb.AuthResponse{Token: token, Message: "Signup successful"}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	row := s.db.QueryRow("SELECT password FROM users WHERE email=?", req.Email)
	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid password")
	}
	
	// For testing, we'll use a simple JWT generation
	token := "test_jwt_token_" + req.Email
	return &pb.AuthResponse{Token: token, Message: "Login successful"}, nil
}

// Test database helper
func setupTestDB(t *testing.T) (*sql.DB, *AuthServer) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		password TEXT
	);`

	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	return db, &AuthServer{db: db}
}

// Cleanup helper to delete test user
func cleanupUser(t *testing.T, db *sql.DB, email string) {
	_, err := db.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		t.Logf("Warning: Failed to cleanup user %s: %v", email, err)
	}
}

// Cleanup Redis OTP
func cleanupOTP(t *testing.T, email string) {
	err := config.Rdb.Del(context.Background(), "otp:"+email).Err()
	if err != nil {
		t.Logf("Warning: Failed to cleanup OTP for %s: %v", email, err)
	}
}

// Setup test environment
func setupTest(t *testing.T) (*sql.DB, *AuthServer) {
	// Try to load environment variables, but don't fail if .env doesn't exist
	if err := config.LoadEnv(); err != nil {
		t.Logf("Warning: Could not load .env file: %v", err)
		// In CI, environment variables should be set directly
	}

	// Try to initialize Redis, but don't fail if it's not available
	if config.Rdb == nil {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Redis not available, tests may fail: %v", r)
			}
		}()
		config.InitRedis()
	}

	return setupTestDB(t)
}