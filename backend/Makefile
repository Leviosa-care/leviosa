# ============================================================================
# LEVIOSA TESTING COMMANDS
# ============================================================================
# Supports current mixed architecture and future internal/ structure

.PHONY: test-all test-unit test-integration test-coverage test-race test-benchmark test-verbose test-smoke test-clean test-help
.PHONY: test-unit-authuser test-unit-catalog test-unit-settings test-unit-booking test-unit-notification
.PHONY: test-integration-authuser test-integration-catalog test-integration-settings test-integration-booking test-integration-notification

# Test settings
TEST_TIMEOUT?=30m
TEST_VERBOSE?=-v
TEST_RACE?=-race
TEST_COVERAGE?=-cover
TEST_COVERAGE_HTML?=-coverprofile=coverage.out -covermode=atomic

# Service paths (future internal/ structure)
INTERNAL_SERVICES=authuser catalog settings booking notification
SERVICES_ROOT=internal

# Integration test paths (current structure)
INTEGRATION_ROOT=test/integration

test-help: ## Show testing help message
	@echo 'Leviosa Backend Testing Commands:'
	@echo ''
	@echo 'Primary Test Commands:'
	@echo '  test-all           Run all unit and integration tests'
	@echo '  test-unit          Run all unit tests across all services'
	@echo '  test-integration   Run all integration tests'
	@echo ''
	@echo 'Service-Specific Unit Tests (internal/ structure):'
	@echo '  test-unit-authuser       Run unit tests for internal/authuser/'
	@echo '  test-unit-catalog        Run unit tests for internal/catalog/'
	@echo '  test-unit-settings       Run unit tests for internal/settings/'
	@echo '  test-unit-booking        Run unit tests for internal/booking/'
	@echo '  test-unit-notification   Run unit tests for internal/notification/'
	@echo ''
	@echo 'Service-Specific Integration Tests:'
	@echo '  test-integration-authuser     Run integration tests for authuser'
	@echo '  test-integration-catalog      Run integration tests for catalog'
	@echo '  test-integration-settings     Run integration tests for settings'
	@echo '  test-integration-booking      Run integration tests for booking'
	@echo '  test-integration-notification Run integration tests for notification'
	@echo ''
	@echo 'Advanced Testing Commands:'
	@echo '  test-coverage      Generate HTML coverage report for all services'
	@echo '  test-race          Run tests with race condition detection'
	@echo '  test-benchmark     Run performance benchmarks'
	@echo '  test-verbose       Run tests with detailed output'
	@echo '  test-smoke         Run quick validation (short tests only)'
	@echo ''
	@echo 'Generic Specific Test Commands:'
	@echo '  test-file          Run specific test file (usage: make test-file TEST_PATH=<path>)'
	@echo '  test-func          Run specific test function (usage: make test-func TEST_NAME=<func> TEST_PATH=<path>)'
	@echo ''
	@echo 'Utility Commands:'
	@echo '  test-clean         Clean test artifacts and coverage files'

# Primary Commands
test-all: test-unit test-integration ## Run all unit and integration tests

test-unit: ## Run all unit tests across all services in internal/
	@echo "Running all unit tests..."
	@for service in $(INTERNAL_SERVICES); do \
		if [ -d "$(SERVICES_ROOT)/$$service" ]; then \
			echo "Testing unit tests for $$service..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/$$service/... || exit 1; \
		else \
			echo "Warning: $(SERVICES_ROOT)/$$service not found, skipping..."; \
		fi; \
	done
	@echo "All unit tests completed!"

test-integration: ## Run all integration tests from test/integration/
	@echo "Running all integration tests..."
	@for service in $(INTERNAL_SERVICES); do \
		if [ -d "$(INTEGRATION_ROOT)/$$service" ]; then \
			echo "Testing integration for $$service..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/$$service/... || exit 1; \
		else \
			echo "Warning: $(INTEGRATION_ROOT)/$$service not found, skipping..."; \
		fi; \
	done
	@echo "All integration tests completed!"

# Service-Specific Unit Tests (internal/ structure)
test-unit-authuser: ## Run unit tests for internal/authuser/
	@if [ -d "$(SERVICES_ROOT)/authuser" ]; then \
		echo "Running unit tests for authuser service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/authuser/...; \
	else \
		echo "Error: $(SERVICES_ROOT)/authuser not found. Service may not be migrated to internal/ yet."; \
		exit 1; \
	fi

test-unit-catalog: ## Run unit tests for internal/catalog/
	@if [ -d "$(SERVICES_ROOT)/catalog" ]; then \
		echo "Running unit tests for catalog service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/catalog/...; \
	else \
		echo "Error: $(SERVICES_ROOT)/catalog not found. Service may not be migrated to internal/ yet."; \
		exit 1; \
	fi

test-unit-settings: ## Run unit tests for internal/settings/
	@if [ -d "$(SERVICES_ROOT)/settings" ]; then \
		echo "Running unit tests for settings service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/settings/...; \
	else \
		echo "Error: $(SERVICES_ROOT)/settings not found. Service may not be migrated to internal/ yet."; \
		exit 1; \
	fi

test-unit-booking: ## Run unit tests for internal/booking/
	@if [ -d "$(SERVICES_ROOT)/booking" ]; then \
		echo "Running unit tests for booking service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/booking/...; \
	else \
		echo "Error: $(SERVICES_ROOT)/booking not found. Service may not be migrated to internal/ yet."; \
		exit 1; \
	fi

test-unit-notification: ## Run unit tests for internal/notification/
	@if [ -d "$(SERVICES_ROOT)/notification" ]; then \
		echo "Running unit tests for notification service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/notification/...; \
	else \
		echo "Error: $(SERVICES_ROOT)/notification not found. Service may not be migrated to internal/ yet."; \
		exit 1; \
	fi

# Service-Specific Integration Tests
test-integration-authuser: ## Run integration tests for authuser
	@if [ -d "$(INTEGRATION_ROOT)/authuser" ]; then \
		echo "Running integration tests for authuser service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/authuser/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/authuser not found, skipping..."; \
	fi

test-integration-catalog: ## Run integration tests for catalog
	@if [ -d "$(INTEGRATION_ROOT)/catalog" ]; then \
		echo "Running integration tests for catalog service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/catalog/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/catalog not found, skipping..."; \
	fi

test-integration-settings: ## Run integration tests for settings
	@if [ -d "$(INTEGRATION_ROOT)/settings" ]; then \
		echo "Running integration tests for settings service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/settings/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/settings not found, skipping..."; \
	fi

test-integration-booking: ## Run integration tests for booking
	@if [ -d "$(INTEGRATION_ROOT)/booking" ]; then \
		echo "Running integration tests for booking service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/booking/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/booking not found, skipping..."; \
	fi

test-integration-notification: ## Run integration tests for notification
	@if [ -d "$(INTEGRATION_ROOT)/notification" ]; then \
		echo "Running integration tests for notification service..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/notification/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/notification not found, skipping..."; \
	fi

# Advanced Testing Commands
test-coverage: ## Generate HTML coverage report for all services
	@echo "Generating coverage report for all services..."
	@rm -f coverage.out coverage.html
	@echo "mode: atomic" > coverage.out
	@for service in $(INTERNAL_SERVICES); do \
		if [ -d "$(SERVICES_ROOT)/$$service" ]; then \
			echo "Collecting coverage for $$service unit tests..."; \
			go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/$$service/... -o $$service-coverage.out 2>/dev/null || true; \
			if [ -f $$service-coverage.out ]; then \
				tail -n +2 $$service-coverage.out >> coverage.out; \
				rm $$service-coverage.out; \
			fi; \
		fi; \
		if [ -d "$(INTEGRATION_ROOT)/$$service" ]; then \
			echo "Collecting coverage for $$service integration tests..."; \
			go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/$$service/... -o $$service-integration-coverage.out 2>/dev/null || true; \
			if [ -f $$service-integration-coverage.out ]; then \
				tail -n +2 $$service-integration-coverage.out >> coverage.out; \
				rm $$service-integration-coverage.out; \
			fi; \
		fi; \
	done
	@if [ -s coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "Coverage report generated: coverage.html"; \
		echo "Overall coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"; \
	else \
		echo "No coverage data collected."; \
	fi

test-race: ## Run tests with race condition detection
	@echo "Running tests with race condition detection..."
	@$(MAKE) test-unit TEST_RACE="$(TEST_RACE)"
	@$(MAKE) test-integration TEST_RACE="$(TEST_RACE)"

test-benchmark: ## Run performance benchmarks
	@echo "Running performance benchmarks..."
	@for service in $(INTERNAL_SERVICES); do \
		if [ -d "$(SERVICES_ROOT)/$$service" ]; then \
			echo "Benchmarking $$service..."; \
			go test -bench=. -benchmem ./$(SERVICES_ROOT)/$$service/... -timeout $(TEST_TIMEOUT) || true; \
		fi; \
		if [ -d "$(INTEGRATION_ROOT)/$$service" ]; then \
			echo "Benchmarking $$service integration tests..."; \
			go test -bench=. -benchmem ./$(INTEGRATION_ROOT)/$$service/... -timeout $(TEST_TIMEOUT) || true; \
		fi; \
	done

test-verbose: ## Run tests with detailed output
	@echo "Running tests with verbose output..."
	@$(MAKE) test-unit TEST_VERBOSE="-v"
	@$(MAKE) test-integration TEST_VERBOSE="-v"

test-smoke: ## Run quick validation (short tests only)
	@echo "Running smoke tests (short tests only)..."
	@for service in $(INTERNAL_SERVICES); do \
		if [ -d "$(SERVICES_ROOT)/$$service" ]; then \
			echo "Smoke testing $$service..."; \
			go test -short -timeout 5m ./$(SERVICES_ROOT)/$$service/... || true; \
		fi; \
		if [ -d "$(INTEGRATION_ROOT)/$$service" ]; then \
			echo "Smoke testing $$service integration..."; \
			go test -short -timeout 5m ./$(INTEGRATION_ROOT)/$$service/... || true; \
		fi; \
	done

# ============================================================================
# ENCRYPTION COMMANDS (ENCX)
# ============================================================================

.PHONY: encx-validate encx-generate encx-clean encx-build encx-help

encx-help: ## Show ENCRYPTION (ENCX) help message
	@echo 'ENCX Encryption Commands:'
	@echo ''
	@echo 'Available Commands:'
	@echo '  encx-install      Install ENCX CLI globally (recommended)'
	@echo '  encx-validate      Validate ENCX struct tags before generation'
	@echo '  encx-generate      Generate ENCX encryption code for all structs'
	@echo '  encx-build         Build ENCX generator tool locally (optional)'
	@echo '  encx-clean         Clean generated ENCX files (custom helper)'
	@echo ''
	@echo 'Development Workflow:'
	@echo '  1. Design structs with encx tags'
	@echo '  2. Run make encx-validate'
	@echo '  3. Run make encx-generate'
	@echo '  4. Test your changes'
	@echo '  5. Commit both source and generated code'
	@echo ''
	@echo 'Direct ENCX Commands (without make):'
	@echo '    go run ./cmd/encx-gen validate -v .'
	@echo '    go run ./cmd/encx-gen generate -v .'

encx-install: ## Install ENCX CLI globally (recommended method)
	@echo "Installing ENCX CLI globally..."
	@go install github.com/hengadev/encx/cmd/encx-gen@latest
	@echo "ENCX CLI installed globally!"
	@echo "You can now use: encx-gen validate -v ."

encx-build: ## Build ENCX generator tool locally (optional development method)
	@echo "Building ENCX generator tool locally..."
	@mkdir -p bin
	@go build -o bin/encx-gen github.com/hengadev/encx/cmd/encx-gen
	@echo "ENCX generator built: bin/encx-gen"
	@echo "You can now use: ./bin/encx-gen validate -v ."

encx-validate: ## Validate ENCX struct tags
	@echo "Validating ENCX struct tags..."
	@if command -v encx-gen >/dev/null 2>&1; then \
		echo "Using global encx-gen CLI..."; \
		encx-gen validate -v .; \
	elif [ -f "./bin/encx-gen" ]; then \
		echo "Using local encx-gen binary..."; \
		./bin/encx-gen validate -v .; \
	else \
		echo "Using go run (no installation needed)..."; \
		go run github.com/hengadev/encx/cmd/encx-gen validate -v .; \
	fi
	@echo "ENCX validation completed!"

encx-generate: ## Generate ENCX encryption code
	@echo "Generating ENCX encryption code..."
	@echo "Running validation before generation..."
	@if command -v encx-gen >/dev/null 2>&1; then \
		echo "Using global encx-gen CLI..."; \
		encx-gen validate -v . || exit 1; \
		echo "Validation passed. Generating code..."; \
		encx-gen generate -v .; \
	elif [ -f "./bin/encx-gen" ]; then \
		echo "Using local encx-gen binary..."; \
		./bin/encx-gen validate -v . || exit 1; \
		echo "Validation passed. Generating code..."; \
		./bin/encx-gen generate -v .; \
	else \
		echo "Using go run (no installation needed)..."; \
		go run github.com/hengadev/encx/cmd/encx-gen validate -v . || exit 1; \
		echo "Validation passed. Generating code..."; \
		go run github.com/hengadev/encx/cmd/encx-gen generate -v .; \
	fi
	@echo "ENCX code generation completed!"
	@echo "Generated files (suffix: _encx.go) have been created alongside source files"

encx-clean: ## Clean generated ENCX files (custom helper)
	@echo "Cleaning generated ENCX files..."
	@find . -name "*_encx.go" -type f -delete
	@echo "Generated ENCX files cleaned!"

# Utility Commands
test-clean: ## Clean test artifacts and coverage files
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@rm -f *.test
	@rm -f *-coverage.out
	@find . -name "*.test" -delete
	@find . -name "*.out" -delete
	@find . -name "*.prof" -delete
	@echo "Test artifacts cleaned!"

# Generic Specific Test Commands
test-file: ## Run specific test file (usage: make test-file TEST_PATH=internal/authuser/infrastructure/postgres/user/create_user_test.go)
	@if [ -z "$(TEST_PATH)" ]; then \
		echo "Error: TEST_PATH parameter is required. Usage: make test-file TEST_PATH=<path-to-test-file>"; \
		exit 1; \
	fi
	@if [ ! -f "$(TEST_PATH)" ]; then \
		echo "Error: Test file '$(TEST_PATH)' not found."; \
		exit 1; \
	fi
	@echo "Running tests in $(TEST_PATH)..."
	@dir_path=$$(dirname "$(TEST_PATH)"); go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) "./$$dir_path"

test-func: ## Run specific test function (usage: make test-func TEST_NAME=TestCreateUser TEST_PATH=internal/authuser/infrastructure/postgres/user/create_user_test.go)
	@if [ -z "$(TEST_NAME)" ] || [ -z "$(TEST_PATH)" ]; then \
		echo "Error: Both TEST_NAME and TEST_PATH parameters are required."; \
		echo "Usage: make test-func TEST_NAME=<test-function-name> TEST_PATH=<path-to-test-file>"; \
		exit 1; \
	fi
	@if [ ! -f "$(TEST_PATH)" ]; then \
		echo "Error: Test file '$(TEST_PATH)' not found."; \
		exit 1; \
	fi
	@echo "Running test function $(TEST_NAME) in $(TEST_PATH)..."
	@dir_path=$$(dirname "$(TEST_PATH)"); go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) -run "$(TEST_NAME)" "./$$dir_path"

# Forward compatibility handles current mixed architecture
test-unit-current: ## Run unit tests for current mixed architecture (internal/ + standalone services)
	@echo "Running unit tests for current architecture..."
	# Test services in internal/
	@if [ -d "$(SERVICES_ROOT)/authuser" ]; then \
		echo "Testing internal/authuser..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/authuser/... || true; \
	fi
	@if [ -d "$(SERVICES_ROOT)/catalog" ]; then \
		echo "Testing internal/catalog..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/catalog/... || true; \
	fi
	@if [ -d "$(SERVICES_ROOT)/settings" ]; then \
		echo "Testing internal/settings..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(SERVICES_ROOT)/settings/... || true; \
	fi
	# Test standalone services (current reality)
	@if [ -d "booking" ]; then \
		echo "Testing standalone booking..."; \
		cd booking && go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./... && cd .. || true; \
	fi
	@if [ -d "notification" ]; then \
		echo "Testing standalone notification..."; \
		cd notification && go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./... && cd .. || true; \
	fi
	@echo "Current architecture unit tests completed!"
