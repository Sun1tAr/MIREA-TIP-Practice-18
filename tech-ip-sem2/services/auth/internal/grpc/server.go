package grpc

import (
	"context"
	"log"

	"github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/auth/internal/service"
	pb "github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	log.Printf("gRPC Verify called with token: %s", req.Token)

	valid, subject := service.VerifyToken(req.Token)
	if !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &pb.VerifyResponse{
		Valid:   true,
		Subject: subject,
	}, nil
}
