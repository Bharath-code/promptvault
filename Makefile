BINARY=promptvault
VERSION=0.1.0
BUILD_DIR=./dist
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"
GO_TAGS=-tags "sqlite_fts5"

.PHONY: all build install clean deps run test seed init

all: deps build

## Install dependencies
deps:
	go mod tidy
	go mod download

## Build for current platform
build:
	CGO_ENABLED=1 go build $(GO_TAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) .

## Build for all platforms
build-all:
	CGO_ENABLED=1 GOOS=darwin  GOARCH=amd64  go build $(GO_TAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 .
	CGO_ENABLED=1 GOOS=darwin  GOARCH=arm64  go build $(GO_TAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 .
	CGO_ENABLED=1 GOOS=linux   GOARCH=amd64  go build $(GO_TAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .
	CGO_ENABLED=1 GOOS=linux   GOARCH=arm64  go build $(GO_TAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64 .
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64  go build $(GO_TAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .

## Install to $GOPATH/bin
install:
	CGO_ENABLED=1 go install $(GO_TAGS) $(LDFLAGS) .

## Run the TUI directly
run:
	CGO_ENABLED=1 go run $(GO_TAGS) .

## Run tests
test:
	CGO_ENABLED=1 go test -tags "fts5" ./...

## Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

## Initialize vault with curated prompts
init:
	go run . init

## Seed with example prompts (alias for init --force)
seed:
	go run . init --force

## Show help
help:
	@echo "PromptVault Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## /  /'
