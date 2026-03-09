package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client/Server", func() {
	var (
		server *Server
		client *Client
		port   = ":9090"
	)

	Context("Happy Path", func() {
		BeforeEach(func() {
			server = NewServer(port)
			Expect(server.Start()).To(Succeed())
		})

		AfterEach(func() {
			if client != nil {
				_ = client.Close()
			}
			if server != nil {
				_ = server.Close()
			}
		})

		It("should be able to connect and receive an echo response", func() {
			client = NewClient()
			Expect(client.Connect("localhost" + port)).To(Succeed())

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			msg := "Hello Distributed World"
			resp, err := client.CallServer(ctx, msg)
			Expect(err).NotTo(HaveOccurred())
			// The server echoes the message with a newline appended
			Expect(resp).To(Equal(msg + "\n"))
		})
	})

	Context("Error Scenarios", func() {
		It("should fail to connect if server is not running", func() {
			client = NewClient()
			// Ensure we are trying to connect to a port where no server is running
			// We use a different port just to be safe, though the BeforeEach isn't running here.
			err := client.Connect("localhost:9999")
			Expect(err).To(HaveOccurred())
		})
	})
})
