.PHONY: all build test clean run-gateway run-user docker-up docker-down

# Build all binaries
build:
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/user ./cmd/user
	go build -o bin/message ./cmd/message
	go build -o bin/relation ./cmd/relation
	go build -o bin/ai ./cmd/ai
	go build -o bin/file ./cmd/file

# Run all tests
test:
	go test ./internal/... -v -cover
	go test ./test/... -v

# Clean binaries
clean:
	rm -rf bin/

# Run individual services (for local dev)
run-gateway:
	go run ./cmd/gateway

run-user:
	go run ./cmd/user

run-message:
	go run ./cmd/message

run-relation:
	go run ./cmd/relation

run-ai:
	go run ./cmd/ai

run-file:
	go run ./cmd/file

# Docker
docker-build:
	docker build -f docker/Dockerfile.gateway -t aim-gateway .
	docker build -f docker/Dockerfile.user -t aim-user .
	docker build -f docker/Dockerfile.message -t aim-message .
	docker build -f docker/Dockerfile.relation -t aim-relation .
	docker build -f docker/Dockerfile.ai -t aim-ai .
	docker build -f docker/Dockerfile.file -t aim-file .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Tidy modules
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...
