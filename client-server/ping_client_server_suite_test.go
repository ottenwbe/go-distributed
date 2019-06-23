package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClientServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ping Server And Client Suite")
}
