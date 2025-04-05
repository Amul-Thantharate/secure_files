# Binary name
BINARY_NAME=secure-files

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test run fmt help

all: clean fmt build test

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

test:
	$(GOTEST) -v ./...

run: build
	./$(BINARY_NAME)

fmt:
	$(GOFMT) ./...

tidy:
	$(GOMOD) tidy

help:
	@echo "Make targets:"
	@echo "  build    - Build the application"
	@echo "  clean    - Clean build files"
	@echo "  test     - Run tests"
	@echo "  run      - Build and run the application"
	@echo "  fmt      - Format Go code"
	@echo "  tidy     - Tidy Go modules"
	@echo "  all      - Clean, format, build, and test"
	@echo "  help     - Show this help message"
