DAEMON_BIN =./_dist/bin/daemon
CLIENT_BIN =./_dist/bin/client
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-server:
	go build -v -o $(DAEMON_BIN) -ldflags "$(LDFLAGS)" ./cmd/daemon

build-client:
	go build -v -o $(CLIENT_BIN) -ldflags "$(LDFLAGS)" ./cmd/client

build: build-server build-client

run-server: build-server
	$(DAEMON_BIN)

run-client: build-client
	$(CLIENT_BIN)	

run: run-server run-client

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.6

lint: install-lint-deps
	golangci-lint run --timeout=90s ./...

unit-tests:
	go test -race -count 100 ./internal/... -v

integration-tests:
	go test  -count 1 ./tests/integration/... -v

tests: unit-tests integration-tests