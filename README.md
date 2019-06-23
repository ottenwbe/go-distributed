# go-distributed

[![Build Status](https://travis-ci.org/ottenwbe/go-distributed.svg?branch=master)](https://travis-ci.org/ottenwbe/go-distributed)

Go implementations of distributed paradigms and algorithms.

## What paradigms/algorithms are implemented?

* client / server

## Develop

Get this repo 

```
git clone https://github.com/ottenwbe/go-distributed.git
```

### Structure

```
.
├── client-server   // client server demo
├── go.mod          
├── go.sum
├── LICENSE         // MIT
├── Makefile
├── README.md       // this
└── vendor          // all vendored files
```
 
# Build

```
go build -o clientServer client-server/*.go 
```

# Test

We use [ginkgo](https://github.com/onsi/ginkgo) for testing.

```
go get -u github.com/onsi/ginkgo/ginkgo  
go get -u github.com/onsi/gomega/...     
cd client-server
ginkgo
```
