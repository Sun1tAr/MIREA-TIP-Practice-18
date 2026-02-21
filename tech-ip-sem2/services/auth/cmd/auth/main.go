package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	grp "github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/auth/internal/grpc"
	pb "github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/proto/auth"
)

func main() {
	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &grp.Server{})

	go func() {
		log.Printf("Auth gRPC server starting on :%s", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Auth gRPC server...")
	s.GracefulStop()
}
