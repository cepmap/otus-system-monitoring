DAEMON_BIN =./_dist/bin/daemon

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(DAEMON_BIN) -ldflags "$(LDFLAGS)" ./cmd/daemon

run: build
	$(DAEMON_BIN)


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.6

lint: install-lint-deps
	golangci-lint run --timeout=90s ./...

test:
	go test ./... -race -count 2
