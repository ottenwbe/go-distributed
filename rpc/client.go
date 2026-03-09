package main

import (
	"net/rpc"
	"sync"
)

// Client is a wrapper for an RPC client.
type Client struct {
	client    *rpc.Client
	closeOnce sync.Once
	closeErr  error
}

// NewClient creates a new RPC client.
func NewClient(serverAddress string) (*Client, error) {
	c, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}

// Multiply calls the remote Arith.Multiply method.
func (c *Client) Multiply(a, b int) (int, error) {
	args := &Args{A: a, B: b}
	var reply Reply
	err := c.client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		return 0, err
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
