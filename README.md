# go-distributed

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/ottenwbe/go-distributed.svg?branch=master)](https://travis-ci.org/ottenwbe/go-distributed)

In the age of cloud computing, micro services, and IoT, distributed systems are omnipresent. 
Oftentimes the heavy lifting that distributed systems do is hidden behind layers of layers of abstractions. 
To this end, we are going to lift the curtain and see how distributed systems tick. 

This is a repository for Go implementations of distributed paradigms and algorithms. 
All implementations are small demo apps that show the power of distributed computing.

## What paradigms/algorithms are implemented?

* client / server

## Develop

Get this repo:

```
git clone https://github.com/ottenwbe/go-distributed.git
```

### Structure

```
.
├── client-server   // client server demo where the server reverses all string messages
├── rpc             // rpc demo
├── go.mod          
├── go.sum
├── LICENSE         // MIT
├── Makefile        // build all projects with make
├── README.md       // this
└── vendor          // all vendored files
```
 
### Build

Use the makefile to build all demos

```
make build
```

or build just one of the deoms directly

```
go build -o clientServer client-server/*.go 
```

### Test

We use [ginkgo](https://github.com/onsi/ginkgo) for testing.

```
go get -u github.com/onsi/ginkgo/ginkgo  
go get -u github.com/onsi/gomega/...     

make test
```

## License

go-distributed is MIT-Licensed