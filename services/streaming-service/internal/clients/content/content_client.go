package content

import (
	"context"

	"google.golang.org/grpc"
	"github.com/streamverse/proto/gen/go/content"
)

// Client is a gRPC client for the content service.
type Client struct {
	conn    *grpc.ClientConn
	client  content.ContentServiceClient
}

// NewClient creates a new content service client.
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := content.NewContentServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// GetContent retrieves content by ID.
func (c *Client) GetContent(ctx context.Context, id string) (*content.GetContentResponse, error) {
	return c.client.GetContent(ctx, &content.GetContentRequest{ContentId: id})
}
