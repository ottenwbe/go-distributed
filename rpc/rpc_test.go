package main

import (
	"context"
	"net"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RPC Client/Server", func() {
	var (
		listener net.Listener
		client   *Client
		port     = ":9091"
		err      error
	)

	Context("with a running server", func() {
		BeforeEach(func() {
			listener, err = startServer(port)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if client != nil {
				_ = client.Close()
			}
			if listener != nil {
				_ = listener.Close()
			}
		})

		It("should successfully call the remote method", func() {
			client, err = NewClient("localhost"+port, 1*time.Second)
			Expect(err).NotTo(HaveOccurred())

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			a, b := 7, 8
			result, err := client.Multiply(ctx, a, b)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(a * b))
		})
	})

	Context("without a running server", func() {
		It("should fail to connect", func() {
			// Use a short timeout to make the test fail quickly.
			_, err := NewClient("localhost"+port, 100*time.Millisecond)
			Expect(err).To(HaveOccurred())
		})
	})
})
