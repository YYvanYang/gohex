.PHONY: build run test lint migrate

# Build variables
BINARY_NAME=gohex
MAIN_PATH=cmd/api/main.go
CONFIG_PATH=config/config.yaml
DOCKER_IMAGE=gohex
DOCKER_TAG=latest

# Build the application
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Build docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run docker development environment
docker-dev:
	docker-compose up --build

# Run docker production environment
docker-prod:
	docker-compose -f docker-compose.prod.yml up --build -d

# Stop docker containers
docker-down:
	docker-compose down -v

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