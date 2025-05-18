# Go parameters
BINARY_NAME=tech-challenge-fiap161
MAIN_PATH=./cmd/api
GO=go
GOLANGCI_LINT=golangci-lint
BUILD_DIR=/tmp/gobuild

.PHONY: all build clean test coverage lint tidy

all: tidy build test coverage lint

build:
	mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	$(GO) clean
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

test:
	$(GO) test -v ./...

coverage:
	$(GO) test -coverprofile=$(BUILD_DIR)/coverage.out $(shell go list ./... | grep -v ./cmd/api)
	$(GO) tool cover -func=$(BUILD_DIR)/coverage.out | grep total | awk '{print "Total coverage: " $$3}'

lint:
	$(GOLANGCI_LINT) run ./...

tidy:
	$(GO) mod tidy

# Run the application (builds in temp dir and runs from there)
run: build
	$(BUILD_DIR)/$(BINARY_NAME)