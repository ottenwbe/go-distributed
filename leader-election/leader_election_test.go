package main

import (
	"fmt"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Bully Leader Election", func() {
	var (
		nodes    []*Node
		peers    map[int]string
		numNodes int
	)

	BeforeEach(func() {
		// Suppress log output during tests to keep the output clean.
		log.SetOutput(io.Discard)

		numNodes = 3
		peers = make(map[int]string)
		for i := 1; i <= numNodes; i++ {
			// Use a high port range to avoid conflicts with other services.
			peers[i] = fmt.Sprintf("localhost:%d", 10000+i)
		}

		nodes = make([]*Node, numNodes)
		for i := 1; i <= numNodes; i++ {
			// Create a new node. The peers map includes the node itself,
			// which is fine as the node's logic filters itself out.
			nodes[i-1] = NewNode(i, peers[i], peers)
			Expect(nodes[i-1].Start()).To(Succeed())
		}
	})

	AfterEach(func() {
		for _, node := range nodes {
			if node != nil {
				node.Stop()
			}
		}
	})

	Context("when the node with the highest ID starts an election", func() {
		It("should become the coordinator immediately", func() {
			highestNode := nodes[numNodes-1] // Node with ID 3

			highestNode.HoldElection()

			// It should become the coordinator because it has no higher-ID peers.
			Eventually(func() int {
				highestNode.mutex.Lock()
				defer highestNode.mutex.Unlock()
				return highestNode.coordinator
			}, "2s", "100ms").Should(Equal(highestNode.id))

			// It should announce its victory to other nodes.
			for _, node := range nodes {
				Eventually(func() int {
					node.mutex.Lock()
					defer node.mutex.Unlock()
					return node.coordinator
				}, "3s", "100ms").Should(Equal(highestNode.id))
			}
		})
	})

	Context("when a lower ID node starts an election with all nodes up", func() {
		It("should result in the highest ID node becoming the coordinator", func() {
			lowestNode := nodes[0]           // Node with ID 1
			highestNode := nodes[numNodes-1] // Node with ID 3

			// Start an election from the lowest node.
			lowestNode.HoldElection()

			// The election process should result in the highest-ID node winning.
			// Eventually, all nodes should agree on the new coordinator.
			for _, node := range nodes {
				Eventually(func() int {
					node.mutex.Lock()
					defer node.mutex.Unlock()
					return node.coordinator
				}, "5s", "200ms").Should(Equal(highestNode.id))
			}
		})
	})

	Context("when the coordinator fails", func() {
		It("should elect the next highest ID node as the new coordinator", func() {
			// 1. Set the initial coordinator for all nodes.
			highestNode := nodes[numNodes-1] // Node 3
			for _, node := range nodes {
				node.SetCoordinator(highestNode.id)
			}
			Expect(nodes[0].coordinator).To(Equal(highestNode.id))
			Expect(nodes[1].coordinator).To(Equal(highestNode.id))

			// 2. Simulate coordinator failure.
			highestNode.Stop()
			nodes[numNodes-1] = nil // Prevent it from being stopped again in AfterEach.

			// 3. Wait for another node to detect the failure via heartbeat and start an election.
			// The heartbeat is 3s, so an election should be triggered shortly after.
			newHighestNode := nodes[numNodes-2] // Node 2

			// 4. Verify that the remaining nodes elect the new highest node as coordinator.
			for i := 0; i < numNodes-1; i++ {
				node := nodes[i]
				Eventually(func() int {
					node.mutex.Lock()
					defer node.mutex.Unlock()
					return node.coordinator
				}, "8s", "500ms").Should(Equal(newHighestNode.id))
			}
		})
	})
})
