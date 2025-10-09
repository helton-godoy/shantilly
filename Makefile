# Shantilly Makefile - Quality-First Build System

.PHONY: help fmt lint test test-race build coverage clean install dev

# Variables
BINARY_NAME=shantilly
BUILD_DIR=bin
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
COVERAGE_THRESHOLD=85

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format code with gofmt and goimports
	@echo "Running gofmt..."
	@gofmt -s -w .
	@echo "Running goimports..."
	@go run golang.org/x/tools/cmd/goimports@latest -w .
	@echo "✓ Code formatted"

lint: ## Run golangci-lint with errcheck as fatal
	@echo "Running golangci-lint..."
	@golangci-lint run --config .golangci.yml ./...
	@echo "✓ Linting passed"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...
	@echo "✓ Tests passed"

test-race: ## Run tests with race detector (Guard Rail)
	@echo "Running tests with race detector..."
	@go test -race -v ./...
	@echo "✓ Tests with race detector passed"

coverage: ## Generate coverage report and enforce 85% threshold
	@echo "Generating coverage report..."
	@go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@./scripts/coverage-report.sh $(COVERAGE_THRESHOLD)

build: fmt lint test-race ## Build binary with quality checks
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install ./cmd/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) installed"

clean: ## Clean build artifacts and coverage reports
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "✓ Cleaned"

dev: ## Run in development mode (build and run)
	@make build
	@./$(BUILD_DIR)/$(BINARY_NAME)

report: ## Generate automated status report with coverage, metrics, and code analysis
	@echo "Generating automated status report..."
	@./scripts/generate-status-report.sh
	@echo "✓ Status report generated: docs/reports/status/current-status.md"

# Dynamic Reports System Targets
capture-reports: ## Capture dynamic reports interactively
	@echo "Starting dynamic report capture..."
	@./scripts/capture-dynamic-reports.sh capture

list-dynamic-reports: ## List all captured dynamic reports
	@echo "Listing dynamic reports..."
	@./scripts/capture-dynamic-reports.sh list

backup-reports: ## Create backup of all dynamic reports
	@echo "Creating backup of dynamic reports..."
	@./scripts/backup-dynamic-reports.sh create

restore-reports: ## Restore reports from latest backup
	@echo "Restoring reports from latest backup..."
	@./scripts/backup-dynamic-reports.sh restore-latest

reports-stats: ## Show statistics of reports system
	@echo "Reports system statistics..."
	@echo ""
	@echo "Automated Reports:"
	@./scripts/backup-dynamic-reports.sh stats
	@echo ""
	@echo "Dynamic Reports:"
	@./scripts/capture-dynamic-reports.sh index

# CI Pipeline Target
ci: fmt lint coverage build ## Run full CI pipeline (format, lint, coverage, build)
	@echo "✓ CI Pipeline completed successfully"
