# Modular Monolith Migration Plan

## Overview
This document outlines the plan to migrate the current microservices-style workspace (authuser, catalog, core) into a single modular monolith with clear boundaries.

## Current Architecture Issues

### Problem 1: Hybrid Architecture Confusion
- Code is organized as separate Go modules (workspace with multiple go.mod files)
- But services communicate via **direct imports** (not HTTP/gRPC)
- This creates confusion: looks like microservices, behaves like monolith

### Problem 2: Unused HTTP Handlers Between Modules
- `catalogHandler` was registered in auth tests but never used
- Partner service uses `ports.CatalogService` directly (not HTTP)
- RabbitMQ catalog consumer syncs cache, but validations should use direct queries

### Problem 3: Import Path Complexity
- Current: `github.com/Leviosa-care/authuser/internal/domain`
- Current: `github.com/Leviosa-care/catalog/internal/domain`
- Current: `github.com/Leviosa-care/core/errs`
- These suggest separate deployable services when they're not

## Target Architecture

### Single Go Module
```
github.com/Leviosa-care/leviosa/backend
```

### Directory Structure
```
backend/
├── go.mod                     # Single module, single source of truth
├── go.sum
├── cmd/
│   └── app/
│       └── main.go            # Application entry point
├── internal/
│   ├── app/
│   │   ├── server.go          # HTTP server setup
│   │   └── wiring.go          # Dependency injection
│   ├── authuser/              # Logical module (not Go module)
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/    # Renamed from adapters
│   │   │   ├── postgres/
│   │   │   ├── redis/
│   │   │   ├── rabbitmq/
│   │   │   ├── stripe/
│   │   │   └── migrations/    # Per-module migrations
│   │   ├── interface/         # HTTP handlers
│   │   │   ├── auth/
│   │   │   ├── user/
│   │   │   └── partner/
│   │   └── ports/
│   ├── catalog/
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   │   ├── postgres/
│   │   │   ├── rabbitmq/
│   │   │   └── migrations/
│   │   ├── interface/
│   │   └── ports/
│   └── common/                # Shared kernel (current core/)
│       ├── auth/
│       ├── contracts/
│       ├── ctxutil/
│       ├── envmode/
│       ├── errs/
│       ├── httpx/
│       ├── logger/
│       ├── messaging/
│       ├── middleware/
│       ├── testutils/
│       └── validation/
└── test/                      # Integration tests
    ├── integration/
    │   ├── authuser/
    │   └── catalog/
    └── helpers/
```

### New Import Paths
```go
// Old
import "github.com/Leviosa-care/authuser/internal/domain"
import "github.com/Leviosa-care/catalog/internal/domain"
import "github.com/Leviosa-care/core/errs"

// New
import "github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
import "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
import "github.com/Leviosa-care/leviosa/backend/internal/common/errs"
```

## Migration Phases

### Phase 1: Create Root go.mod (1 session)
**Goal**: Establish single Go module

**Steps**:
1. Create `backend/go.mod` with module path: `github.com/Leviosa-care/leviosa/backend`
2. Consolidate dependencies from authuser, catalog, and core go.mod files
3. Delete `backend/go.work`
4. Delete individual go.mod files in authuser/, catalog/, core/
5. Run `go mod tidy`

**Verification**:
- `go mod tidy` succeeds
- No duplicate dependencies

---

### Phase 2: Move core/ to internal/common/ (1 session)
**Goal**: Make shared code part of the monolith

**Steps**:
1. `mv core/ internal/common/`
2. Automated import path updates:
   ```bash
   find . -type f -name "*.go" -exec sed -i \
     's|github.com/Leviosa-care/core/|github.com/Leviosa-care/leviosa/backend/internal/common/|g' {} +
   ```
3. Run `go mod tidy`
4. Run `go test ./...` to verify

**Impact**:
- ~981 import statements updated
- No functional changes

**Verification**:
- All tests pass
- No compilation errors

---

### Phase 3: Restructure authuser/ Module (1 session)
**Goal**: Move authuser into internal/ with proper naming

**Steps**:
1. Rename directories:
   - `authuser/internal/adapters/` → `infrastructure/`
   - `authuser/internal/adapters/http/` → `interface/`
2. Move: `mv authuser/internal/* internal/authuser/`
3. Move: `mv authuser/test/ test/integration/authuser/`
4. Split migrations: Move auth/user/partner migrations to `internal/authuser/infrastructure/migrations/`
5. Update import paths:
   ```bash
   find . -type f -name "*.go" -exec sed -i \
     's|github.com/Leviosa-care/authuser/internal/|github.com/Leviosa-care/leviosa/backend/internal/authuser/|g' {} +
   ```
6. Update migration loading to use new embed path
7. Run `go mod tidy && go test ./...`

**Verification**:
- Tests pass
- Migrations load correctly

---

### Phase 4: Restructure catalog/ Module (1 session)
**Goal**: Move catalog into internal/

**Steps**:
1. Same directory renaming as authuser
2. Move: `mv catalog/internal/* internal/catalog/`
3. Move: `mv catalog/test/ test/integration/catalog/`
4. Split migrations: Move catalog migrations to `internal/catalog/infrastructure/migrations/`
5. Update import paths
6. Run `go mod tidy && go test ./...`

**Verification**:
- Tests pass
- Migrations load correctly

---

### Phase 5: Establish Cross-Module Dependencies (1 session)
**Goal**: Implement hybrid validation strategy

**Steps**:
1. Create catalog ports for authuser:
   ```go
   // internal/catalog/ports/repository.go
   type CategoryRepository interface {
       ExistsByIDs(ctx context.Context, ids []uuid.UUID) (bool, error)
   }
   ```

2. Implement in catalog infrastructure:
   ```go
   // internal/catalog/infrastructure/postgres/category/exists_by_ids.go
   func (r *Repository) ExistsByIDs(ctx context.Context, ids []uuid.UUID) (bool, error) {
       query := `SELECT COUNT(*) FROM categories WHERE id = ANY($1)`
       var count int
       err := r.pool.QueryRow(ctx, query, ids).Scan(&count)
       return count == len(ids), err
   }
   ```

3. Update authuser partner service:
   ```go
   import catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

   type Service struct {
       categoryRepo catalogPorts.CategoryRepository  // Direct validation
       catalogCache ports.CatalogCache              // Read performance
   }
   ```

4. Implement hybrid validation in create/update partner:
   ```go
   // Critical: Direct query
   valid, err := s.categoryRepo.ExistsByIDs(ctx, req.CategoryIDs)

   // Non-critical: Cache
   categories := s.catalogCache.GetCategories()
   ```

**Verification**:
- Partner validation uses direct queries
- Tests insert catalog data via SQL (no RabbitMQ events)

---

### Phase 6: Make RabbitMQ Optional (1 session)
**Goal**: Preserve RabbitMQ code but disable by default

**Steps**:
1. Add feature flag to partner service:
   ```go
   func New(..., mqConn *amqp.Connection, ...) (*Service, error) {
       if os.Getenv("ENABLE_CATALOG_MQ_SYNC") == "true" && mqConn != nil {
           // Initialize catalog consumer
           s.catalogConsumer = setupCatalogConsumer(...)
       } else {
           log.Info("Catalog MQ consumer disabled (monolith mode)")
       }
   }
   ```

2. Add package comment to rabbitmq files:
   ```go
   // Package rabbitmq provides optional RabbitMQ integration.
   // In monolith mode (default), consumers are disabled.
   // Enable: ENABLE_CATALOG_MQ_SYNC=true
   ```

3. Update tests to pass `nil` for mqConn

**Verification**:
- Catalog consumer not initialized by default
- Can re-enable with environment variable

---

### Phase 7: Create Application Wiring Layer (1 session)
**Goal**: Centralized dependency injection

**Steps**:
1. Create `internal/app/wiring.go`:
   - Container struct with all repositories and services
   - NewContainer function that wires dependencies
   - Cross-module dependencies passed via constructor

2. Create `internal/app/server.go`:
   - SetupRouter that registers all module routes
   - Start method with middleware

3. Create `cmd/app/main.go`:
   - Load config
   - Create container
   - Start server

**Example**:
```go
// internal/app/wiring.go
func NewContainer(ctx context.Context, cfg *Config) (*Container, error) {
    // Infrastructure
    pool := initPostgres(ctx, cfg.DatabaseURL)
    crypto := initCrypto(ctx, cfg.VaultAddr)

    // Repositories (both modules)
    categoryRepo := catalogInfra.NewCategoryRepository(ctx, pool)
    partnerRepo := authuserInfra.NewPartnerRepository(ctx, pool)

    // Services with cross-module deps
    partnerService := partner.New(
        partnerRepo,
        categoryRepo,  // Direct catalog dependency
        nil,          // No RabbitMQ in monolith
        crypto,
    )

    return &Container{...}, nil
}
```

**Verification**:
- Application starts
- All routes registered
- Cross-module calls work

---

### Phase 8: Update Integration Tests (1 session)
**Goal**: Update tests for new structure

**Steps**:
1. Update import paths in all test files
2. Update test setup to use direct repositories:
   ```go
   categoryRepo := catalogInfra.NewCategoryRepository(ctx, testPool)
   partnerService := partner.New(
       partnerRepo,
       categoryRepo,  // Direct dependency
       nil,          // No RabbitMQ
       crypto,
   )
   ```

3. Update test data helpers:
   ```go
   func InsertTestCategory(...) uuid.UUID {
       // Direct SQL insertion, no events
   }
   ```

4. Remove RabbitMQ catalog queue setup from tests

**Verification**:
- All integration tests pass
- Tests use direct SQL for catalog data

---

### Phase 9: Update Build and Deploy (1 session)
**Goal**: Single binary deployment

**Steps**:
1. Update Makefile:
   ```makefile
   build:
       go build -o bin/app ./cmd/app/main.go
   ```

2. Update Dockerfile:
   ```dockerfile
   RUN go build -o /app/bin/app ./cmd/app/main.go
   ```

3. Update docker-compose.yaml to use single service

**Verification**:
- Binary builds successfully
- Docker image builds
- Application starts in container

---

### Phase 10: Documentation (1 session)
**Goal**: Update all documentation

**Files to update**:
- `backend/CLAUDE.md` - Architecture overview
- `backend/README.md` - Build and run instructions
- Root `CLAUDE.md` - Backend reference
- Add `ARCHITECTURE.md` - Detailed design decisions

**Verification**:
- Documentation accurate
- Commands work as documented

---

## Migration Strategy

### Approach: Incremental with Testing
1. **One phase per session** - Don't rush
2. **Test after each phase** - Ensure everything works
3. **Commit after each phase** - Easy rollback if needed
4. **Document issues** - Track problems for future reference

### Risk Mitigation
- **Automated tools** - Use `sed`, `gofmt` for bulk changes
- **Compilation checks** - `go build ./...` after each phase
- **Test verification** - `go test ./...` after each phase
- **Git commits** - Commit working state after each phase

### Rollback Plan
- Each phase is in git
- Can revert to previous phase if issues occur
- Incremental approach limits blast radius

---

## Benefits After Migration

### Developer Experience
✅ **Simpler imports** - `backend/internal/authuser/domain`
✅ **Single dependency version** - No version conflicts
✅ **Free refactoring** - Move code between modules easily
✅ **Clear boundaries** - `internal/` enforces encapsulation
✅ **Easier testing** - Shared test utilities in one place

### Operational
✅ **Single binary deployment** - Simpler ops
✅ **Horizontal scaling** - Still scales like microservices
✅ **Shared database transactions** - Better consistency
✅ **Faster builds** - Single module, single build
✅ **Lower memory footprint** - No duplicate dependencies

### Architecture
✅ **True modular monolith** - Correct pattern implementation
✅ **Hybrid validation** - Direct queries + cache
✅ **Preserved RabbitMQ** - Easy microservice migration path
✅ **Per-module migrations** - Clear ownership
✅ **Ports pattern** - Testable, flexible dependencies

---

## Migration Timeline Estimate

### Conservative Approach (10 sessions)
- Phase 1: 30 minutes
- Phase 2: 1 hour (automated import updates)
- Phase 3: 1.5 hours
- Phase 4: 1 hour
- Phase 5: 1.5 hours
- Phase 6: 30 minutes
- Phase 7: 2 hours
- Phase 8: 1.5 hours
- Phase 9: 30 minutes
- Phase 10: 1 hour

**Total**: ~11 hours over 10 sessions

### Aggressive Approach (3 sessions)
- Session 1: Phases 1-3 (structure)
- Session 2: Phases 4-7 (dependencies + wiring)
- Session 3: Phases 8-10 (tests + deploy + docs)

**Total**: ~6 hours over 3 sessions (higher risk)

---

## Next Steps

When ready to proceed:
1. Review this plan
2. Choose migration approach (incremental vs aggressive)
3. Start with Phase 1 in a fresh session
4. Commit after each successful phase

## Current Status

✅ **Completed**: Removed unused catalogHandler from integration tests
⏸️ **Paused**: Full migration to be done in future sessions
📝 **Documented**: This plan for future reference
