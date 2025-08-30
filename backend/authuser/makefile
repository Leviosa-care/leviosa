# Test execution commands
.PHONY: test test-unit test-integration test-all test-verbose test-coverage test-smoke test-race benchmark
.PHONY: test-postgres test-redis test-unit-user test-unit-otp test-unit-session
.PHONY: test-integration-auth test-integration-user test-integration-otp test-integration-session

# Unit tests by adapter type
test-unit:
	go test -v -count=1 ./internal/adapters/...

test-postgres:
	go test -v -count=1 ./internal/adapters/postgres/...

test-redis:
	go test -v -count=1 ./internal/adapters/redis/...

# Unit tests by specific adapter
test-unit-user:
	go test -v -count=1 ./internal/adapters/postgres/user/...

test-unit-otp:
	go test -v -count=1 ./internal/adapters/redis/otp/...

test-unit-session:
	go test -v -count=1 ./internal/adapters/redis/session/...

# Integration tests by module
test-integration:
	go test -v -count=1 ./test/integration/...

test-integration-auth:
	go test -v -count=1 ./test/integration/auth/...

test-integration-user:
	go test -v -count=1 ./test/integration/user/...

test-integration-otp:
	go test -v -count=1 ./test/integration/otp/...

test-integration-session:
	go test -v -count=1 ./test/integration/session/...

# All tests
test-all:
	go test -v -count=1 ./...

# Run specific test by name within adapters
test-unit-user-%:
	go test -v -count=1 ./internal/adapters/postgres/user/... -run $(TEST)

test-unit-otp-%:
	go test -v -count=1 ./internal/adapters/redis/otp/... -run $(TEST)

test-unit-session-%:
	go test -v -count=1 ./internal/adapters/redis/session/... -run $(TEST)

# Run specific integration test by module
test-integration-auth-%:
	go test -v -count=1 ./test/integration/auth/... -run $(TEST)

test-integration-user-%:
	go test -v -count=1 ./test/integration/user/... -run $(TEST)

test-integration-otp-%:
	go test -v -count=1 ./test/integration/otp/... -run $(TEST)

test-integration-session-%:
	go test -v -count=1 ./test/integration/session/... -run $(TEST)

# Advanced testing
test-verbose:
	go test -v -count=1 -json ./... | tee test-output.json

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-smoke:
	go test -v -count=1 -short ./...

test-race:
	go test -v -race -count=1 ./...

benchmark:
	go test -v -bench=. -benchmem ./...

# Default test command (backwards compatibility)
test: test-unit 
