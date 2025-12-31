package main

import (
	"fmt"
	"goAuth/internal/api/handlers"
	"goAuth/internal/api/interceptors"
	"goAuth/pkg/utils"
	pb "goAuth/proto/gen"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
		return
	}
}

func main() {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AuthenticationInterceptor),
	)

	// Triggers every 2 minutes and cleans up all the expired tokens
	go utils.JwtStore.CleanUpExpiredTokens()

	pb.RegisterAuthServiceServer(s, &handlers.Server{})

	reflection.Register(s)

	port := os.Getenv("PORT")

	fmt.Printf("gRPC server running on port %s\n", port)

	// The TCP listener acts as a means for our gRPC server to communicate with the outside world
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening on the specified port: %v", err)
		return
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve", err)
		return
	}
}
