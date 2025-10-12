package main

import (
	"net"
	"log"
	"fmt"

	"google.golang.org/grpc"
	"github.com/joho/godotenv"

	"goAuth/proto"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	db := InitDB()
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, &AuthServer{db: db})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	fmt.Println("Auth gRPC server running on port 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
