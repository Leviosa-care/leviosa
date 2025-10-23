# Catalog Cache Refactor Plan

## Problem Analysis

We implemented the catalog cache incorrectly by using Postgres storage and HTTP endpoints, when it should follow the Settings pattern with in-memory caching for internal use by the partner service.

## Current Incorrect Implementation

❌ **Postgres-backed catalog cache** with database tables
❌ **HTTP endpoints** exposing catalog data publicly
❌ **Generic catalog operations** not tied to partner service needs
❌ **Repository pattern** with full CRUD operations

## Desired Implementation

✅ **In-memory catalog cache** (like `TokenDurationCache`)
✅ **Private to AuthUser** for internal use only
✅ **Integrated with partner service** for validation
✅ **Event-driven updates** via RabbitMQ consumer
✅ **Only active items** (status = "published") stored in cache

## Reference Pattern: Settings Integration

The catalog cache should follow the exact same pattern as the Settings → AuthUser integration:

1. **Settings Service** publishes updates to RabbitMQ
2. **AuthUser** has a `SettingsConsumer` that listens to these updates
3. **AuthUser** has a `TokenDurationCache` (in-memory cache) that gets updated by the consumer
4. **Session Service** uses this cached data for business logic

## Implementation Tasks

### Phase 1: Remove Incorrect Implementation

#### Task 1.1: Remove HTTP Endpoints
- [ ] Delete `authuser/internal/adapters/http/catalog/` directory completely
  - `endpoints.go`
  - `get_category_by_id.go`
  - `list_categories.go`
  - `get_product_by_id.go`
  - `list_products.go`
  - `list_products_by_category.go`
  - `handler.go`
  - `helpers.go`
  - `routes.go`
- [ ] Remove `authuser/internal/ports/catalog_service.go` file
- [ ] Remove `authuser/internal/application/catalog/service.go` file
- [ ] Delete catalog HTTP integration test `authuser/test/integration/catalog_http_test.go`

#### Task 1.2: Remove Postgres Repository Implementation
- [ ] Delete `authuser/internal/adapters/postgres/catalog_cache/` directory completely
  - `repository.go`
  - `upsert_category.go` and `upsert_category_test.go`
  - `upsert_product.go` and `upsert_product_test.go`
  - `delete_category.go` and `delete_category_test.go`
  - `delete_product.go` and `delete_product_test.go`
  - `get_category_by_id.go` and `get_category_by_id_test.go`
  - `get_product_by_id.go` and `get_product_by_id_test.go`
  - `list_categories.go` and `list_categories_test.go`
  - `list_products.go` and `list_products_test.go`
  - `list_products_by_category.go` and `list_products_by_category_test.go`
  - `main_test.go`
- [ ] Remove `authuser/internal/ports/catalog_cache_repository.go` file
- [ ] Remove catalog domain models:
  - `authuser/internal/domain/catalog_category.go`
  - `authuser/internal/domain/catalog_product.go`
  - `authuser/internal/domain/catalog_dto.go`
- [ ] Remove catalog test helpers `authuser/test/helpers/catalog_cache.go`

#### Task 1.3: Remove Database Migration
- [ ] Delete `core/migrations/20251017161249_authuser_add_catalog_cache_tables.sql`
- [ ] Update `CATALOG-EVENT-INTEGRATION-PLAN.md` to mark removed phases

### Phase 2: Implement Correct In-Memory Cache

#### Task 2.1: Create In-Memory Catalog Cache
- [ ] Create `authuser/internal/application/catalog/cache.go` with:
  - `CatalogCache` struct with `sync.RWMutex` for thread safety
  - `categories` map: `map[uuid.UUID]domain.CachedCategory`
  - `products` map: `map[uuid.UUID]domain.CachedProduct`
  - `NewCatalogCache()` constructor function
  - **Read methods:**
    - `GetCategoryByID(uuid.UUID) (*domain.CachedCategory, bool)` - returns category and existence flag
    - `GetProductByID(uuid.UUID) (*domain.CachedProduct, bool)` - returns product and existence flag
    - `ListCategories() []domain.CachedCategory` - returns all cached categories
    - `ListProducts() []domain.CachedProduct` - returns all cached products
    - `ListProductsByCategory(uuid.UUID) []domain.CachedProduct` - filters products by category
    - `IsValidCategory(uuid.UUID) bool` - checks if category exists
    - `IsValidProduct(uuid.UUID) bool` - checks if product exists
  - **Update methods:**
    - `UpsertCategory(ctx, *domain.CachedCategory) error` - add/update category (only if published)
    - `UpsertProduct(ctx, *domain.CachedProduct) error` - add/update product (only if published)
    - `DeleteCategory(ctx, uuid.UUID) error` - remove category and cascade delete products
    - `DeleteProduct(ctx, uuid.UUID) error` - remove product
  - **Filter logic:** Only store items with `status == "published"`, remove items if status changes

#### Task 2.2: Create Catalog Cache Interface
- [ ] Create `authuser/internal/ports/catalog_cache.go` with:
  - `CatalogCache` interface - read-only methods for services:
    - `GetCategoryByID(id uuid.UUID) (*domain.CachedCategory, bool)`
    - `GetProductByID(id uuid.UUID) (*domain.CachedProduct, bool)`
    - `ListCategories() []domain.CachedCategory`
    - `ListProducts() []domain.CachedProduct`
    - `ListProductsByCategory(categoryID uuid.UUID) []domain.CachedProduct`
    - `IsValidCategory(categoryID uuid.UUID) bool`
    - `IsValidProduct(productID uuid.UUID) bool`
  - `CatalogCacheUpdater` interface - write methods for consumer:
    - `UpsertCategory(ctx context.Context, category *domain.CachedCategory) error`
    - `UpsertProduct(ctx context.Context, product *domain.CachedProduct) error`
    - `DeleteCategory(ctx context.Context, categoryID uuid.UUID) error`
    - `DeleteProduct(ctx context.Context, productID uuid.UUID) error`

#### Task 2.3: Create Simplified Catalog Domain Models
- [ ] Create `authuser/internal/domain/catalog.go` with:
  - `CachedCategory` struct (simplified for in-memory storage):
    - `ID uuid.UUID`
    - `Name string`
    - `Description string`
    - `Status string` (used for filtering)
    - `Metadata map[string]any`
  - `CachedProduct` struct (simplified for in-memory storage):
    - `ID uuid.UUID`
    - `Name string`
    - `Description string`
    - `CategoryID uuid.UUID`
    - `Duration int`
    - `Status string` (used for filtering)
    - `Availability string`
    - `BufferTime int`
    - `CancellationHours int`
    - `StripeProductID string`
    - `Metadata map[string]any`
  - **Remove** database-specific fields: `SyncedAt`, `CreatedAt`, `UpdatedAt`

### Phase 3: Update RabbitMQ Consumer

#### Task 3.1: Modify Catalog Consumer
- [ ] Update `authuser/internal/adapters/rabbitmq/catalog_consumer.go`:
  - Change dependency from `ports.CatalogCacheRepository` to `ports.CatalogCacheUpdater`
  - Update constructor: `NewCatalogConsumer(conn *amqp.Connection, cache ports.CatalogCacheUpdater)`
  - Update `handleProductCreated`:
    - Convert event to simplified domain model
    - Call `cache.UpsertProduct(ctx, product)` (filter handled by cache)
  - Update `handleProductUpdated`:
    - Convert event to simplified domain model
    - Call `cache.UpsertProduct(ctx, product)` (will remove if not published)
  - Update `handleProductDeleted`:
    - Call `cache.DeleteProduct(ctx, productID)`
  - Update `handleCategoryCreated`:
    - Convert event to simplified domain model
    - Call `cache.UpsertCategory(ctx, category)` (filter handled by cache)
  - Update `handleCategoryUpdated`:
    - Convert event to simplified domain model
    - Call `cache.UpsertCategory(ctx, category)` (will remove if not published)
  - Update `handleCategoryDeleted`:
    - Call `cache.DeleteCategory(ctx, categoryID)` (cascades to products)

#### Task 3.2: Update Consumer Setup
- [ ] Verify `authuser/internal/adapters/rabbitmq/setup_catalog_consumer.go` remains correct
- [ ] Update `authuser/internal/adapters/rabbitmq/setup.go` if needed

### Phase 4: Integrate Cache with Partner Service

#### Task 4.1: Update Partner Service Interface
- [ ] Update `authuser/internal/ports/partner_service.go`:
  - Add new validation methods to interface:
    - `ValidatePartnerSpecializations(ctx context.Context, specializationIDs []uuid.UUID) error`
    - `ValidatePartnerProducts(ctx context.Context, productIDs []uuid.UUID) error`

#### Task 4.2: Update Partner Service Implementation
- [ ] Update `authuser/internal/application/partner/service.go`:
  - Add `catalogCache ports.CatalogCache` field to `PartnerService` struct
  - Update `New()` constructor to accept `catalogCache ports.CatalogCache` parameter
  - Implement `ValidatePartnerSpecializations()`:
    - Loop through specialization IDs
    - Call `catalogCache.IsValidCategory(specID)` for each
    - Return error if any not found: `"specialization {id} not found in catalog"`
  - Implement `ValidatePartnerProducts()`:
    - Loop through product IDs
    - Call `catalogCache.IsValidProduct(productID)` for each
    - Return error if any not found: `"product {id} not found in catalog"`
  - Update `CreatePartner()` to validate specializations before creation
  - Update `UpdatePartner()` to validate specializations and products before update
  - Update `AddPartnerSpecialization()` to validate specialization exists

#### Task 4.3: Update Partner HTTP Handlers
- [ ] Update `authuser/internal/adapters/http/partner/create_partner.go`:
  - Call `service.ValidatePartnerSpecializations()` after decoding request
  - Call `service.ValidatePartnerProducts()` if products provided
  - Return `http.StatusBadRequest` if validation fails
- [ ] Update `authuser/internal/adapters/http/partner/update_partner.go`:
  - Call `service.ValidatePartnerSpecializations()` if specializations provided
  - Call `service.ValidatePartnerProducts()` if products provided
  - Return `http.StatusBadRequest` if validation fails
- [ ] Update `authuser/internal/adapters/http/partner/specialization_management.go`:
  - Call `service.ValidatePartnerSpecializations()` before adding specialization
  - Return `http.StatusBadRequest` if validation fails

### Phase 5: Update Service Startup & Wiring

#### Task 5.1: Update Main Service Initialization
- [ ] Locate main service startup file (likely `cmd/leviosa/main.go` or similar)
- [ ] Remove catalog cache repository initialization
- [ ] Add in-memory catalog cache initialization:
  ```go
  catalogCache := catalog.NewCatalogCache()
  ```
- [ ] Update catalog consumer creation:
  ```go
  catalogConsumer := rabbitmq.NewCatalogConsumer(mqConn, catalogCache)
  ```
- [ ] Update partner service creation to include catalog cache:
  ```go
  partnerService := partner.New(
      partnerRepo,
      userRepo,
      specializationRepo,
      crypto,
      stripe,
      catalogCache,
  )
  ```
- [ ] Ensure catalog consumer starts in background goroutine

#### Task 5.2: Update Dependency Injection
- [ ] Update any dependency injection/wiring code to use in-memory cache
- [ ] Remove Postgres catalog cache repository from service container
- [ ] Ensure proper initialization order: cache → consumer → partner service

### Phase 6: Add Testing Infrastructure

#### Task 6.1: Create In-Memory Cache Tests
- [ ] Create `authuser/internal/application/catalog/cache_test.go` with tests for:
  - Cache initialization (should start empty)
  - `UpsertCategory` and `UpsertProduct` with published status filter
  - `GetCategoryByID` and `GetProductByID` return correct data
  - `ListCategories` and `ListProducts` return only published items
  - `ListProductsByCategory` filters correctly
  - `DeleteCategory` cascades to delete products
  - Concurrent access with goroutines (thread safety)
  - Draft/archived items are not stored
  - Status change from published to draft removes item

#### Task 6.2: Update Consumer Tests
- [ ] Create new consumer tests using in-memory cache:
  - Test product created events (published vs draft)
  - Test product updated events (status changes)
  - Test product deleted events
  - Test category created events
  - Test category deleted events (cascade delete products)
  - Test error handling for invalid UUIDs
  - Test that only published items are cached

#### Task 6.3: Update Partner Service Tests
- [ ] Update partner service tests:
  - Test partner creation with valid specializations (in cache)
  - Test partner creation fails with invalid specializations (not in cache)
  - Test partner update with valid products (in cache)
  - Test partner update fails with invalid products (not in cache)
  - Test specialization management validates against cache
  - Test validation error messages

#### Task 6.4: Update HTTP Integration Tests
- [ ] Update partner HTTP integration tests:
  - Test `POST /partners` returns 400 for invalid specializations
  - Test `POST /partners` succeeds with valid specializations
  - Test `PUT /partners/{id}` returns 400 for invalid products/specializations
  - Test `PUT /partners/{id}` succeeds with valid products/specializations
  - Test `POST /partners/{id}/specializations` validates against catalog
  - Test error response contains "not found in catalog" message

### Phase 7: Documentation & Cleanup

#### Task 7.1: Update Implementation Plan
- [ ] Update `CATALOG-EVENT-INTEGRATION-PLAN.md`:
  - Mark Phase 2 (Database Schema) as removed/not needed
  - Mark Phase 5 (Repository Implementation) as removed
  - Mark Phase 8 (Test Helpers) as removed
  - Mark Phase 9 (HTTP Endpoints) as removed
  - Add note explaining in-memory cache approach
  - Update architecture diagram to show in-memory cache
  - Update success criteria

#### Task 7.2: Update Code Documentation
- [ ] Add godoc comments to all cache methods
- [ ] Update partner service documentation with validation behavior
- [ ] Add comments explaining catalog cache pattern follows Settings pattern

### Phase 8: Validation & Testing

#### Task 8.1: Run Unit Tests
- [ ] Run catalog cache unit tests: `go test ./internal/application/catalog/...`
- [ ] Run updated consumer tests: `go test ./internal/adapters/rabbitmq/...`
- [ ] Run updated partner service tests: `go test ./internal/application/partner/...`

#### Task 8.2: Run Integration Tests
- [ ] Run partner HTTP integration tests: `make test-integration`
- [ ] Verify catalog events are consumed correctly
- [ ] Verify partner validation works against cached catalog

#### Task 8.3: Manual Testing
- [ ] Start AuthUser service and verify catalog consumer starts
- [ ] Test partner creation with valid/invalid specializations via API
- [ ] Publish catalog events and verify cache updates in logs
- [ ] Verify service restart maintains empty cache (populated by events)

## Success Criteria

- [ ] All incorrect Postgres/HTTP implementation removed
- [ ] In-memory catalog cache implemented following Settings pattern
- [ ] Catalog consumer updated to use in-memory cache
- [ ] Partner service validates against catalog cache
- [ ] Only published items stored in cache
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Partner API returns proper validation errors
- [ ] Catalog events successfully update cache
- [ ] Code follows existing patterns (Settings integration)
- [ ] Thread-safe concurrent access to cache
- [ ] Service startup works correctly

## Implementation Notes

### Thread Safety
- Use `sync.RWMutex` for all cache operations
- Use `RLock()`/`RUnlock()` for read operations (concurrent reads allowed)
- Use `Lock()`/`Unlock()` for write operations (exclusive access)

### Active Items Only
- Filter events to only cache items with `status == "published"`
- When status changes from "published" to "draft", remove from cache
- When status changes from "draft" to "published", add to cache

### Validation Errors
- Return clear error messages when catalog items not found
- Error format: `"specialization {uuid} not found in catalog"`
- HTTP handlers should return `400 Bad Request` for validation failures

### Event-Driven Only
- No initial data loading from database
- Cache starts empty on service startup
- Rely solely on RabbitMQ events to populate cache
- This is consistent with Settings pattern

### Graceful Degradation
- Service should work even when catalog cache is empty
- Partner operations may fail validation if cache not yet populated
- Consider adding startup delay or initial sync if needed (future enhancement)

## Partner Service Use Cases

### Use Case 1: Partner Creation with Specializations
1. Partner HTTP handler receives create request with specialization IDs
2. Handler calls `partnerService.ValidatePartnerSpecializations(specializationIDs)`
3. Service checks each ID against `catalogCache.IsValidCategory(specID)`
4. If all valid, proceed with partner creation
5. If any invalid, return error: "specialization {id} not found in catalog"

### Use Case 2: Partner Update with Products
1. Partner HTTP handler receives update request with product IDs
2. Handler calls `partnerService.ValidatePartnerProducts(productIDs)`
3. Service checks each ID against `catalogCache.IsValidProduct(productID)`
4. If all valid, proceed with partner update
5. If any invalid, return error: "product {id} not found in catalog"

### Use Case 3: Add Specialization to Partner
1. Partner HTTP handler receives add specialization request
2. Handler calls `partnerService.AddPartnerSpecialization(partnerID, specializationID)`
3. Service calls `catalogCache.IsValidCategory(specializationID)`
4. If valid, proceed to add association
5. If invalid, return error: "specialization {id} not found in catalog"

## Architecture Comparison

### Before (Incorrect)
```
Catalog Service → RabbitMQ → AuthUser Consumer → Postgres Cache → HTTP API → External Clients
                                                ↓
                                         Partner Service (no validation)
```

### After (Correct)
```
Catalog Service → RabbitMQ → AuthUser Consumer → In-Memory Cache → Partner Service (validation)
                                                                  → Partner HTTP Handlers
```

This matches the Settings pattern:
```
Settings Service → RabbitMQ → AuthUser Consumer → In-Memory Cache → Session Service (uses cache)
```

## Files to Create
- `authuser/internal/application/catalog/cache.go`
- `authuser/internal/application/catalog/cache_test.go`
- `authuser/internal/ports/catalog_cache.go`
- `authuser/internal/domain/catalog.go`

## Files to Modify
- `authuser/internal/adapters/rabbitmq/catalog_consumer.go`
- `authuser/internal/ports/partner_service.go`
- `authuser/internal/application/partner/service.go`
- `authuser/internal/adapters/http/partner/create_partner.go`
- `authuser/internal/adapters/http/partner/update_partner.go`
- `authuser/internal/adapters/http/partner/specialization_management.go`
- Main service startup file (e.g., `cmd/leviosa/main.go`)
- `CATALOG-EVENT-INTEGRATION-PLAN.md`

## Files to Delete
- `authuser/internal/adapters/http/catalog/` (entire directory)
- `authuser/internal/adapters/postgres/catalog_cache/` (entire directory)
- `authuser/internal/ports/catalog_service.go`
- `authuser/internal/ports/catalog_cache_repository.go`
- `authuser/internal/application/catalog/service.go`
- `authuser/internal/domain/catalog_category.go`
- `authuser/internal/domain/catalog_product.go`
- `authuser/internal/domain/catalog_dto.go`
- `authuser/test/helpers/catalog_cache.go`
- `authuser/test/integration/catalog_http_test.go`
- `core/migrations/20251017161249_authuser_add_catalog_cache_tables.sql`

## Estimated Timeline
- Phase 1 (Removal): 30 minutes
- Phase 2 (In-Memory Cache): 1 hour
- Phase 3 (Consumer Update): 30 minutes
- Phase 4 (Partner Integration): 1 hour
- Phase 5 (Service Startup): 30 minutes
- Phase 6 (Testing): 2 hours
- Phase 7 (Documentation): 30 minutes
- Phase 8 (Validation): 1 hour

**Total: ~7 hours**
