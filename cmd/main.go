package main

import (
	"log"
	"net"

	"github.com/sarthak0714/dbz/internal/database"
	pb "github.com/sarthak0714/dbz/pkg/api"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	db := database.NewDatabase()
	pb.RegisterDatabaseServer(s, db)

	// healthServer := health.NewServer()
	// to be implemented
	// pb.RegisterHealthServer(s, healthServer)

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
