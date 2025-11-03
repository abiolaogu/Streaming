package payment

import (
	"context"

	"google.golang.org/grpc"
	"github.com/streamverse/proto/gen/go/payment"
)

// Client is a gRPC client for the payment service.
type Client struct {
	conn    *grpc.ClientConn
	client  payment.PaymentServiceClient
}

// NewClient creates a new payment service client.
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// GetSubscription retrieves subscription status for a user.
func (c *Client) GetSubscription(ctx context.Context, userID string) (*payment.GetSubscriptionResponse, error) {
	return c.client.GetSubscription(ctx, &payment.GetSubscriptionRequest{UserId: userID})
}

// CheckConcurrentStreams checks the concurrent stream limit for a user.
func (c *Client) CheckConcurrentStreams(ctx context.Context, userID string) (*payment.CheckConcurrentStreamsResponse, error) {
	return c.client.CheckConcurrentStreams(ctx, &payment.CheckConcurrentStreamsRequest{UserId: userID})
}
