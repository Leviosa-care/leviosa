# Test execution commands
.PHONY: test test-unit test-integration test-all test-verbose test-coverage test-smoke test-race benchmark

# Unit tests (repository layer)
test-unit:
	go test -v -count=1 ./internal/adapters/...

# Integration tests (full stack)  
test-integration:
	go test -v -count=1 ./test/integration/...

# All tests
test-all:
	go test -v -count=1 ./...

# Run specific unit test
test-unit-%:
	go test -v -count=1 ./internal/adapters/$(subst test-unit-,,$@)/... -run $(TEST)

# Run specific integration test
test-integration-%:
	go test -v -count=1 ./test/integration/... -run $(TEST)

# Verbose test output with timing
test-verbose:
	go test -v -count=1 -json ./... | tee test-output.json

# Test coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Quick smoke test (fastest critical tests only)
test-smoke:
	go test -v -count=1 -short ./...

# Test with race detection
test-race:
	go test -v -race -count=1 ./...

# Benchmark tests
benchmark:
	go test -v -bench=. -benchmem ./...

# Default test command (backwards compatibility)
test: test-unit 
