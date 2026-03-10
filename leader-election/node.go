package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Node represents a process in the distributed system.
type Node struct {
	id          int
	address     string
	peers       map[int]string // Map of node ID to address
	coordinator int
	isElection  bool
	listener    net.Listener
	server      *rpc.Server
	wg          sync.WaitGroup
	quit        chan struct{}
	mutex       sync.Mutex
}

// NewNode creates a new node.
func NewNode(id int, address string, peers map[int]string) *Node {
	return &Node{
		id:      id,
		address: address,
		peers:   peers,
		quit:    make(chan struct{}),
	}
}

// Start begins the node's RPC server and heartbeating.
func (n *Node) Start() error {
	n.server = rpc.NewServer()
	bullyService := &BullyService{node: n}
	if err := n.server.Register(bullyService); err != nil {
		return fmt.Errorf("failed to register RPC service: %w", err)
	}

	var err error
	n.listener, err = net.Listen("tcp", n.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", n.address, err)
	}

	log.Infof("Node %d listening on %s", n.id, n.address)

	n.wg.Go(func() {
		n.server.Accept(n.listener)
	})

	n.wg.Go(func() {
		n.heartbeat()
	})

	return nil
}

// Stop shuts down the node.
func (n *Node) Stop() {
	close(n.quit)
	if n.listener != nil {
		n.listener.Close()
	}
	n.wg.Wait()
	log.Infof("Node %d shut down", n.id)
}

// HoldElection starts an election process.
func (n *Node) HoldElection() {
	n.mutex.Lock()
	if n.isElection {
		n.mutex.Unlock()
		return
	}
	n.isElection = true
	n.mutex.Unlock()

	log.Infof("Node %d is holding an election", n.id)

	responses := make(chan bool, 1)

	higherNodes := 0
	for peerID, peerAddress := range n.peers {
		if peerID > n.id {
			higherNodes++
			go func(addr string) {
				log.Infof("Node %d sending ELECTION to %s", n.id, addr)
				err := call(addr, "BullyService.HandleElection", &ElectionArgs{SenderID: n.id}, &ElectionReply{})
				if err == nil {
					// Received a response from a higher node
					responses <- true
				} else {
					log.Warnf("Node %d failed to send election message to %s: %v", n.id, addr, err)
				} // If there's an error, we assume the higher node is down.
			}(peerAddress)
		}
	}

	// If there are no higher nodes, this node wins the election.
	if higherNodes == 0 {
		n.becomeCoordinator()
		return
	}

	// Wait for a response from a higher node or a timeout.
	select {
	case <-responses:
		// A higher node has responded and will take over the election.
		// This node's part in the election is over. It will now wait for a VICTORY message.
		log.Infof("Node %d yielding election to a higher node.", n.id)
		// The isElection flag will be reset when a VICTORY message is received via SetCoordinator.
		return
	case <-time.After(2 * time.Second):
		// No higher node responded within the timeout period. This node wins.
		n.becomeCoordinator()
	}
}

func (n *Node) becomeCoordinator() {
	log.Infof("Node %d is the new coordinator", n.id)
	n.SetCoordinator(n.id)

	for peerID, peerAddress := range n.peers {
		if peerID != n.id {
			go func(id int, addr string) {
				log.Infof("Node %d announcing victory to node %d", n.id, id)
				err := call(addr, "BullyService.AnnounceVictory", &VictoryArgs{CoordinatorID: n.id}, &VictoryReply{})
				if err != nil {
					log.Warnf("Node %d failed to announce victory to node %d: %v", n.id, id, err)
				}
			}(peerID, peerAddress)
		}
	}

	n.mutex.Lock()
	n.isElection = false
	n.mutex.Unlock()
}

// SetCoordinator sets the coordinator for the node.
func (n *Node) SetCoordinator(coordinatorID int) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if n.coordinator != coordinatorID {
		n.coordinator = coordinatorID
		log.Infof("Node %d recognizes Node %d as the coordinator", n.id, coordinatorID)
	}
	n.isElection = false
}

func (n *Node) heartbeat() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-n.quit:
			return
		case <-ticker.C:
			n.mutex.Lock()
			coordID := n.coordinator
			isElection := n.isElection
			n.mutex.Unlock()

			if coordID == 0 || isElection {
				continue
			}

			if n.id == coordID {
				continue // I am the coordinator
			}

			coordAddress, ok := n.peers[coordID]
			if !ok {
				log.Errorf("Node %d: Coordinator %d not in peer list", n.id, coordID)
				continue
			}

			err := call(coordAddress, "BullyService.Ping", &PingArgs{}, &PingReply{})
			if err != nil {
				log.Warnf("Node %d: Coordinator %d is down, starting election.", n.id, coordID)
				n.HoldElection()
			} else {
				log.Infof("Node %d: Ping to coordinator %d successful", n.id, coordID)
			}
		}
	}
}
