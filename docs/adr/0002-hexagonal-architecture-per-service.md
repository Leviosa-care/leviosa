# ADR-0002: Hexagonal Architecture Per Service Module

**Status**: Accepted  
**Date**: 2025-07

---

## Context

Each service module needs a structure that isolates business logic from infrastructure, allows the integration test suite to swap real adapters for test doubles, and keeps the domain model free of framework dependencies.

## Decision

Each service module follows the hexagonal (ports & adapters) pattern with four layers:

```
internal/[service]/
├── domain/         # Pure business logic — no external imports
├── ports/          # Go interfaces for all dependencies (repository, external services)
├── application/    # Use cases that orchestrate domain + ports
└── infrastructure/ # Concrete adapter implementations (postgres, redis, stripe, s3)
```

The DI container (`internal/app/container.go`) wires adapters to ports at startup. Application code only depends on the `ports` interfaces; it never imports infrastructure packages directly.

## Consequences

**Positive:**
- Domain layer is pure Go with zero external dependencies — trivially testable
- Integration tests can use real adapters (PostgreSQL, Redis, S3, Vault) without touching application logic
- Swapping an adapter (e.g., replacing Redis with Memcached) requires no changes to application or domain code
- Error classification (`errs.ClassifyPgError`) stays at the adapter boundary, keeping domain errors clean

**Negative:**
- More files and packages per feature compared to a flat structure
- New developers must understand the layering before contributing effectively

**Enforcement**: Import cycle detection in Go naturally enforces the dependency rule — `domain` cannot import `application`, `application` cannot import `infrastructure`.
