package handlers

import (
	"context"
	"goAuth/internal/repositories/mongodb"
	"goAuth/pkg/utils"
	pb "goAuth/proto/gen"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := mongodb.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = utils.VerifyPassword(req.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Incorrect username or password")
	}

	tokenString, err := utils.SignToken(user.Id, user.Username, user.Role)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Could not create token")
	}

	return &pb.LoginResponse{
		Status: true,
		Token:  tokenString,
	}, nil
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.LoginResponse, error) {
	user, err := mongodb.AddUserToDB(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokenString, err := utils.SignToken(user.Id, user.Username, user.Role)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Could not create token")
	}

	return &pb.LoginResponse{
		Status: true,
		Token:  tokenString,
	}, nil
}

func (s *Server) ChangeRole(ctx context.Context, req *pb.ChangeRoleRequest) (*pb.ChangeRoleResponse, error) {
	// Allow only super admin to be able to change user roles
	err := utils.AuthorizeUser(ctx, "super_admin")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userId := req.GetId()
	updatedRole := req.GetRole()
	err = mongodb.ModifyUserRoleInDB(ctx, userId, updatedRole)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ChangeRoleResponse{
		Status: true,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.EmptyRequest) (*pb.LogoutResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "No metadata found")
	}

	val, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized access")
	}

	token := strings.TrimPrefix(val[0], "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized access")
	}

	// Get expiry time from context (set by authentication interceptor)
	expiryTimeStamp := ctx.Value(utils.ContextKey("expiresAt"))
	expiryTimeInt, ok := expiryTimeStamp.(int64)
	if !ok {
		return nil, status.Error(codes.Internal, "Failed to retrieve token expiry time")
	}

	expiryTime := time.Unix(expiryTimeInt, 0)

	utils.JwtStore.AddToken(token, expiryTime)

	return &pb.LogoutResponse{
		Status: true,
	}, nil
}
