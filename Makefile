.PHONY: all build test fmt lint clean install docs

# Variables
MODULE = github.com/juststeveking/sevalla-go
GO = go
GOFMT = gofmt
GOLINT = golangci-lint
GOTEST = $(GO) test

# Default target
all: fmt lint test build

# Verify the package builds
build:
	@echo "Verifying package builds..."
	@$(GO) build ./...

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -s -w .
	@$(GO) mod tidy

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not found. Running go vet instead..."; \
		echo "💡 To install golangci-lint, run: make install"; \
		$(GO) vet ./...; \
	fi

# Install dependencies
install:
	@echo "Installing dependencies..."
	@$(GO) mod download
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate documentation
docs:
	@echo "Generating documentation..."
	@$(GO) doc -all > API.md

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ coverage.out coverage.html

# Check for updates
check-updates:
	@echo "Checking for dependency updates..."
	@$(GO) list -u -m all

# Create a new release
release:
	@echo "Creating release..."
	@read -p "Enter version (e.g., v0.1.0): " version; \
	git tag -a $$version -m "Release $$version"; \
	git push origin $$version

# Run security checks
security:
	@echo "Running security checks..."
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...

# Generate mocks for testing
mocks:
	@echo "Generating mocks..."
	@$(GO) install github.com/golang/mock/mockgen@latest
	@mockgen -source=sevalla.go -destination=mocks/client_mock.go -package=mocks

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./...

# Quick check before commit
pre-commit: fmt lint test
	@echo "✅ Pre-commit checks passed!"

# Help target
help:
	@echo "Available targets:"
	@echo "  make all          - Format, lint, test, and verify build"
	@echo "  make build        - Verify the package compiles"
	@echo "  make test         - Run tests with coverage"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make install      - Install dependencies"
	@echo "  make docs         - Generate documentation"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make check-updates- Check for dependency updates"
	@echo "  make release      - Create a new release tag"
	@echo "  make security     - Run security checks"
	@echo "  make mocks        - Generate mocks for testing"
	@echo "  make bench        - Run benchmarks"
	@echo "  make pre-commit   - Run pre-commit checks"
	@echo "  make help         - Show this help message"
