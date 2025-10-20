package main

import (
	"fmt"
	"log"
	"net"

	"goAuth/config"
	pb "goAuth/proto"

	"google.golang.org/grpc"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}
	config.InitRedis()
}

func main() {
	db := InitDB()
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, &AuthServer{db: db})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	fmt.Println("Auth gRPC server running on port 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
