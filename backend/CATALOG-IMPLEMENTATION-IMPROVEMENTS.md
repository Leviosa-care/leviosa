# Catalog Cache Implementation - Improvement Plan

## Executive Summary

The catalog event integration and in-memory cache implementation follows the correct architectural pattern (similar to Settings integration). However, there are several improvements needed in the areas of integration completeness, code quality, testing, and operational robustness.

**Overall Assessment**: ✅ **Good foundation** with room for improvement in integration, testing, and error handling.

---

## 1. Core Implementation Review

### ✅ What's Working Well

#### 1.1 Event Contracts (`core/contracts/catalog/events.go`)
- **Status**: ✅ Excellent
- Clean event definitions with proper validation methods
- Good JSON serialization structure
- Comprehensive validation logic for UUIDs, names, and business rules
- Well-documented event types

#### 1.2 In-Memory Cache (`authuser/internal/application/catalog/cache.go`)
- **Status**: ✅ Very Good
- Proper use of `sync.RWMutex` for thread safety
- Correct implementation of published-only filtering
- Returns copies to prevent external mutation
- Cascade deletion for categories → products
- Implements both read and write interfaces

#### 1.3 Consumer Implementation (`authuser/internal/adapters/rabbitmq/catalog_consumer.go`)
- **Status**: ✅ Good
- Proper event routing by routing key
- Good error handling with Nack/Ack
- Comprehensive logging
- Event validation before processing

#### 1.4 Domain Models (`authuser/internal/domain/catalog.go`)
- **Status**: ✅ Good
- Simplified models without database-specific fields
- Helper methods for status checking
- Follows the refactor plan correctly

---

## 2. Critical Issues

### ❌ 2.1 CreatePartner Method Not Using Catalog Cache Validation

**Location**: `authuser/internal/application/partner/create_partner.go:58-82`

**Problem**: The `CreatePartner` method still validates specializations using the old repository pattern:
```go
// Current implementation
for _, specializationID := range request.SpecializationIDs {
    // Verify specialization exists and is active
    specEncx, err := s.specializationRepo.GetSpecializationByID(ctx, specializationID)
    // ... more repository operations
}
```

**Issue**: This bypasses the catalog cache and queries the database directly, which:
- Defeats the purpose of the in-memory cache
- Creates unnecessary database load
- Doesn't validate against catalog published items
- Inconsistent with the refactor plan

**Expected Behavior**:
```go
// Should use catalog cache validation
if err := s.ValidatePartnerSpecializations(ctx, request.SpecializationIDs); err != nil {
    return nil, err
}
```

**Impact**: HIGH - Core business logic not using the cache

---

### ❌ 2.2 Missing Service Startup Wiring

**Problem**: The catalog cache and consumer are not wired into the main service startup.

**Missing in `cmd/leviosa/main.go`** (or equivalent):
```go
// Missing initialization
catalogCache := catalog.NewCatalogCache()
catalogConsumer := rabbitmq.NewCatalogConsumer(mqConn, catalogCache)

// Missing goroutine to start consumer
go func() {
    if err := catalogConsumer.Start(ctx); err != nil {
        if err != context.Canceled {
            log.Fatalf("Catalog consumer error: %v", err)
        }
    }
}()

// Missing catalog cache injection into partner service
partnerService := partner.New(
    partnerRepo,
    userRepo,
    specializationRepo,
    catalogCache,  // Missing parameter
    crypto,
    stripe,
)
```

**Impact**: CRITICAL - Implementation not functional without wiring

---

### ⚠️ 2.3 Exchange Type Inconsistency

**Locations**:
- `core/contracts/rabbitmq/catalog.go` - No exchange type specified
- `authuser/internal/adapters/rabbitmq/setup_catalog_consumer.go:15` - Uses `"topic"`
- `authuser/internal/adapters/rabbitmq/catalog_consumer.go:42` - Uses `"topic"`

**Problem**: The routing keys are simple strings like `"catalog.product.created"`, not wildcard patterns.

**Analysis**:
- **Topic exchange** is used for pattern matching (e.g., `catalog.*.created`, `catalog.product.*`)
- **Direct exchange** is used for exact matching

**Current routing keys**:
```go
ProductCreatedRoutingKey  = "catalog.product.created"
CategoryCreatedRoutingKey = "catalog.category.created"
```

**Recommendation**:
- **Option A** (Recommended): Use **direct exchange** since routing keys are exact matches
- **Option B**: Keep topic exchange but document that it's for future wildcard support

**Impact**: LOW - Works correctly but semantically inconsistent

---

### ⚠️ 2.4 Unnecessary Catalog Service Wrapper

**Location**: `authuser/internal/application/catalog/service.go`

**Problem**: The `Service` struct wraps the cache and consumer unnecessarily:
```go
type Service struct {
    cache     *CatalogCache
    consumer  *authRabbitMQ.CatalogConsumer
    conn      *amqp.Connection
}
```

**Issues**:
1. **Single Responsibility Violation**: Manages both cache and consumer lifecycle
2. **Tight Coupling**: Consumer is created inside the service constructor
3. **Limited Value**: Provides simple getter methods with no business logic
4. **Inconsistent Pattern**: Settings integration doesn't have this wrapper

**Recommendation**: Remove this file and wire cache/consumer directly in main

**Impact**: MEDIUM - Adds unnecessary complexity

---

## 3. Integration Gaps

### ❌ 3.1 Missing Partner HTTP Handler Validation

**Expected in**:
- `authuser/internal/adapters/http/partner/create_partner.go`
- `authuser/internal/adapters/http/partner/update_partner.go`
- `authuser/internal/adapters/http/partner/specialization_management.go`

**Missing Validation**:
```go
// In CreatePartner HTTP handler
if len(request.SpecializationIDs) > 0 {
    if err := h.service.ValidatePartnerSpecializations(ctx, request.SpecializationIDs); err != nil {
        httpx.RespondWithError(w, err, http.StatusBadRequest)
        return
    }
}
```

**Impact**: HIGH - Validation happens too late in the request lifecycle

---

### ❌ 3.2 Missing Error Message Standardization

**Current error messages** (`authuser/internal/application/partner/service.go:116,126`):
```go
return errs.NewInvalidValueErr("specialization ID " + specializationID.String() + " does not exist in catalog")
```

**Issues**:
1. String concatenation instead of `fmt.Sprintf`
2. Generic error type - should be `ErrRepositoryNotFound` or custom catalog error
3. No context about which operation failed

**Recommendation**:
```go
func (s *PartnerService) ValidatePartnerSpecializations(ctx context.Context, specializationIDs []uuid.UUID) error {
    for _, specializationID := range specializationIDs {
        if !s.catalogCache.IsValidCategory(specializationID) {
            return fmt.Errorf("specialization %s not found in catalog: %w",
                specializationID, errs.ErrRepositoryNotFound)
        }
    }
    return nil
}
```

**Impact**: MEDIUM - Poor error messages and incorrect error types

---

## 4. Testing Gaps

### ⚠️ 4.1 Missing Concurrency Tests

**Location**: `authuser/internal/application/catalog/cache_test.go`

**Missing Tests**:
```go
func TestCatalogCache_ConcurrentAccess(t *testing.T) {
    cache := NewCatalogCache()
    ctx := context.Background()

    t.Run("concurrent reads and writes", func(t *testing.T) {
        var wg sync.WaitGroup

        // Spawn 10 writers
        for i := 0; i < 10; i++ {
            wg.Add(1)
            go func(index int) {
                defer wg.Done()
                category := &domain.CachedCategory{
                    ID:     uuid.New(),
                    Name:   fmt.Sprintf("Category %d", index),
                    Status: "published",
                }
                cache.UpsertCategory(ctx, category)
            }(i)
        }

        // Spawn 10 readers
        for i := 0; i < 10; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                cache.ListCategories()
            }()
        }

        wg.Wait()

        // No race conditions should occur
        categories := cache.ListCategories()
        assert.Len(t, categories, 10)
    })
}
```

**Impact**: MEDIUM - Thread safety not verified under concurrent load

---

### ⚠️ 4.2 Missing Status Change Tests

**Location**: `authuser/internal/application/catalog/cache_test.go`

**Missing Tests**:
```go
func TestCatalogCache_StatusChangeRemoval(t *testing.T) {
    cache := NewCatalogCache()
    ctx := context.Background()

    t.Run("should remove category when status changes from published to draft", func(t *testing.T) {
        categoryID := uuid.New()

        // Insert published category
        published := &domain.CachedCategory{
            ID:     categoryID,
            Name:   "Published Category",
            Status: "published",
        }
        err := cache.UpsertCategory(ctx, published)
        require.NoError(t, err)

        // Verify exists
        _, exists := cache.GetCategoryByID(categoryID)
        assert.True(t, exists)

        // Update to draft status
        draft := &domain.CachedCategory{
            ID:     categoryID,
            Name:   "Draft Category",
            Status: "draft",
        }
        err = cache.UpsertCategory(ctx, draft)
        require.NoError(t, err)

        // Verify removed
        _, exists = cache.GetCategoryByID(categoryID)
        assert.False(t, exists)
    })
}
```

**Impact**: MEDIUM - Critical filtering logic not tested

---

### ❌ 4.3 Missing Consumer Integration Tests

**Missing File**: `authuser/internal/adapters/rabbitmq/catalog_consumer_test.go`

**Required Tests**:
1. Test event processing with published items
2. Test event processing with draft items (should not cache)
3. Test status change events (published → draft removal)
4. Test cascade deletion (category deletion removes products)
5. Test malformed JSON handling
6. Test invalid UUID handling
7. Test Nack behavior on processing errors

**Impact**: HIGH - Consumer logic not verified

---

### ❌ 4.4 Missing Partner Service Integration Tests

**Missing Tests** in `authuser/test/integration/partner/`:
```go
func TestCreatePartner_CatalogValidation(t *testing.T) {
    t.Run("should reject partner creation with invalid specializations", func(t *testing.T) {
        // Setup: Catalog cache is empty

        // Attempt to create partner with non-existent specialization
        request := domain.CreatePartnerRequest{
            // ... partner data
            SpecializationIDs: []uuid.UUID{uuid.New()},
        }

        // Should fail with 400 Bad Request
        resp := makeCreatePartnerRequest(t, request)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

        // Verify error message
        body := parseErrorResponse(t, resp)
        assert.Contains(t, body.Error, "not found in catalog")
    })

    t.Run("should accept partner creation with valid specializations", func(t *testing.T) {
        // Setup: Populate catalog cache with published category
        categoryID := uuid.New()
        publishCategoryEvent(t, rabbitmqCh, categoryID, "published")

        // Wait for event processing
        time.Sleep(100 * time.Millisecond)

        // Create partner with valid specialization
        request := domain.CreatePartnerRequest{
            // ... partner data
            SpecializationIDs: []uuid.UUID{categoryID},
        }

        // Should succeed
        resp := makeCreatePartnerRequest(t, request)
        assert.Equal(t, http.StatusCreated, resp.StatusCode)
    })
}
```

**Impact**: HIGH - End-to-end flow not validated

---

## 5. Code Quality Improvements

### ⚠️ 5.1 Map Copy Issues

**Location**: `authuser/internal/application/catalog/cache.go:38-40,54-56`

**Current Code**:
```go
func (c *CatalogCache) GetCategoryByID(id uuid.UUID) (*domain.CachedCategory, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    category, exists := c.categories[id]
    if !exists {
        return nil, false
    }

    // Return a copy to prevent external mutation
    copy := category  // ⚠️ Shallow copy - map fields are still references!
    return &copy, true
}
```

**Problem**: The `Metadata` field is a `map[string]any`, and shallow copying doesn't protect against mutation of the map itself.

**Fix**:
```go
func (c *CatalogCache) GetCategoryByID(id uuid.UUID) (*domain.CachedCategory, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    category, exists := c.categories[id]
    if !exists {
        return nil, false
    }

    // Deep copy including metadata map
    copy := category
    if category.Metadata != nil {
        copy.Metadata = make(map[string]any, len(category.Metadata))
        for k, v := range category.Metadata {
            copy.Metadata[k] = v
        }
    }
    return &copy, true
}
```

**Impact**: MEDIUM - Potential for external mutation of cached data

---

### ⚠️ 5.2 Missing Context Usage

**Location**: Throughout cache methods

**Observation**: Context is accepted but never used:
```go
func (c *CatalogCache) UpsertCategory(ctx context.Context, category *domain.CachedCategory) error {
    // ctx is never used
}
```

**Options**:
1. **Remove context** if not needed (simpler)
2. **Use context** for cancellation (more robust)
3. **Keep context** for future observability/tracing

**Recommendation**: Keep context for future extensibility (logging, metrics, tracing)

**Impact**: LOW - API design consideration

---

### ⚠️ 5.3 Missing Logging in Cache Operations

**Current State**: Only consumer has logging, cache operations are silent

**Recommendation**:
```go
func (c *CatalogCache) UpsertCategory(ctx context.Context, category *domain.CachedCategory) error {
    if category == nil {
        return nil
    }

    if !category.IsPublished() {
        c.mu.Lock()
        if _, exists := c.categories[category.ID]; exists {
            log.Printf("[CatalogCache] Removing unpublished category: id=%s, status=%s",
                category.ID, category.Status)
        }
        delete(c.categories, category.ID)
        c.mu.Unlock()
        return nil
    }

    c.mu.Lock()
    defer c.mu.Unlock()

    action := "updated"
    if _, exists := c.categories[category.ID]; !exists {
        action = "added"
    }

    c.categories[category.ID] = *category
    log.Printf("[CatalogCache] Category %s: id=%s, name=%s",
        action, category.ID, category.Name)

    return nil
}
```

**Impact**: MEDIUM - Observability for debugging

---

## 6. Operational Improvements

### ⚠️ 6.1 Missing Metrics

**Recommended Metrics**:
```go
// In cache.go
type CatalogCache struct {
    mu         sync.RWMutex
    categories map[uuid.UUID]domain.CachedCategory
    products   map[uuid.UUID]domain.CachedProduct

    // Metrics
    categoryCount int64  // Updated atomically
    productCount  int64  // Updated atomically
}

func (c *CatalogCache) GetMetrics() CacheMetrics {
    c.mu.RLock()
    defer c.mu.RUnlock()

    return CacheMetrics{
        CategoryCount: len(c.categories),
        ProductCount:  len(c.products),
    }
}
```

**Impact**: MEDIUM - Production observability

---

### ⚠️ 6.2 Missing Health Check

**Recommendation**: Add catalog cache health indicator
```go
// In health check endpoint
func (h *HealthHandler) checkCatalogCache() HealthStatus {
    metrics := h.catalogCache.GetMetrics()

    return HealthStatus{
        Service:  "catalog_cache",
        Status:   "healthy",
        Metadata: map[string]any{
            "categories_cached": metrics.CategoryCount,
            "products_cached":   metrics.ProductCount,
        },
    }
}
```

**Impact**: LOW - Nice to have for operations

---

## 7. Documentation Improvements

### ⚠️ 7.1 Missing Package Documentation

**Add to `authuser/internal/application/catalog/cache.go`**:
```go
// Package catalog provides an in-memory cache for catalog data consumed from
// the Catalog service via RabbitMQ events.
//
// This implementation follows the same pattern as the Settings cache integration.
// The cache maintains a read-only copy of published categories and products,
// which are used by the Partner service for validation during partner creation
// and updates.
//
// Architecture:
//   - Catalog Service publishes events to RabbitMQ
//   - CatalogConsumer listens to events and updates the in-memory cache
//   - Partner Service validates against the cached data
//
// Thread Safety:
// All cache operations are protected by sync.RWMutex, allowing concurrent reads
// and exclusive writes.
//
// Data Filtering:
// Only items with status="published" are stored in the cache. When an item's
// status changes to "draft" or "archived", it is automatically removed.
package catalog
```

**Impact**: LOW - Developer experience

---

### ⚠️ 7.2 Missing Architecture Diagram

**Add to `CATALOG-EVENT-INTEGRATION-PLAN.md`**:
```markdown
## Architecture Diagram

```
┌─────────────────┐         ┌──────────────────┐
│  Catalog Service│         │  RabbitMQ Broker │
│                 │         │                  │
│  - Categories   │────────▶│  catalog.exchange│
│  - Products     │ Publish │  (topic)         │
└─────────────────┘         └────────┬─────────┘
                                     │
                                     │ Consume
                                     ▼
                            ┌────────────────────┐
                            │ AuthUser Service   │
                            │                    │
                            │  CatalogConsumer   │
                            │         │          │
                            │         ▼          │
                            │  CatalogCache      │
                            │  (In-Memory)       │
                            │         │          │
                            │         ▼          │
                            │  Partner Service   │
                            │  (Validation)      │
                            └────────────────────┘
```
```

**Impact**: LOW - Communication and onboarding

---

## 8. Priority Implementation Plan

### Phase 1: Critical Fixes (Week 1)
**Priority**: CRITICAL

1. ✅ **Fix CreatePartner to use catalog cache validation**
   - Update `create_partner.go` to call `ValidatePartnerSpecializations`
   - Remove repository-based specialization validation
   - Test: Verify partner creation fails with invalid specializations

2. ✅ **Wire catalog cache into service startup**
   - Initialize cache in main.go
   - Start consumer in background goroutine
   - Inject cache into partner service
   - Test: Verify consumer starts and processes events

3. ✅ **Add partner HTTP handler validation**
   - Update create/update handlers
   - Add catalog validation before service calls
   - Return 400 on validation failures
   - Test: Integration tests for validation

**Success Criteria**:
- Partner creation validates against catalog cache
- Consumer processes catalog events on startup
- HTTP API returns proper errors for invalid catalog IDs

---

### Phase 2: Integration & Testing (Week 2)
**Priority**: HIGH

4. ✅ **Implement consumer integration tests**
   - Test event processing with testcontainers
   - Test published vs draft filtering
   - Test cascade deletion
   - Test error handling

5. ✅ **Implement partner integration tests**
   - Test catalog validation in create/update flows
   - Test with empty cache (should fail)
   - Test with populated cache (should succeed)
   - Test status changes (published → draft)

6. ✅ **Add concurrency tests**
   - Test concurrent reads/writes
   - Test with race detector enabled
   - Test under load (100+ goroutines)

**Success Criteria**:
- Full test coverage for consumer logic
- End-to-end partner validation tests passing
- No race conditions detected

---

### Phase 3: Code Quality (Week 3)
**Priority**: MEDIUM

7. ✅ **Fix map copy issues**
   - Deep copy metadata maps
   - Update product copy logic
   - Add tests for mutation protection

8. ✅ **Standardize error handling**
   - Use proper sentinel errors
   - Add context to error messages
   - Update HTTP error responses

9. ✅ **Remove unnecessary service wrapper**
   - Delete `catalog/service.go`
   - Wire cache/consumer directly
   - Update all references

10. ✅ **Resolve exchange type inconsistency**
    - Decide: topic vs direct exchange
    - Update all declarations consistently
    - Document reasoning

**Success Criteria**:
- No shallow copy vulnerabilities
- Consistent error messages
- Simplified service wiring
- Clear exchange type decision

---

### Phase 4: Observability (Week 4)
**Priority**: LOW

11. ✅ **Add logging to cache operations**
    - Log cache hits/misses
    - Log status change removals
    - Log cache size changes

12. ✅ **Add metrics**
    - Track cache size (categories, products)
    - Track validation failures
    - Track consumer processing rate

13. ✅ **Add health check**
    - Report cache status
    - Report consumer status
    - Include cache metrics

**Success Criteria**:
- Operators can monitor cache state
- Alerts possible on cache issues
- Health check includes catalog data

---

### Phase 5: Documentation (Ongoing)
**Priority**: LOW

14. ✅ **Add package documentation**
    - Document cache behavior
    - Document filtering logic
    - Document thread safety

15. ✅ **Update implementation plan**
    - Mark completed phases
    - Add architecture diagram
    - Document design decisions

16. ✅ **Write operational runbook**
    - How to verify cache is working
    - How to debug cache issues
    - How to recover from event loss

**Success Criteria**:
- New developers can understand the system
- Operations team can troubleshoot issues
- Design decisions are documented

---

## 9. Risk Assessment

### High Risk
1. **Service startup not wired** - Implementation not functional
2. **CreatePartner not using cache** - Defeats the purpose
3. **Missing integration tests** - Bugs may go undetected

### Medium Risk
4. **Map copy issues** - Potential data corruption
5. **Missing error standardization** - Poor user experience
6. **No concurrency tests** - Race conditions possible

### Low Risk
7. **Exchange type inconsistency** - Works but confusing
8. **Missing observability** - Hard to debug in production
9. **Missing documentation** - Slower onboarding

---

## 10. Summary

**Overall Grade**: B+ (Good foundation, needs completion)

**Strengths**:
- ✅ Correct architectural pattern (follows Settings integration)
- ✅ Clean in-memory cache implementation
- ✅ Good event contracts and validation
- ✅ Thread-safe cache operations
- ✅ Proper published-only filtering

**Weaknesses**:
- ❌ Service startup not wired
- ❌ Partner creation not using cache validation
- ❌ Missing comprehensive tests
- ⚠️ Some code quality issues (map copying, error messages)
- ⚠️ Minimal observability (logging, metrics)

**Recommendation**: Focus on **Phase 1 (Critical Fixes)** immediately to make the implementation functional. Then proceed with testing and code quality improvements.

---

## Appendix A: Code Examples

### A.1 Corrected CreatePartner Method
```go
func (s *PartnerService) CreatePartner(ctx context.Context, request *domain.CreatePartnerRequest) (*domain.CompletePartnerResponse, error) {
    // Validate request
    if err := request.Valid(ctx); err != nil {
        return nil, errs.ErrInvalidInput
    }

    // ADDED: Validate specializations against catalog cache
    if len(request.SpecializationIDs) > 0 {
        if err := s.ValidatePartnerSpecializations(ctx, request.SpecializationIDs); err != nil {
            return nil, err
        }
    }

    // Check if user with email already exists
    exists, err := s.userRepo.ExistsByEmailHash(ctx, request.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errs.ErrUniqueViolation
    }

    // ... rest of the method (user creation, partner creation)

    // REMOVED: The loop that validates specializations via repository
    // Add specializations directly without re-validation
    for _, specializationID := range request.SpecializationIDs {
        if err := s.partnerRepo.AddPartnerSpecialization(ctx, partnerEncx.ID, specializationID); err != nil {
            return nil, err
        }
    }

    // ... rest of the method
}
```

### A.2 Service Startup Wiring
```go
// In cmd/leviosa/main.go (or equivalent)
func main() {
    // ... existing setup

    // Initialize catalog cache
    catalogCache := catalog.NewCatalogCache()
    log.Println("Catalog cache initialized")

    // Initialize catalog consumer
    catalogConsumer := rabbitmq.NewCatalogConsumer(mqConn, catalogCache)

    // Start catalog consumer in background
    consumerCtx, consumerCancel := context.WithCancel(ctx)
    defer consumerCancel()

    go func() {
        log.Println("Starting catalog consumer...")
        if err := catalogConsumer.Start(consumerCtx); err != nil {
            if err != context.Canceled {
                log.Fatalf("Catalog consumer error: %v", err)
            }
        }
        log.Println("Catalog consumer stopped")
    }()

    // Initialize partner service with catalog cache
    partnerService := partner.New(
        partnerRepo,
        userRepo,
        specializationRepo,
        catalogCache,  // Injected here
        crypto,
        stripe,
    )

    // ... rest of setup
}
```

### A.3 Partner HTTP Handler Validation
```go
// In authuser/internal/adapters/http/partner/create_partner.go
func (h *Handler) CreatePartner(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Decode request
    var request domain.CreatePartnerRequest
    if err := httpx.DecodeJSON(r, &request); err != nil {
        httpx.RespondWithError(w, err, http.StatusBadRequest)
        return
    }

    // ADDED: Validate specializations against catalog cache
    if len(request.SpecializationIDs) > 0 {
        if err := h.service.ValidatePartnerSpecializations(ctx, request.SpecializationIDs); err != nil {
            httpx.RespondWithError(w,
                fmt.Errorf("invalid specializations: %w", err),
                http.StatusBadRequest)
            return
        }
    }

    // Call service
    response, err := h.service.CreatePartner(ctx, &request)
    if err != nil {
        // ... existing error handling
    }

    httpx.RespondWithJSON(w, response, http.StatusCreated)
}
```

---

## Appendix B: Testing Checklist

### Unit Tests
- [x] Cache initialization
- [x] Category upsert (published)
- [x] Category upsert (draft - not stored)
- [x] Category update
- [x] Category deletion
- [x] Category cascade deletion
- [x] Product upsert (published)
- [x] Product upsert (draft - not stored)
- [x] Product update
- [x] Product deletion
- [x] List operations
- [x] Validation methods
- [ ] **Status change removal (published → draft)**
- [ ] **Concurrent access**
- [ ] **Map mutation protection**

### Integration Tests
- [ ] Consumer processes category created
- [ ] Consumer processes category updated
- [ ] Consumer processes category deleted
- [ ] Consumer processes product created
- [ ] Consumer processes product updated
- [ ] Consumer processes product deleted
- [ ] Consumer filters draft items
- [ ] Consumer cascade deletes products
- [ ] Partner creation with valid specializations
- [ ] Partner creation with invalid specializations
- [ ] Partner update with valid products
- [ ] Partner update with invalid products
- [ ] HTTP handler validation

### Performance Tests
- [ ] Cache with 1000+ categories
- [ ] Cache with 10000+ products
- [ ] 100 concurrent readers
- [ ] 100 concurrent writers
- [ ] Mixed read/write load

---

## Appendix C: Metrics to Track

### Cache Metrics
```go
type CacheMetrics struct {
    CategoryCount       int64     `json:"category_count"`
    ProductCount        int64     `json:"product_count"`
    LastUpdateTimestamp time.Time `json:"last_update_timestamp"`
}
```

### Consumer Metrics
```go
type ConsumerMetrics struct {
    EventsProcessed     int64     `json:"events_processed"`
    EventsFailed        int64     `json:"events_failed"`
    LastEventTimestamp  time.Time `json:"last_event_timestamp"`
    ConsumerRunning     bool      `json:"consumer_running"`
}
```

### Validation Metrics
```go
type ValidationMetrics struct {
    SpecializationValidationSuccesses int64 `json:"specialization_validation_successes"`
    SpecializationValidationFailures  int64 `json:"specialization_validation_failures"`
    ProductValidationSuccesses        int64 `json:"product_validation_successes"`
    ProductValidationFailures         int64 `json:"product_validation_failures"`
}
```
