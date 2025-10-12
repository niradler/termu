.PHONY: build run install test clean fmt lint

BINARY_NAME=olloco
MAIN_PATH=cmd/olloco/main.go

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH) chat

install:
	go install $(MAIN_PATH)

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run

deps:
	go mod download
	go mod tidy

chat: build
	./$(BINARY_NAME) chat

sandbox: build
	./$(BINARY_NAME) --sandbox chat

yolo: build
	@echo "⚠️  ⚠️  ⚠️  WARNING: YOLO MODE ⚠️  ⚠️  ⚠️"
	@echo "All security checks will be DISABLED!"
	./$(BINARY_NAME) --yolo chat

install-tools: build
	./$(BINARY_NAME) install-tools

help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  run           - Run in chat mode"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  test          - Run tests"
	@echo "  clean         - Remove binary and clean"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  deps          - Download dependencies"
	@echo "  chat          - Build and run chat mode"
	@echo "  sandbox       - Build and run in sandbox mode"
	@echo "  yolo          - ⚠️  Build and run in YOLO mode (NO SECURITY!)"
	@echo "  install-tools - Install modern CLI tools (fd, rg, bat, etc.)"



