# This Makefile provides convenient methods for common operations

.PHONY: add commit push build-docs run run-worker test fmt vet clean docker-build docker-run

# The default commit message
MESSAGE = "updated: `date '+%Y-%m-%d %H:%M:%S'`"

# Development commands
dev: build-docs run

# Add all changes to the staging area
add:
	git add .

# Commit changes with the provided message
commit:
	git commit -m $(MESSAGE)

# Push changes to the master branch
push:
	git push origin master

# Run add, commit, and push in one command
all: add commit push

# Run HTTP API server
run:
	@go run main.go

# Run background worker
run-worker:
	@go run main.go worker

# Build the swagger docs
build-docs:
	@swag fmt
	@swag init --parseDependency

# Run tests
test:
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Format Go code
fmt:
	@go fmt ./...

# Run static analysis
vet:
	@go vet ./...

# Verify dependencies
verify:
	@go mod verify

# Clean up dependencies
tidy:
	@go mod tidy

# Clean build artifacts
clean:
	@rm -f app coverage.out coverage.html

# Build the application
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app
	@docker build -t ccr.ccs.tencentyun.com/wwww/bill:latest .
	@docker push ccr.ccs.tencentyun.com/wwww/bill:latest

# Development checks
check: fmt vet verify test

# Build and run locally
local: build-docs build
	@./app

# Help command
help:
	@echo "Available commands:"
	@echo "  dev           - Build docs and run server"
	@echo "  run           - Run HTTP API server"
	@echo "  run-worker    - Run background worker"
	@echo "  build-docs    - Generate API documentation"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format Go code"
	@echo "  vet           - Run static analysis"
	@echo "  verify        - Verify dependencies"
	@echo "  tidy          - Clean up dependencies"
	@echo "  build         - Build application binary"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  check         - Run all checks (fmt, vet, verify, test)"
	@echo "  clean         - Clean build artifacts"
	@echo "  add           - Git add all changes"
	@echo "  commit        - Git commit with timestamp"
	@echo "  push          - Git push to master"
	@echo "  all           - Git add, commit, and push"