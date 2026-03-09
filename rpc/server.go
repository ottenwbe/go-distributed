package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"
)

var registerOnce sync.Once

// Arith is a type that we will register as an RPC service.
type Arith int

// Multiply is the remote method we will expose.
func (t *Arith) Multiply(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func startServer(port string) (net.Listener, error) {
	var regErr error
	registerOnce.Do(func() {
		arith := new(Arith)
		regErr = rpc.Register(arith)
	})
	if regErr != nil {
		return nil, regErr
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	log.Printf("RPC server listening on %s", listener.Addr())
	go rpc.Accept(listener)
	return listener, nil
}
