BIN 	= $(GOPATH)/bin
GO      = go
GOTEST  = $(BIN)/ginkgo
GOFMT   = gofmt
GOLINT  = golint
GOVET  	= $(GO) vet
DATE    = $(shell date +'%Y.%m.%d-%H:%M:%S')
VERSION = $(shell git describe --tags --always 2> /dev/null)
M 		= $(shell printf "\e[35;1m>\e[0m")

.PHONY: list
list: ; $(info $(M) listing all distributed projects…)
	@echo All distributed projects:; \
	for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		echo $$d ; \
	 done

.PHONY: test
test: ; $(info shell printf $(M) running tests..) @ ## Run tests
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOTEST) -cover $$d ; \
	 done

.PHONY: lint
lint: ; $(info $(M) running golint…) @ ## Run golint
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOLINT) $${d};  \
	done

.PHONY: vet
vet: ; $(info $(M) running vet…) @ ## Run go vet
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOVET) $${d};  \
	done

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go  ; \
	 done

.PHONY: build
build: ; $(info $(M) building release executables…) @ ## Build the apps' binary version
	for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GO) build \
		$$d/*.go ; \
	done

.PHONY: release
release: ; $(info $(M) building release executables…) @ ## Build the apps' binary release version
	@for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GO) build \
		-tags release \
		-ldflags "-s -w" \
		$$d/*.go ; \
	 done

.PHONY: version
version:
	@echo $(VERSION)	 
		
.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[35m%-15s\033[0m %s\n", $$1, $$2}'
