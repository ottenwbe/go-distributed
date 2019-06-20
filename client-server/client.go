package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {

	// get server address
	serverAddress := determineServerAddress()

	// connect to server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.WithError(err).Error("Error while connecting to server.")
		os.Exit(1)
	}

	// send and receive messages in a loop
	err = messageLoop(conn)
	if err != nil {
		log.WithError(err).Error("Error in Message Loop.")
		os.Exit(2)
	}
}

func determineServerAddress() string {
	const defaultAddress = ":8081"
	serverAddress := readAddress(defaultAddress)
	serverAddress = checkForDefaultAddress(serverAddress, defaultAddress)
	return serverAddress
}

func readAddress(defaultAddress string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter the server's address [default: '%v']: ", defaultAddress)
	serverAddress, err := reader.ReadString('\n')
	if err != nil {
		log.WithError(err).Warn("Could not read server address from input")
		serverAddress = ""
	}
	return serverAddress
}

func checkForDefaultAddress(serverAddress string, defaultAddress string) string {
	if serverAddress == "" {
		return defaultAddress
	}
	return serverAddress
}

func messageLoop(conn net.Conn) error {
	var (
		err     error
		scanner = bufio.NewScanner(os.Stdin)
	)
	for err == nil && scanner.Scan() {
		err = sendMessage(scanner, conn)
		if err != nil {
			err = waitForReplyMessage(conn)
		}
	}
	if err != nil {
		err = scanner.Err()
	}
	return err
}

func sendMessage(scanner *bufio.Scanner, conn net.Conn) error {
	fmt.Print("Text to send: ")
	msg := scanner.Text()
	_, err := fmt.Fprintf(conn, msg+"\n")
	return err
}

func waitForReplyMessage(conn net.Conn) error {
	message, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("Message from server: %v", message)
	return err
}
