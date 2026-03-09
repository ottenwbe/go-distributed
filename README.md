# go-distributed

[![CI](https://github.com/ottenwbe/go-distributed/actions/workflows/ci.yml/badge.svg)](https://github.com/ottenwbe/go-distributed/actions/workflows/ci.yml)

This repository provides demo and examples of Go implementations of distributed paradigms and algorithms.

## What paradigms/algorithms are implemented?

* client / server
* rpc (Remote Procedure Call)

## Develop

Get this repo 

```bash
git clone https://github.com/ottenwbe/go-distributed.git
```

### Structure

```text
.
├── client-server   # Client/Server pattern implementation
├── rpc             # RPC pattern implementation
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```
 
 
# How to Run

The `client-server` example can be run in different modes using the `-mode` flag.

### Demo Mode (Default)

This runs both the client and server in the same process for a quick demonstration.

```bash
go run ./client-server
```

### Server & Client Mode

To demonstrate the distributed nature, run the server and client in separate terminals.

```bash
# Terminal 1: Run the server
go run ./client-server -mode server

# Terminal 2: Run the client
go run ./client-server -mode client
```

# Test

We use ginkgo for testing.


```bash
go test ./...
```

or 

```bash
ginkgo  -v ./...
```