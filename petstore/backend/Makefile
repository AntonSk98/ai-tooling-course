.PHONY: run test clean generate swagger models help install-tools dev build

# Default target - show help
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo ""
	@echo "  Development:"
	@echo "    make run              - Run the application"
	@echo "    make dev              - Run with live reload (requires air)"
	@echo "    make build            - Build binary to bin/petstore"
	@echo ""
	@echo "  Code Generation:"
	@echo "    make models           - Regenerate models from OpenAPI spec"
	@echo "    make swagger          - Regenerate Swagger documentation"
	@echo "    make generate         - Regenerate both models and Swagger docs"
	@echo ""
	@echo "  Testing:"
	@echo "    make test             - Run integration tests"
	@echo "    make test-coverage    - Run tests with coverage report"
	@echo ""
	@echo "  Dependencies:"
	@echo "    make deps             - Download and tidy Go dependencies"
	@echo "    make install-tools    - Install required development tools"
	@echo ""
	@echo "  Maintenance:"
	@echo "    make clean            - Clean build artifacts and database"
	@echo "    make verify           - Verify everything builds and tests pass"
	@echo ""

## run: Run the application
run:
	go run cmd/api/main.go

## test: Run integration tests
test:
	go test ./tests/integration/... -v

## test-coverage: Run tests with coverage report
test-coverage:
	go test ./tests/integration/... -v -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## models: Regenerate API models from OpenAPI specification
models:
	@echo "Regenerating models from petstore_swagger.yml..."
	oapi-codegen -config oapi-codegen.yaml petstore_swagger.yml
	@echo "✓ Models regenerated in internal/models/generated.go"

## swagger: Regenerate Swagger documentation
swagger:
	@echo "Regenerating Swagger documentation..."
	swag init -g cmd/api/main.go --output docs --parseDependency --parseInternal
	@echo "✓ Swagger docs regenerated in docs/"
	@echo "  - docs/docs.go"
	@echo "  - docs/swagger.json"
	@echo "  - docs/swagger.yaml"

## generate: Regenerate both models and Swagger documentation
generate: models swagger
	@echo "✓ All code generation complete!"

## clean: Clean build artifacts and database
clean:
	@echo "Cleaning build artifacts..."
	rm -f petstore.db
	rm -f coverage.out coverage.html
	rm -rf bin/
	@echo "✓ Clean complete"

## deps: Install and tidy Go dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies updated"

## install-tools: Install required development tools
install-tools:
	@echo "Installing development tools..."
	@echo "Installing oapi-codegen..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Installing swag..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing air (optional - for live reload)..."
	go install github.com/cosmtrek/air@latest
	@echo "✓ All tools installed"
	@echo ""
	@echo "Installed tools:"
	@echo "  - oapi-codegen: Generate models from OpenAPI spec"
	@echo "  - swag: Generate Swagger documentation"
	@echo "  - air: Live reload for development"

## build: Build binary
build:
	@echo "Building binary..."
	mkdir -p bin
	go build -o bin/petstore cmd/api/main.go
	@echo "✓ Binary built: bin/petstore"

## verify: Verify everything builds and tests pass
verify:
	@echo "Verifying project..."
	@echo "1. Building all packages..."
	go build ./...
	@echo "2. Running tests..."
	go test ./tests/integration/... -v
	@echo "✓ Verification complete - all checks passed!"

## dev: Run with live reload (requires air)
dev:
	@echo "Starting development server with live reload..."
	@echo "Make sure 'air' is installed: make install-tools"
	air
