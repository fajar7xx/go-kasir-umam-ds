GOBIN ?= $$(go env GOPATH)/bin
APP_PORT ?= 8080
APP_NAME := go-kasir-umam-ds

# ==================== APP COMMANDS ====================
.PHONY: run
run:
	@echo "🚀 Starting $(APP_NAME) on port $(APP_PORT)..."
	go run main.go

.PHONY: build
build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) main.go
	@echo "✅ Build complete! Run with: ./bin/$(APP_NAME)"

.PHONY: dev
dev:
	@echo "🔄 Running in development mode with auto-reload..."
	@command -v air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

.PHONY: stop
stop:
	@echo "⏹️  Stopping $(APP_NAME)..."
	@pkill -f "go-kasir-umam-ds" || echo "No process found"

# ==================== TESTING COMMANDS ====================
.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: test
test:
	@echo "🧪 Running all tests with race detector..."
	go test -v -race ./...

.PHONY: test-fast
test-fast:
	@echo "⚡ Running tests (fast mode, no race detector)..."
	go test -v ./...

.PHONY: test-verbose
test-verbose:
	@echo "📝 Running tests with verbose output..."
	go test -v -race -count=1 ./...

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	@echo "📊 Running tests with coverage analysis..."
	go test ./... \
		-coverprofile=./cover.out \
		-covermode=atomic \
		-coverpkg=./handlers/...,./internal/services/...,./internal/repositories/...,./utils/...
	@echo "Checking coverage thresholds..."
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

.PHONY: coverage-html
coverage-html:
	@echo "📈 Generating HTML coverage report..."
	go test ./... -coverprofile=cover.out -covermode=atomic
	go tool cover -html=cover.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

.PHONY: coverage-report
coverage-report:
	@echo "=== Coverage Summary ==="
	go test ./... -cover | grep -E "coverage:|ok"

# ==================== CODE QUALITY ====================
.PHONY: lint
lint:
	@echo "🔍 Running linter..."
	@command -v golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

.PHONY: fmt
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

.PHONY: vet
vet:
	@echo "🔎 Running go vet..."
	go vet ./...
	@echo "✅ No issues found"

# ==================== DEVELOPMENT TOOLS ====================
.PHONY: deps
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated"

.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f bin/$(APP_NAME)
	rm -f cover.out coverage.html
	go clean
	@echo "✅ Clean complete"

# ==================== DATABASE ====================
.PHONY: db-migrate
db-migrate:
	@echo "🗄️  Running migrations..."
	@echo "Note: Update this target with your migration command"

# ==================== HELP ====================
.PHONY: help
help:
	@echo ""
	@echo "╔════════════════════════════════════════════════════════════════╗"
	@echo "║        $(APP_NAME) - Makefile Commands              ║"
	@echo "╚════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "🚀 APPLICATION"
	@echo "  make run              - Run the application ($(APP_PORT))"
	@echo "  make build            - Build binary to ./bin/"
	@echo "  make dev              - Run with auto-reload (requires air)"
	@echo "  make stop             - Stop running application"
	@echo ""
	@echo "🧪 TESTING"
	@echo "  make test             - Run all tests with race detector"
	@echo "  make test-fast        - Run tests without race detector"
	@echo "  make test-verbose     - Run tests with verbose output"
	@echo "  make check-coverage   - Check coverage meets thresholds"
	@echo "  make coverage-html    - Generate HTML coverage report"
	@echo "  make coverage-report  - Show coverage summary"
	@echo ""
	@echo "✨ CODE QUALITY"
	@echo "  make lint             - Run linter (golangci-lint)"
	@echo "  make fmt              - Format all code"
	@echo "  make vet              - Run go vet"
	@echo ""
	@echo "🛠️  TOOLS"
	@echo "  make deps             - Download and tidy dependencies"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make db-migrate       - Run database migrations"
	@echo ""
	@echo "💡 EXAMPLES:"
	@echo "  make run              # Quick start"
	@echo "  make dev              # Development with auto-reload"
	@echo "  make test coverage    # Test + coverage in one go"
	@echo "  make clean run        # Fresh build and run"
	@echo ""

# ==================== DEFAULT ====================
.DEFAULT_GOAL := help
