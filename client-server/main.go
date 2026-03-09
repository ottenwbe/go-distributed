package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	mode := flag.String("mode", "demo", "Mode to run: 'server', 'client', or 'demo'")
	port := flag.String("port", ":8080", "Port to listen/connect to")
	flag.Parse()

	switch *mode {
	case "server":
		runServer(*port)
	case "client":
		runClient(*port)
	case "demo":
		runDemo(*port)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

func runServer(port string) {
	server := NewServer(port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer server.Close()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
}

func runClient(port string) {
	client := NewClient()
	if err := client.Connect("localhost" + port); err != nil {
		log.Fatalf("Failed to connect client: %v", err)
	}
	defer client.Close()

	msg := "Hello Distributed World"
	log.Infof("Client sending: %q", msg)

	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.CallServer(ctx, msg)
	if err != nil {
		log.Fatalf("CallServer failed: %v", err)
	}

	log.Infof("Client received: %q", response)
}

func runDemo(port string) {
	// 1. Start the Server
	server := NewServer(port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer server.Close()

	time.Sleep(100 * time.Millisecond)

	// 2. Create the Client
	client := NewClient()
	if err := client.Connect("localhost" + port); err != nil {
		log.Fatalf("Failed to connect client: %v", err)
	}
	defer client.Close()

	// 3. Send a message
	msg := "Hello Distributed World"
	log.Infof("Client sending: %q", msg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.CallServer(ctx, msg)
	if err != nil {
		log.Fatalf("CallServer failed: %v", err)
	}

	log.Infof("Client received: %q", response)
}
