package main

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"goAuth/proto"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	db *sql.DB
}

func (s *AuthServer) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.AuthResponse, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	_, err := s.db.Exec("INSERT INTO users(email, password) VALUES(?, ?)", req.Email, string(hash))
	if err != nil {
		return nil, errors.New("Email already exists!")
	}
	token, _ := GenerateJWT(req.Email)
	return &proto.AuthResponse{Token: token, Message: "Signup successful"}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
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
	return &proto.AuthResponse{Token: token, Message: "Login successful"}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	email, valid := ValidateJWT(req.Token)
	return &proto.ValidateResponse{Valid: valid, UserEmail: email}, nil
}