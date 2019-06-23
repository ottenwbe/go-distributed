package main

import (
	"bufio"
	"fmt"
	"net"
)

// Client is a simple ping server
type Client struct {
	conn net.Conn
}

// NewClient creates a non initialized client
func NewClient() *Client {
	return &Client{}
}

// CallServer sends a message to the server and returns the corresponding response
func (c *Client) CallServer(msg string) (string, error) {
	// send a message
	_, err := fmt.Fprintf(c.conn, msg+"\n")
	if err != nil {
		return "", err
	}

	// receive and handle message
	message, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return message, err
}

// Connect to the server with the given address
func (c *Client) Connect(serverAddress string) (err error) {
	// connect to server
	c.conn, err = net.Dial("tcp", serverAddress)
	return err
}

// Close shuts down the connection to the server
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
