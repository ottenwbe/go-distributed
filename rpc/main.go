package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mode := flag.String("mode", "demo", "Mode to run: 'server', 'client', or 'demo'")
	port := flag.String("port", ":8081", "Port to listen/connect to")
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
	listener, err := startServer(port)
	if err != nil {
		log.Fatalf("Failed to start RPC server: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	listener.Close()
}

func runClient(port string) {
	client, err := NewClient("localhost" + port)
	if err != nil {
		log.Fatalf("Failed to connect to RPC server: %v", err)
	}
	defer client.Close()

	a, b := 7, 8
	log.Printf("Client calling Arith.Multiply with %d and %d", a, b)

	result, err := client.Multiply(a, b)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}

	log.Printf("Client received result: %d * %d = %d", a, b, result)
}

func runDemo(port string) {
	listener, err := startServer(port)
	if err != nil {
		log.Fatalf("Failed to start RPC server: %v", err)
	}
	defer listener.Close()

	time.Sleep(100 * time.Millisecond) // Give server time to start
	runClient(port)
}
