package main

import (
	"bufio"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

func main() {

	ln, _ := net.Listen("tcp", ":8081")

	log.Info("Ping server started and is listening on port :8081")

	conn, _ := ln.Accept()

	// run loop until it is interrupted
	for {
		// listen for new messages
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("Message Received from Client: %v", string(msg))

		// send new string back to client
		fmt.Fprintf(conn, msg+"\n")
	}
}
