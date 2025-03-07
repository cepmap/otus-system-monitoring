DAEMON_BIN =./_dist/bin/daemon

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(DAEMON_BIN) -ldflags "$(LDFLAGS)" ./cmd/daemon

run: build
	$(DAEMON_BIN)
