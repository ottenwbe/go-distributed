BIN 	= $(GOPATH)/bin
GO      = go
GOTEST  = $(BIN)/ginkgo
GOFMT   = gofmt
GOLINT  = $(BIN)/golint
GOVET  	= $(GO) vet
DATE    = $(shell date +'%Y.%m.%d-%H:%M:%S')
VERSION = $(shell git describe --tags --always 2> /dev/null)
MARKER 	= $(shell printf "\e[35;1m>\e[0m")

.PHONY: list
list: ; $(info $(MARKER) listing all go-distributed demos...) @ ## List all go-distributed demos
	@echo All distributed projects:; \
	for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		echo $$(basename $$d); \
	 done

.PHONY: test
test: ; $(info shell printf $(MARKER) running tests..) @ ## Run tests for all demos
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOTEST) -cover $$d; \
	 done

.PHONY: lint
lint: ; $(info $(MARKER) running golint…) @ ## Run golint over all demo sources
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOLINT) $${d}; \
	done

.PHONY: vet
vet: ; $(info $(MARKER) running vet…) @ ## Run go vet for all demo sources
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOVET) $${d}/*.go;  \
	done

.PHONY: fmt
fmt: ; $(info $(MARKER) running gofmt…) @ ## Run gofmt on all demo source files
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go; \
	 done

.PHONY: build
build: ; $(info $(MARKER) building executables…) @ ## Build the apps' binary version
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GO) build \
		-o $$(basename $$d)-$(VERSION)-$(DATE) \
		$$d/*.go; \
	done

.PHONY: release
release: ; $(info $(MARKER) building release executables…) @ ## Build the apps' binary release version
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GO) build \
		-tags release \
		-ldflags "-s -w" \
		-o $$(basename $$d)-$(VERSION) \
		$$d/*.go; \
	 done

.PHONY: version
version: ; $(info $(MARKER) getting the current version…) @ ## Print the current version of go-distributed
	@echo $(VERSION)	 
		
.PHONY: help
help: ; @ ## Show this help menu
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[35m%-15s\033[0m %s\n", $$1, $$2}'
