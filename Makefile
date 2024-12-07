.PHONY: build run test lint migrate

# Build variables
BINARY_NAME=gohex
MAIN_PATH=cmd/api/main.go
CONFIG_PATH=config/config.yaml

# Build the application
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	go run $(MAIN_PATH) -config $(CONFIG_PATH)

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Run database migrations
migrate:
	go run cmd/migrate/main.go -config $(CONFIG_PATH)

# Clean build artifacts
clean:
	rm -rf bin/ 