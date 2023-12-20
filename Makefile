.PHONY: build test

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-zendesk.version=${VERSION}'" -o conduit-connector-zendesk cmd/connector/main.go

test:
	go test $(GOTEST_FLAGS) -count=1 -race ./...

lint:
	golangci-lint run -c .golangci.yml --go=1.18

dep:
	go mod download
	go mod tidy

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools.go
	@go list -e -f '{{ join .Imports "\n" }}' tools.go | xargs -tI % go install %
	@go mod tidy

