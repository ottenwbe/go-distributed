package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"
)

// Client is a simple ping server
type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

// NewClient creates a non initialized client
func NewClient() *Client {
	return &Client{}
}

// CallServer sends a message to the server and returns the corresponding response
func (c *Client) CallServer(ctx context.Context, msg string) (string, error) {
	// Set deadline based on context
	if deadline, ok := ctx.Deadline(); ok {
		c.conn.SetDeadline(deadline)
		defer c.conn.SetDeadline(time.Time{})
	}

	// send a message
	_, err := fmt.Fprintf(c.conn, "%s\n", msg)
	if err != nil {
		return "", err
	}

	// receive and handle message
	message, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return message, err
}

// Connect to the server with the given address
func (c *Client) Connect(serverAddress string) (err error) {
	// connect to server
	c.conn, err = net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}
	c.reader = bufio.NewReader(c.conn)
	return nil
}

// Close shuts down the connection to the server
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
