GOBIN ?= $$(go env GOPATH)/bin

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: test
test:
	go test -v -race ./...

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	@echo "Running tests with coverage for business logic packages..."
	go test ./... \
		-coverprofile=./cover.out \
		-covermode=atomic \
		-coverpkg=./handlers/...,./internal/services/...,./internal/repositories/...,./utils/...
	@echo "Checking coverage thresholds..."
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

.PHONY: coverage-html
coverage-html:
	go test ./... -coverprofile=cover.out -covermode=atomic
	go tool cover -html=cover.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: coverage-report
coverage-report:
	@echo "=== Coverage Summary ==="
	go test ./... -cover | grep -E "coverage:|ok"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests with race detector"
	@echo "  make check-coverage    - Check coverage meets thresholds (for CI/CD)"
	@echo "  make coverage-html     - Generate HTML coverage report"
	@echo "  make coverage-report   - Show coverage summary"
	@echo "  make help              - Show this help message"
