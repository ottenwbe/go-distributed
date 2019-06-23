package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

// Server for a simple ping application. NOTE: The current implementation avoids synchronization of the go function and close for simplicity.
type Server struct {
	listenAddress string
	listener      net.Listener
	conn          net.Conn
	closed        bool
}

// NewServer creates a simple server object and initializes the address the server should listen on
func NewServer(listenAddress string) *Server {
	return &Server{
		closed:        true,
		listenAddress: listenAddress,
	}
}

// Start the server
func (s *Server) Start() error {
	var err error

	s.closed = false

	s.listener, err = net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	go s.run()

	log.Infof("Ping server started and is listening on port %v", s.listenAddress)

	return nil
}

func (s *Server) run() {
	var err error
	s.conn, err = s.listener.Accept()
	if err == nil {
		reader := bufio.NewReader(s.conn)
		for !s.closed && err == nil {
			var msg string
			// listen for new messages
			msg, err = reader.ReadString('\n')
			err = s.handleMessage(msg)
		}
	}

	log.WithError(err).Info("Server closed down")

}

func (s *Server) handleMessage(msg string) (err error) {
	if s.conn != nil {
		// send new string back to the client
		_, err = fmt.Fprintf(s.conn, msg+"\n")
	}
	return err
}

// Close shuts down the server
func (s *Server) Close() error {
	var err error

	if !s.closed {
		s.closed = true
		if s.conn != nil {
			err = s.conn.Close()
		}
		err = s.listener.Close()
	}

	return err
}
