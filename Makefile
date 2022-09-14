.PHONY: build test

build:
	go build -o conduit-connector-zendesk cmd/main.go

test:
	go test $(GOTEST_FLAGS) -count=1 -race ./...

lint:
	golangci-lint run -c .golangci.yml --go=1.18

dep:
	go mod download
	go mod tidy

