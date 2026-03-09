package main

import (
	"net"

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
			client, err = NewClient("localhost" + port)
			Expect(err).NotTo(HaveOccurred())

			a, b := 7, 8
			result, err := client.Multiply(a, b)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(a * b))
		})
	})

	Context("without a running server", func() {
		It("should fail to connect", func() {
			_, err := NewClient("localhost" + port)
			Expect(err).To(HaveOccurred())
		})
	})
})
