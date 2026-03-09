package main

import (
	"context"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Client is a wrapper for an RPC client.
type Client struct {
	client    *rpc.Client
	closeOnce sync.Once
	closeErr  error
}

// NewClient creates a new RPC client.
func NewClient(serverAddress string, timeout time.Duration) (*Client, error) {
	conn, err := net.DialTimeout("tcp", serverAddress, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{client: rpc.NewClient(conn)}, nil
}

// Multiply calls the remote Arith.Multiply method.
// It respects the context for cancellation or deadlines.
func (c *Client) Multiply(ctx context.Context, a, b int) (int, error) {
	args := &Args{A: a, B: b}
	var reply Reply

	call := c.client.Go("Arith.Multiply", args, &reply, nil)

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case replyCall := <-call.Done:
		// RPC call completed.
		if replyCall.Error != nil {
			return 0, replyCall.Error
		}
	}
	return reply.C, nil
}

// Close closes the connection to the server.
func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		if c.client != nil {
			c.closeErr = c.client.Close()
		}
	})
	return c.closeErr
}
