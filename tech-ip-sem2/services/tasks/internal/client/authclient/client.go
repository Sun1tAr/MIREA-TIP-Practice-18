package authclient

import (
	"context"
	"fmt"
	"time"

	pb "github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn    *grpc.ClientConn
	client  pb.AuthServiceClient
	timeout time.Duration
}

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &Client{
		conn:    conn,
		client:  pb.NewAuthServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) VerifyToken(ctx context.Context, token string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Verify(ctx, &pb.VerifyRequest{Token: token})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return false, "", fmt.Errorf("auth service unavailable: %w", err)
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return false, "", nil
		case codes.DeadlineExceeded:
			return false, "", fmt.Errorf("auth service timeout")
		default:
			return false, "", fmt.Errorf("auth service error: %v", st.Message())
		}
	}

	return resp.Valid, resp.Subject, nil
}
