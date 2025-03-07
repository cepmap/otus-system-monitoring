DAEMON_BIN =./_dist/bin/daemon

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(DAEMON_BIN) -ldflags "$(LDFLAGS)" ./cmd/daemon

run: build
	$(DAEMON_BIN)

.PHONY: proto
proto:
	protoc --plugin=protoc-gen-go=$(shell go env GOPATH)/bin/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(shell go env GOPATH)/bin/protoc-gen-go-grpc \
		--go_out=api/stats_service \
		--go_opt=paths=source_relative \
		--go-grpc_out=api/stats_service \
		--go-grpc_opt=paths=source_relative \
		api/stats_service/stats.proto
