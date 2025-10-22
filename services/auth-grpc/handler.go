package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	
	"goAuth/shared/config"
	pb "goAuth/shared/proto"
	"goAuth/shared/utils"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

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

	subject := "Your Signup OTP"
	body := fmt.Sprintf("<p>Your OTP is <b>%s</b>. It expires in 5 minutes.</p>", otp)
	if err := utils.SendMail(req.Email, subject, body); err != nil {
		return nil, errors.New("Error sending OTP")
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
	config.Rdb.Del(ctx, "otp:" + req.Email)
	
	// Generate JWT token after successful verification
	token, _ := GenerateJWT(req.Email)
	
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
	token, _ := GenerateJWT(req.Email)
	return &pb.AuthResponse{Token: token, Message: "Login successful"}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	email, valid := ValidateJWT(req.Token)
	return &pb.ValidateResponse{Valid: valid, UserEmail: email}, nil
}