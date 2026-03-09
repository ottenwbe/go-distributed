package main

// Args holds arguments for the Arith.Multiply RPC call.
type Args struct {
	A, B int
}

// Reply holds the result of the Arith.Multiply RPC call.
type Reply struct {
	C int
}
