# ADR-0005: Integration Tests Use Real Adapters via Testcontainers

**Status**: Accepted  
**Date**: 2025-07

---

## Context

The application's correctness depends heavily on database queries, Redis operations, S3 interactions, and Vault key retrieval. Mocking these in tests provides coverage of application logic but cannot catch schema mismatches, query bugs, migration regressions, or adapter configuration errors.

## Decision

Integration tests in `test/integration/` spin up real infrastructure using **testcontainers-go**: PostgreSQL, Redis, RabbitMQ, Vault (dev mode), and S3 (Localstack). Tests exercise the full stack from HTTP handler to database and back.

Unit tests (`*_test.go` alongside source) test pure domain logic with no external dependencies.

Mock adapters (e.g., `StripeMock`) are used only for external paid APIs where real calls would have side effects or costs.

The `testutils/` package provides helpers:
- `SetupServiceVault()` — provisions per-service `encx` keys in a test Vault instance
- `SetupDB()` — runs migrations and seeds test data
- Auth helpers for injecting test sessions

## Consequences

**Positive:**
- Migration regressions are caught before they reach staging
- Query bugs surface in CI, not in production
- Tests document the real behavior of adapters, not an idealized mock

**Negative:**
- Integration tests are slower than unit tests (container startup time)
- Requires Docker in the CI environment
- Test isolation requires careful teardown to avoid state bleed between test cases

**Policy**: Do not mock the database in integration tests. A test that only mocks the repository interface tells us nothing about whether the SQL is correct.
