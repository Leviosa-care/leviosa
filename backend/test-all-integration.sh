#!/bin/bash
# test-all-integration.sh
# Script to run all integration tests after ENCX v0.6.0 migration

set -e

echo "=== Running All Integration Tests (ENCX v0.6.0) ==="
echo "Started at: $(date)"
echo

# Function to run integration tests for a service
run_service_tests() {
    local service=$1
    echo "Testing $service..."

    if [ -d "$service/test/integration" ]; then
        cd "$service"
        echo "  Running integration tests for $service..."
        if go test ./test/integration/... -v; then
            echo "  ✓ $service integration tests passed"
        else
            echo "  ✗ $service integration tests failed"
            return 1
        fi
        cd ..
    else
        echo "  No integration tests found for $service"
    fi
    echo
}

# Function to run repository tests
run_repo_tests() {
    local service=$1
    local repo=$2
    echo "Testing $service $repo repository..."

    if [ -f "$service/internal/adapters/$repo/main_test.go" ]; then
        cd "$service/internal/adapters/$repo"
        echo "  Running repository tests for $service/$repo..."
        if go test -v .; then
            echo "  ✓ $service/$repo repository tests passed"
        else
            echo "  ✗ $service/$repo repository tests failed"
            return 1
        fi
        cd ../../..
    else
        echo "  No repository tests found for $service/$repo"
    fi
    echo
}

# Test results
PASSED=0
FAILED=0

# Run AuthUser integration tests
echo "=== AuthUser Service Tests ==="
if run_service_tests "authuser"; then
    ((PASSED++))
else
    ((FAILED++))
fi

# Run AuthUser repository tests
echo "=== AuthUser Repository Tests ==="
if run_repo_tests "authuser" "postgres/user"; then
    ((PASSED++))
else
    ((FAILED++))
fi

if run_repo_tests "authuser" "redis/session"; then
    ((PASSED++))
else
    ((FAILED++))
fi

# Run Settings integration tests
echo "=== Settings Service Tests ==="
if run_service_tests "settings"; then
    ((PASSED++))
else
    ((FAILED++))
fi

# Run Booking integration tests
echo "=== Booking Service Tests ==="
if run_service_tests "booking"; then
    ((PASSED++))
else
    ((FAILED++))
fi

# Summary
echo "=== Test Summary ==="
echo "Started at: $(date)"
echo "Total test suites: $((PASSED + FAILED))"
echo "Passed: $PASSED"
echo "Failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo "🎉 All integration tests passed!"
    echo "✅ ENCX v0.6.0 migration verification successful"
    exit 0
else
    echo "❌ Some integration tests failed!"
    echo "❌ ENCX v0.6.0 migration needs attention"
    exit 1
fi