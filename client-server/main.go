package main

import (
	"bufio"
	"fmt"
	"os"
)

const serverAddress = ":8888"

func main() {
	fmt.Println("Client Server Demo App")

	server := NewServer(serverAddress)
	err := server.Start()
	exitOnError(err, 1)
	defer closeServer(server)

	client := NewClient()
	err = client.Connect(serverAddress)
	exitOnError(err, 2)
	defer closeClient(client)

	// exchange messages with the ping server
	err = askInLoop("Message to server", func(msg string) {
		message, err := client.CallServer(msg)
		exitOnError(err, 3)
		fmt.Printf(">> Message from server: %v", message)
	})
	exitOnError(err, 4)

}

func closeServer(server *Server) {
	err := server.Close()
	exitOnError(err, 5)
}

func closeClient(client *Client) {
	err := client.Close()
	exitOnError(err, 6)
}

func askInLoop(question string, answerHandler func(answer string)) error {
	const stopIndicator = "quit"
	var (
		scanner = bufio.NewScanner(os.Stdin)
		run     = true
	)
	fmt.Printf(">> Type '%v' to stop the client (or ctrl+D)\n", stopIndicator)
	fmt.Printf(">> %v: ", question)
	for run && scanner.Scan() {
		msg := scanner.Text()
		if msg == stopIndicator {
			run = false
		} else {
			answerHandler(msg)
			fmt.Printf(">> %v: ", question)
		}
	}
	return scanner.Err()
}
