package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Server for a simple ping application.
type Server struct {
	listenAddress string
	listener      net.Listener
	quit          chan struct{}
	wg            sync.WaitGroup
}

// NewServer creates a simple server object and initializes the address the server should listen on
func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
		quit:          make(chan struct{}),
	}
}

// Start the server
func (s *Server) Start() error {
	var err error

	s.listener, err = net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	s.wg.Go(s.serve)

	log.Infof("Ping server started and is listening on port %v", s.listenAddress)

	return nil
}

func (s *Server) serve() {

	for {
		select {
		case <-s.quit:
			log.Info("Server shutting down...")
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			// Check if the error is due to the listener being closed.
			// If so, it's a clean shutdown.
			select {
			case <-s.quit:
				return
			default:
				log.WithError(err).Error("Failed to accept connection")
			}
			continue
		}

		s.wg.Go(func() {
			s.handleConnection(conn)
		})
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// Connection closed or other error
			return
		}
		// send new string back to the client
		fmt.Fprintf(conn, "%s", msg)
	}
}

// Close shuts down the server
func (s *Server) Close() error {
	select {
	case <-s.quit:
		return nil
	default:
		close(s.quit)
	}
	if s.listener != nil {
		err := s.listener.Close()
		s.wg.Wait()
		return err
	}
	return nil
}
