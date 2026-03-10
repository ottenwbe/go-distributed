package main

import "net/rpc"

// BullyService is the RPC service for the bully algorithm.
type BullyService struct {
	node *Node
}

// ElectionArgs are the arguments for an election message.
type ElectionArgs struct {
	SenderID int
}

// ElectionReply is the reply to an election message.
type ElectionReply struct {
	IsAlive bool
}

// VictoryArgs are the arguments for a victory message.
type VictoryArgs struct {
	CoordinatorID int
}

// VictoryReply is the reply to a victory message.
type VictoryReply struct{}

// PingArgs are the arguments for a ping message.
type PingArgs struct{}

// PingReply is the reply to a ping message.
type PingReply struct{}

// HandleElection is called when a node receives an election message.
func (s *BullyService) HandleElection(args *ElectionArgs, reply *ElectionReply) error {
	reply.IsAlive = true
	// The receiver has a higher ID, so it tells the sender it's alive
	// and starts its own election.
	go s.node.HoldElection()
	return nil
}

// AnnounceVictory is called when a new coordinator is elected.
func (s *BullyService) AnnounceVictory(args *VictoryArgs, reply *VictoryReply) error {
	s.node.SetCoordinator(args.CoordinatorID)
	return nil
}

// Ping is used to check if a node is alive.
func (s *BullyService) Ping(args *PingArgs, reply *PingReply) error {
	return nil
}

func call(address string, method string, args interface{}, reply interface{}) error {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Call(method, args, reply)
}
