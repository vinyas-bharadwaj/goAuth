package handlers

import (
	"context"
	"goAuth/internal/repositories/mongodb"
	"goAuth/pkg/utils"
	pb "goAuth/proto/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GoogleLogin handles Google OAuth authentication
func (s *Server) GoogleLogin(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.GoogleLoginResponse, error) {
	// Verify the Google ID token
	googleUser, err := utils.VerifyGoogleIDToken(ctx, req.GetIdToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid Google ID token")
	}

	// Check if user exists by Google ID first
	existingUser, err := mongodb.GetUserByGoogleId(ctx, googleUser.Sub)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error checking for existing user")
	}

	var user *pb.User

	// If user doesn't exist with this Google ID, check by email
	if existingUser == nil {
		existingUser, err = mongodb.GetUserByEmail(ctx, googleUser.Email)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error checking for existing user by email")
		}

		// If user exists by email but not Google ID, update their Google ID
		if existingUser != nil {
			// Update existing user with Google ID and picture
			err = mongodb.UpdateUserGoogleInfo(ctx, existingUser.Id, googleUser.Sub, googleUser.Picture)
			if err != nil {
				return nil, status.Error(codes.Internal, "Error updating user with Google information")
			}

			// Update the model with the new Google info for the response
			existingUser.GoogleId = googleUser.Sub
			existingUser.Picture = googleUser.Picture
			user = mongodb.MapModelUserToPbUser(existingUser)
		} else {
			// Create a new user
			newUser, err := mongodb.CreateGoogleUser(ctx, googleUser.Email, googleUser.Name, googleUser.Sub, googleUser.Picture)
			if err != nil {
				return nil, status.Error(codes.Internal, "Error creating new user")
			}
			user = mongodb.MapModelUserToPbUser(newUser)
		}
	} else {
		// User exists, use their data
		user = mongodb.MapModelUserToPbUser(existingUser)
	}

	// Generate access token (JWT)
	accessToken, err := utils.SignToken(user.Id, user.Username, user.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not create access token")
	}

	// Generate refresh token (for now, using the same token generation)
	// In a production environment, you'd want a separate refresh token mechanism
	refreshToken, err := utils.SignToken(user.Id, user.Username, user.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not create refresh token")
	}

	return &pb.GoogleLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
