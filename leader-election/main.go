package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	numNodes := flag.Int("n", 5, "Number of nodes to create")
	flag.Parse()

	if *numNodes < 1 {
		log.Fatal("Number of nodes must be at least 1")
	}

	peers := make(map[int]string)
	for i := 1; i <= *numNodes; i++ {
		peers[i] = fmt.Sprintf("localhost:%d", 9000+i)
	}

	nodes := make([]*Node, *numNodes)
	for i := 1; i <= *numNodes; i++ {
		nodes[i-1] = NewNode(i, peers[i], peers) // Each node gets the full peer list and filters it internally
	}

	for _, node := range nodes {
		if err := node.Start(); err != nil {
			log.Fatalf("Failed to start node %d: %v", node.id, err)
		}
	}

	// The node with the highest ID is the initial coordinator
	highestID := *numNodes
	for _, node := range nodes {
		node.SetCoordinator(highestID)
	}

	log.Infof("All %d nodes started. Initial coordinator is Node %d.", *numNodes, highestID)
	log.Info("Press Ctrl+C to exit.")
	log.Info("Simulating coordinator failure in 10 seconds...")

	// Simulate coordinator failure and restart
	go func() {
		time.Sleep(10 * time.Second)
		coordinatorNode := nodes[highestID-1]
		log.Warnf(">>> Stopping coordinator Node %d to trigger election <<<", coordinatorNode.id)
		coordinatorNode.Stop()

		// Wait for a minute and then bring the old coordinator back online.
		log.Info("Waiting 0.5 minutes before bringing the old coordinator back...")
		time.Sleep(30 * time.Second)

		log.Warnf(">>> Bringing Node %d back online to trigger another election <<<", highestID)
		// We need to create a new node object to simulate a restart.
		restartedNode := NewNode(highestID, peers[highestID], peers)
		if err := restartedNode.Start(); err != nil {
			log.Errorf("Failed to restart node %d: %v", highestID, err)
			return
		}
		// The restarted node will bully its way to become the coordinator.
		go restartedNode.HoldElection()
		nodes[highestID-1] = restartedNode // Replace the old node with the new one for cleanup.
	}()

	// Wait for termination signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutting down all nodes...")
	for _, node := range nodes {
		node.Stop()
	}
	log.Info("Shutdown complete.")
}
