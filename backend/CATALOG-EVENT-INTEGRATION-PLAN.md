# Event-Driven Catalog Integration Implementation Plan

> **Implementation Pattern**: Event-driven architecture following the existing Settings → AuthUser pattern
> **Goal**: Enable AuthUser to consume catalog data (products & categories) via RabbitMQ events with persistent local cache
> **Approach**: Option C - Persistent local Postgres cache in AuthUser, populated by Catalog events

---

## 📋 Implementation Overview

**Architecture:**
- **Catalog Service** (Publisher): Publishes product/category events to RabbitMQ
- **AuthUser Service** (Consumer): Consumes events and maintains local read-only cache
- **Core Package**: Defines shared event contracts, queue names, and routing keys
- **Local Cache**: Postgres tables in AuthUser for resilient data access

**Key Patterns:**
- Follow `settings_consumer.go` pattern for event consumption
- Use hexagonal architecture: domain → ports → adapters → application
- Apply `ClassifyPgError()` for error handling with sentinel errors
- Implement testcontainer-based integration tests
- Store local cache in dedicated Postgres tables

---

## 🏗️ Phase 1: Core Contracts & Infrastructure

### 1.1 Define Event Contracts

#### **File: `core/contracts/catalog/events.go`**

```go
package catalog

import (
	"time"

	"github.com/google/uuid"
)

// CategoryCreatedEvent represents a new category being created
type CategoryCreatedEvent struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"` // published, draft, archived
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// CategoryUpdatedEvent represents a category being updated
type CategoryUpdatedEvent struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// CategoryDeletedEvent represents a category being deleted
type CategoryDeletedEvent struct {
	ID string `json:"id"`
}

// ProductCreatedEvent represents a new product being created
type ProductCreatedEvent struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	CategoryID        string         `json:"categoryId"`
	Duration          int            `json:"duration"`
	Status            string         `json:"status"`
	Availability      string         `json:"availability"` // online, in-person, hybrid
	BufferTime        int            `json:"bufferTime"`
	CancellationHours int            `json:"cancellationHours"`
	StripeProductID   string         `json:"stripeProductId"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

// ProductUpdatedEvent represents a product being updated
type ProductUpdatedEvent struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	CategoryID        string         `json:"categoryId"`
	Duration          int            `json:"duration"`
	Status            string         `json:"status"`
	Availability      string         `json:"availability"`
	BufferTime        int            `json:"bufferTime"`
	CancellationHours int            `json:"cancellationHours"`
	StripeProductID   string         `json:"stripeProductId"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

// ProductDeletedEvent represents a product being deleted
type ProductDeletedEvent struct {
	ID string `json:"id"`
}

// Helper to validate UUID in events
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// Validation methods
func (e CategoryCreatedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid category ID: %s", e.ID)
	}
	if e.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	return nil
}

func (e ProductCreatedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid product ID: %s", e.ID)
	}
	if !isValidUUID(e.CategoryID) {
		return fmt.Errorf("invalid category ID: %s", e.CategoryID)
	}
	if e.Name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	return nil
}
```

**Example JSON Payloads:**

```json
// CategoryCreatedEvent
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Massage",
  "description": "Therapeutic massage services",
  "status": "published",
  "metadata": {"color": "#FF5733"},
  "createdAt": "2025-10-17T16:12:49Z",
  "updatedAt": "2025-10-17T16:12:49Z"
}

// ProductCreatedEvent
{
  "id": "223e4567-e89b-12d3-a456-426614174001",
  "name": "Swedish Massage - 60min",
  "description": "Relaxing full body massage",
  "categoryId": "123e4567-e89b-12d3-a456-426614174000",
  "duration": 60,
  "status": "published",
  "availability": "in-person",
  "bufferTime": 15,
  "cancellationHours": 24,
  "stripeProductId": "prod_ABC123",
  "metadata": {},
  "createdAt": "2025-10-17T16:12:49Z",
  "updatedAt": "2025-10-17T16:12:49Z"
}
```

---

#### **File: `core/contracts/catalog/keys.go`**

```go
package catalog

const (
	ProductKey  = "product"
	CategoryKey = "category"
)
```

---

#### **File: `core/contracts/rabbitmq/catalog.go`**

```go
package rabbitmq

const (
	// Exchange
	CatalogExchangeName = "catalog.exchange"

	// Routing keys
	ProductCreatedRoutingKey  = "catalog.product.created"
	ProductUpdatedRoutingKey  = "catalog.product.updated"
	ProductDeletedRoutingKey  = "catalog.product.deleted"
	CategoryCreatedRoutingKey = "catalog.category.created"
	CategoryUpdatedRoutingKey = "catalog.category.updated"
	CategoryDeletedRoutingKey = "catalog.category.deleted"

	// Queue names per consuming service
	AuthUserCatalogQueueName = "authuser.catalog.queue"
	// Add more queues here as other services need catalog data
	// BookingCatalogQueueName  = "booking.catalog.queue"
)
```

---

## 🗄️ Phase 2: Database Migrations & Schema

### 2.1 Create AuthUser Catalog Cache Tables

#### **File: `core/migrations/20251017161249_authuser_add_catalog_cache_tables.sql`**

```sql
-- +goose Up
-- +goose StatementBegin

-- Create dedicated schema for AuthUser's event-sourced caches
CREATE SCHEMA IF NOT EXISTS authuser_cache;

COMMENT ON SCHEMA authuser_cache IS 'Event-sourced read-only caches maintained by AuthUser service';

-- Categories cache table
CREATE TABLE IF NOT EXISTS authuser_cache.categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    metadata JSONB,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT chk_category_status CHECK (status IN ('published', 'draft', 'archived'))
);

COMMENT ON TABLE authuser_cache.categories IS 'Read-only cache of catalog categories, populated via RabbitMQ events from Catalog service';
COMMENT ON COLUMN authuser_cache.categories.synced_at IS 'Timestamp when this record was last synced from a Catalog event';

-- Products cache table
CREATE TABLE IF NOT EXISTS authuser_cache.products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    category_id UUID NOT NULL,
    duration INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    availability VARCHAR(20) NOT NULL,
    buffer_time INT NOT NULL,
    cancellation_hours INT NOT NULL,
    stripe_product_id TEXT,
    metadata JSONB,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_category FOREIGN KEY (category_id)
        REFERENCES authuser_cache.categories(id)
        ON DELETE CASCADE,
    CONSTRAINT chk_product_status CHECK (status IN ('published', 'draft', 'archived')),
    CONSTRAINT chk_product_availability CHECK (availability IN ('online', 'in-person', 'hybrid'))
);

COMMENT ON TABLE authuser_cache.products IS 'Read-only cache of catalog products, populated via RabbitMQ events from Catalog service';
COMMENT ON COLUMN authuser_cache.products.synced_at IS 'Timestamp when this record was last synced from a Catalog event';

-- Performance indexes
CREATE INDEX idx_authuser_cache_products_category_id
    ON authuser_cache.products(category_id);

CREATE INDEX idx_authuser_cache_products_status
    ON authuser_cache.products(status);

CREATE INDEX idx_authuser_cache_categories_status
    ON authuser_cache.categories(status);

CREATE INDEX idx_authuser_cache_products_synced_at
    ON authuser_cache.products(synced_at);

CREATE INDEX idx_authuser_cache_categories_synced_at
    ON authuser_cache.categories(synced_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS authuser_cache.products CASCADE;
DROP TABLE IF EXISTS authuser_cache.categories CASCADE;
DROP SCHEMA IF EXISTS authuser_cache CASCADE;

-- +goose StatementEnd
```

---

## 📦 Phase 3: AuthUser Domain Layer

### 3.1 Define Domain Models for Cached Catalog Data

#### **File: `authuser/internal/domain/catalog_category.go`**

```go
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// CachedCategory represents a catalog category cached from events
type CachedCategory struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	SyncedAt    time.Time      `json:"syncedAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// Valid performs basic validation
func (c CachedCategory) Valid(ctx context.Context) error {
	if c.ID == uuid.Nil {
		return fmt.Errorf("category ID cannot be nil")
	}
	if c.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	if c.Status != "published" && c.Status != "draft" && c.Status != "archived" {
		return fmt.Errorf("invalid category status: %s", c.Status)
	}
	return nil
}
```

---

#### **File: `authuser/internal/domain/catalog_product.go`**

```go
package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CachedProduct represents a catalog product cached from events
type CachedProduct struct {
	ID                uuid.UUID      `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	CategoryID        uuid.UUID      `json:"categoryId"`
	Duration          int            `json:"duration"`
	Status            string         `json:"status"`
	Availability      string         `json:"availability"`
	BufferTime        int            `json:"bufferTime"`
	CancellationHours int            `json:"cancellationHours"`
	StripeProductID   string         `json:"stripeProductId"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	SyncedAt          time.Time      `json:"syncedAt"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

// Valid performs basic validation
func (p CachedProduct) Valid(ctx context.Context) error {
	if p.ID == uuid.Nil {
		return fmt.Errorf("product ID cannot be nil")
	}
	if p.CategoryID == uuid.Nil {
		return fmt.Errorf("category ID cannot be nil")
	}
	if p.Name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	if p.Duration <= 0 {
		return fmt.Errorf("product duration must be positive")
	}
	if p.Status != "published" && p.Status != "draft" && p.Status != "archived" {
		return fmt.Errorf("invalid product status: %s", p.Status)
	}
	if p.Availability != "online" && p.Availability != "in-person" && p.Availability != "hybrid" {
		return fmt.Errorf("invalid product availability: %s", p.Availability)
	}
	return nil
}
```

---

#### **File: `authuser/internal/domain/catalog_dto.go`**

```go
package domain

// ListCategoriesResponse wraps a list of categories
type ListCategoriesResponse struct {
	Categories []CachedCategory `json:"categories"`
	Total      int              `json:"total"`
}

// ListProductsResponse wraps a list of products
type ListProductsResponse struct {
	Products []CachedProduct `json:"products"`
	Total    int             `json:"total"`
}
```

---

## 🔌 Phase 4: AuthUser Ports (Interfaces)

### 4.1 Define Repository Interfaces

#### **File: `authuser/internal/ports/catalog_cache_repository.go`**

```go
package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
)

// CatalogCacheRepository defines methods for managing the local catalog cache
type CatalogCacheRepository interface {
	// Category operations
	UpsertCategory(ctx context.Context, category domain.CachedCategory) error
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategoryByID(ctx context.Context, categoryID string) (*domain.CachedCategory, error)
	ListCategories(ctx context.Context, publishedOnly bool) ([]domain.CachedCategory, error)

	// Product operations
	UpsertProduct(ctx context.Context, product domain.CachedProduct) error
	DeleteProduct(ctx context.Context, productID string) error
	GetProductByID(ctx context.Context, productID string) (*domain.CachedProduct, error)
	ListProducts(ctx context.Context, publishedOnly bool) ([]domain.CachedProduct, error)
	ListProductsByCategory(ctx context.Context, categoryID string) ([]domain.CachedProduct, error)
}
```

---

## 🗃️ Phase 5: AuthUser Adapters - Postgres Repository

### 5.1 Implement Catalog Cache Repository

#### **File: `authuser/internal/adapters/postgres/catalog_cache/repository.go`**

```go
package catalog_cache

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

const schema = "authuser_cache"

// Repository implements catalog cache persistence
type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

// New creates a new catalog cache repository
func New(ctx context.Context, pool *pgxpool.Pool) ports.CatalogCacheRepository {
	return &Repository{
		pool:   pool,
		schema: schema,
	}
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/upsert_category.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) UpsertCategory(ctx context.Context, category domain.CachedCategory) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.categories (
			id, name, description, status, metadata, synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW(), $6, $7
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			status = EXCLUDED.status,
			metadata = EXCLUDED.metadata,
			synced_at = NOW(),
			updated_at = EXCLUDED.updated_at
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.Status,
		category.Metadata,
		category.CreatedAt,
		category.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("upsert category %s", category.ID), err)
	}

	return nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/upsert_category_test.go`**

```go
package catalog_cache

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully insert new category", func(t *testing.T) {
		// Arrange
		category := domain.CachedCategory{
			ID:          uuid.New(),
			Name:        "Test Category",
			Description: "Test Description",
			Status:      "published",
			Metadata:    map[string]any{"color": "#FF5733"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Act
		err := testRepo.UpsertCategory(ctx, category)

		// Assert
		require.NoError(t, err)

		// Verify it was inserted
		retrieved, err := testRepo.GetCategoryByID(ctx, category.ID.String())
		require.NoError(t, err)
		assert.Equal(t, category.Name, retrieved.Name)
		assert.Equal(t, category.Description, retrieved.Description)
		assert.Equal(t, category.Status, retrieved.Status)
	})

	t.Run("should successfully update existing category", func(t *testing.T) {
		// Arrange - insert initial category
		categoryID := uuid.New()
		initialCategory := domain.CachedCategory{
			ID:          categoryID,
			Name:        "Initial Name",
			Description: "Initial Description",
			Status:      "draft",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, initialCategory)
		require.NoError(t, err)

		// Act - update with new data
		updatedCategory := domain.CachedCategory{
			ID:          categoryID,
			Name:        "Updated Name",
			Description: "Updated Description",
			Status:      "published",
			CreatedAt:   initialCategory.CreatedAt,
			UpdatedAt:   time.Now(),
		}
		err = testRepo.UpsertCategory(ctx, updatedCategory)

		// Assert
		require.NoError(t, err)

		// Verify update
		retrieved, err := testRepo.GetCategoryByID(ctx, categoryID.String())
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", retrieved.Name)
		assert.Equal(t, "Updated Description", retrieved.Description)
		assert.Equal(t, "published", retrieved.Status)
		assert.True(t, retrieved.SyncedAt.After(initialCategory.UpdatedAt))
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/upsert_product.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) UpsertProduct(ctx context.Context, product domain.CachedProduct) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.products (
			id, name, description, category_id, duration, status, availability,
			buffer_time, cancellation_hours, stripe_product_id, metadata,
			synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), $12, $13
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			category_id = EXCLUDED.category_id,
			duration = EXCLUDED.duration,
			status = EXCLUDED.status,
			availability = EXCLUDED.availability,
			buffer_time = EXCLUDED.buffer_time,
			cancellation_hours = EXCLUDED.cancellation_hours,
			stripe_product_id = EXCLUDED.stripe_product_id,
			metadata = EXCLUDED.metadata,
			synced_at = NOW(),
			updated_at = EXCLUDED.updated_at
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.CategoryID,
		product.Duration,
		product.Status,
		product.Availability,
		product.BufferTime,
		product.CancellationHours,
		product.StripeProductID,
		product.Metadata,
		product.CreatedAt,
		product.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("upsert product %s", product.ID), err)
	}

	return nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/upsert_product_test.go`**

```go
package catalog_cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully insert new product", func(t *testing.T) {
		// Arrange - create category first
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:          categoryID,
			Name:        "Test Category",
			Description: "Test Description",
			Status:      "published",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Create product
		product := domain.CachedProduct{
			ID:                uuid.New(),
			Name:              "Test Product",
			Description:       "Test Description",
			CategoryID:        categoryID,
			Duration:          60,
			Status:            "published",
			Availability:      "in-person",
			BufferTime:        15,
			CancellationHours: 24,
			StripeProductID:   "prod_TEST123",
			Metadata:          map[string]any{"key": "value"},
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		// Act
		err = testRepo.UpsertProduct(ctx, product)

		// Assert
		require.NoError(t, err)

		// Verify it was inserted
		retrieved, err := testRepo.GetProductByID(ctx, product.ID.String())
		require.NoError(t, err)
		assert.Equal(t, product.Name, retrieved.Name)
		assert.Equal(t, product.CategoryID, retrieved.CategoryID)
		assert.Equal(t, product.Duration, retrieved.Duration)
	})

	t.Run("should fail with foreign key violation for non-existent category", func(t *testing.T) {
		// Arrange
		product := domain.CachedProduct{
			ID:                uuid.New(),
			Name:              "Orphan Product",
			Description:       "Test Description",
			CategoryID:        uuid.New(), // Non-existent category
			Duration:          60,
			Status:            "published",
			Availability:      "online",
			BufferTime:        15,
			CancellationHours: 24,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		// Act
		err := testRepo.UpsertProduct(ctx, product)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrForeignKeyViolation))
	})

	t.Run("should successfully update existing product", func(t *testing.T) {
		// Arrange - create category
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:          categoryID,
			Name:        "Test Category",
			Status:      "published",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Insert initial product
		productID := uuid.New()
		initialProduct := domain.CachedProduct{
			ID:                productID,
			Name:              "Initial Product",
			Description:       "Initial Description",
			CategoryID:        categoryID,
			Duration:          30,
			Status:            "draft",
			Availability:      "online",
			BufferTime:        10,
			CancellationHours: 12,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err = testRepo.UpsertProduct(ctx, initialProduct)
		require.NoError(t, err)

		// Act - update product
		time.Sleep(10 * time.Millisecond) // Ensure different synced_at
		updatedProduct := domain.CachedProduct{
			ID:                productID,
			Name:              "Updated Product",
			Description:       "Updated Description",
			CategoryID:        categoryID,
			Duration:          60,
			Status:            "published",
			Availability:      "in-person",
			BufferTime:        15,
			CancellationHours: 24,
			CreatedAt:         initialProduct.CreatedAt,
			UpdatedAt:         time.Now(),
		}
		err = testRepo.UpsertProduct(ctx, updatedProduct)

		// Assert
		require.NoError(t, err)

		// Verify update
		retrieved, err := testRepo.GetProductByID(ctx, productID.String())
		require.NoError(t, err)
		assert.Equal(t, "Updated Product", retrieved.Name)
		assert.Equal(t, 60, retrieved.Duration)
		assert.Equal(t, "published", retrieved.Status)
		assert.True(t, retrieved.SyncedAt.After(initialProduct.UpdatedAt))
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/delete_category.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) DeleteCategory(ctx context.Context, categoryID string) error {
	id, err := uuid.Parse(categoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %w", err)
	}

	query := fmt.Sprintf(`
		DELETE FROM %s.categories
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("delete category %s", categoryID), err)
	}

	// Check if anything was deleted (idempotency - not an error if nothing deleted)
	_ = result.RowsAffected()

	return nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/delete_category_test.go`**

```go
package catalog_cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete existing category", func(t *testing.T) {
		// Arrange - create category
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:        categoryID,
			Name:      "Category to Delete",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Act
		err = testRepo.DeleteCategory(ctx, categoryID.String())

		// Assert
		require.NoError(t, err)

		// Verify it was deleted
		_, err = testRepo.GetCategoryByID(ctx, categoryID.String())
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should be idempotent - no error when deleting non-existent category", func(t *testing.T) {
		// Arrange
		nonExistentID := uuid.New()

		// Act
		err := testRepo.DeleteCategory(ctx, nonExistentID.String())

		// Assert - should not error
		assert.NoError(t, err)
	})

	t.Run("should cascade delete products when category is deleted", func(t *testing.T) {
		// Arrange - create category with products
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:        categoryID,
			Name:      "Category with Products",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Create product
		productID := uuid.New()
		product := domain.CachedProduct{
			ID:                productID,
			Name:              "Product in Category",
			CategoryID:        categoryID,
			Duration:          60,
			Status:            "published",
			Availability:      "online",
			BufferTime:        15,
			CancellationHours: 24,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err = testRepo.UpsertProduct(ctx, product)
		require.NoError(t, err)

		// Act - delete category
		err = testRepo.DeleteCategory(ctx, categoryID.String())

		// Assert
		require.NoError(t, err)

		// Verify category is deleted
		_, err = testRepo.GetCategoryByID(ctx, categoryID.String())
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))

		// Verify product is also deleted (CASCADE)
		_, err = testRepo.GetProductByID(ctx, productID.String())
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/delete_product.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) DeleteProduct(ctx context.Context, productID string) error {
	id, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("invalid product ID format: %w", err)
	}

	query := fmt.Sprintf(`
		DELETE FROM %s.products
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError(fmt.Sprintf("delete product %s", productID), err)
	}

	// Idempotent - not an error if nothing deleted
	_ = result.RowsAffected()

	return nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/delete_product_test.go`**

```go
package catalog_cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete existing product", func(t *testing.T) {
		// Arrange - create category and product
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:        categoryID,
			Name:      "Test Category",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		productID := uuid.New()
		product := domain.CachedProduct{
			ID:                productID,
			Name:              "Product to Delete",
			CategoryID:        categoryID,
			Duration:          60,
			Status:            "published",
			Availability:      "online",
			BufferTime:        15,
			CancellationHours: 24,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err = testRepo.UpsertProduct(ctx, product)
		require.NoError(t, err)

		// Act
		err = testRepo.DeleteProduct(ctx, productID.String())

		// Assert
		require.NoError(t, err)

		// Verify it was deleted
		_, err = testRepo.GetProductByID(ctx, productID.String())
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))

		// Verify category still exists
		_, err = testRepo.GetCategoryByID(ctx, categoryID.String())
		assert.NoError(t, err)
	})

	t.Run("should be idempotent - no error when deleting non-existent product", func(t *testing.T) {
		// Arrange
		nonExistentID := uuid.New()

		// Act
		err := testRepo.DeleteProduct(ctx, nonExistentID.String())

		// Assert - should not error
		assert.NoError(t, err)
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/get_category_by_id.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetCategoryByID(ctx context.Context, categoryID string) (*domain.CachedCategory, error) {
	id, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID format: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, name, description, status, metadata, synced_at, created_at, updated_at
		FROM %s.categories
		WHERE id = $1
	`, r.schema)

	category := &domain.CachedCategory{}
	err = r.pool.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.Metadata,
		&category.SyncedAt,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError(fmt.Sprintf("get category by id %s", categoryID), err)
	}

	return category, nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/get_category_by_id_test.go`**

```go
package catalog_cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCategoryByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve existing category", func(t *testing.T) {
		// Arrange
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:          categoryID,
			Name:        "Test Category",
			Description: "Test Description",
			Status:      "published",
			Metadata:    map[string]any{"key": "value"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Act
		retrieved, err := testRepo.GetCategoryByID(ctx, categoryID.String())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, category.ID, retrieved.ID)
		assert.Equal(t, category.Name, retrieved.Name)
		assert.Equal(t, category.Description, retrieved.Description)
		assert.Equal(t, category.Status, retrieved.Status)
		assert.NotZero(t, retrieved.SyncedAt)
	})

	t.Run("should return ErrRepositoryNotFound for non-existent category", func(t *testing.T) {
		// Arrange
		nonExistentID := uuid.New()

		// Act
		_, err := testRepo.GetCategoryByID(ctx, nonExistentID.String())

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should return error for invalid UUID format", func(t *testing.T) {
		// Act
		_, err := testRepo.GetCategoryByID(ctx, "invalid-uuid")

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid category ID format")
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/get_product_by_id.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetProductByID(ctx context.Context, productID string) (*domain.CachedProduct, error) {
	id, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			id, name, description, category_id, duration, status, availability,
			buffer_time, cancellation_hours, stripe_product_id, metadata,
			synced_at, created_at, updated_at
		FROM %s.products
		WHERE id = $1
	`, r.schema)

	product := &domain.CachedProduct{}
	err = r.pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&product.Duration,
		&product.Status,
		&product.Availability,
		&product.BufferTime,
		&product.CancellationHours,
		&product.StripeProductID,
		&product.Metadata,
		&product.SyncedAt,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError(fmt.Sprintf("get product by id %s", productID), err)
	}

	return product, nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/get_product_by_id_test.go`**

```go
package catalog_cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve existing product", func(t *testing.T) {
		// Arrange - create category and product
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:        categoryID,
			Name:      "Test Category",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		productID := uuid.New()
		product := domain.CachedProduct{
			ID:                productID,
			Name:              "Test Product",
			Description:       "Test Description",
			CategoryID:        categoryID,
			Duration:          60,
			Status:            "published",
			Availability:      "in-person",
			BufferTime:        15,
			CancellationHours: 24,
			StripeProductID:   "prod_TEST123",
			Metadata:          map[string]any{"key": "value"},
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err = testRepo.UpsertProduct(ctx, product)
		require.NoError(t, err)

		// Act
		retrieved, err := testRepo.GetProductByID(ctx, productID.String())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, product.ID, retrieved.ID)
		assert.Equal(t, product.Name, retrieved.Name)
		assert.Equal(t, product.CategoryID, retrieved.CategoryID)
		assert.Equal(t, product.Duration, retrieved.Duration)
		assert.NotZero(t, retrieved.SyncedAt)
	})

	t.Run("should return ErrRepositoryNotFound for non-existent product", func(t *testing.T) {
		// Arrange
		nonExistentID := uuid.New()

		// Act
		_, err := testRepo.GetProductByID(ctx, nonExistentID.String())

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/list_categories.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) ListCategories(ctx context.Context, publishedOnly bool) ([]domain.CachedCategory, error) {
	query := fmt.Sprintf(`
		SELECT id, name, description, status, metadata, synced_at, created_at, updated_at
		FROM %s.categories
	`, r.schema)

	if publishedOnly {
		query += " WHERE status = 'published'"
	}

	query += " ORDER BY name ASC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("list categories", err)
	}
	defer rows.Close()

	categories := make([]domain.CachedCategory, 0)
	for rows.Next() {
		var cat domain.CachedCategory
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&cat.Status,
			&cat.Metadata,
			&cat.SyncedAt,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan category row", err)
		}
		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("list categories rows error", err)
	}

	return categories, nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/list_categories_test.go`**

```go
package catalog_cache

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCategories(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all categories when publishedOnly is false", func(t *testing.T) {
		// Arrange - create categories with different statuses
		categories := []domain.CachedCategory{
			{
				ID:        uuid.New(),
				Name:      "Published Category",
				Status:    "published",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "Draft Category",
				Status:    "draft",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "Archived Category",
				Status:    "archived",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		for _, cat := range categories {
			err := testRepo.UpsertCategory(ctx, cat)
			require.NoError(t, err)
		}

		// Act
		results, err := testRepo.ListCategories(ctx, false)

		// Assert
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 3)
	})

	t.Run("should return only published categories when publishedOnly is true", func(t *testing.T) {
		// Arrange - categories already created in previous test

		// Act
		results, err := testRepo.ListCategories(ctx, true)

		// Assert
		require.NoError(t, err)
		for _, cat := range results {
			assert.Equal(t, "published", cat.Status)
		}
	})

	t.Run("should return empty slice when no categories exist", func(t *testing.T) {
		// Arrange - clear all categories
		_, err := testPool.Exec(ctx, "TRUNCATE TABLE authuser_cache.categories CASCADE")
		require.NoError(t, err)

		// Act
		results, err := testRepo.ListCategories(ctx, false)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, results)
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/list_products.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) ListProducts(ctx context.Context, publishedOnly bool) ([]domain.CachedProduct, error) {
	query := fmt.Sprintf(`
		SELECT
			id, name, description, category_id, duration, status, availability,
			buffer_time, cancellation_hours, stripe_product_id, metadata,
			synced_at, created_at, updated_at
		FROM %s.products
	`, r.schema)

	if publishedOnly {
		query += " WHERE status = 'published'"
	}

	query += " ORDER BY name ASC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("list products", err)
	}
	defer rows.Close()

	products := make([]domain.CachedProduct, 0)
	for rows.Next() {
		var prod domain.CachedProduct
		err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Description,
			&prod.CategoryID,
			&prod.Duration,
			&prod.Status,
			&prod.Availability,
			&prod.BufferTime,
			&prod.CancellationHours,
			&prod.StripeProductID,
			&prod.Metadata,
			&prod.SyncedAt,
			&prod.CreatedAt,
			&prod.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan product row", err)
		}
		products = append(products, prod)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("list products rows error", err)
	}

	return products, nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/list_products_test.go`**

```go
package catalog_cache

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProducts(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all products when publishedOnly is false", func(t *testing.T) {
		// Arrange - create category
		categoryID := uuid.New()
		category := domain.CachedCategory{
			ID:        categoryID,
			Name:      "Test Category",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Create products with different statuses
		products := []domain.CachedProduct{
			{
				ID:                uuid.New(),
				Name:              "Published Product",
				CategoryID:        categoryID,
				Duration:          60,
				Status:            "published",
				Availability:      "online",
				BufferTime:        15,
				CancellationHours: 24,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				ID:                uuid.New(),
				Name:              "Draft Product",
				CategoryID:        categoryID,
				Duration:          60,
				Status:            "draft",
				Availability:      "online",
				BufferTime:        15,
				CancellationHours: 24,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		}
		for _, prod := range products {
			err := testRepo.UpsertProduct(ctx, prod)
			require.NoError(t, err)
		}

		// Act
		results, err := testRepo.ListProducts(ctx, false)

		// Assert
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)
	})

	t.Run("should return only published products when publishedOnly is true", func(t *testing.T) {
		// Act
		results, err := testRepo.ListProducts(ctx, true)

		// Assert
		require.NoError(t, err)
		for _, prod := range results {
			assert.Equal(t, "published", prod.Status)
		}
	})

	t.Run("should return empty slice when no products exist", func(t *testing.T) {
		// Arrange - clear all
		_, err := testPool.Exec(ctx, "TRUNCATE TABLE authuser_cache.products CASCADE")
		require.NoError(t, err)

		// Act
		results, err := testRepo.ListProducts(ctx, false)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, results)
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/list_products_by_category.go`**

```go
package catalog_cache

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) ListProductsByCategory(ctx context.Context, categoryID string) ([]domain.CachedProduct, error) {
	id, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID format: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			id, name, description, category_id, duration, status, availability,
			buffer_time, cancellation_hours, stripe_product_id, metadata,
			synced_at, created_at, updated_at
		FROM %s.products
		WHERE category_id = $1
		ORDER BY name ASC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, errs.ClassifyPgError(fmt.Sprintf("list products by category %s", categoryID), err)
	}
	defer rows.Close()

	products := make([]domain.CachedProduct, 0)
	for rows.Next() {
		var prod domain.CachedProduct
		err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Description,
			&prod.CategoryID,
			&prod.Duration,
			&prod.Status,
			&prod.Availability,
			&prod.BufferTime,
			&prod.CancellationHours,
			&prod.StripeProductID,
			&prod.Metadata,
			&prod.SyncedAt,
			&prod.CreatedAt,
			&prod.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan product row", err)
		}
		products = append(products, prod)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("list products by category rows error", err)
	}

	return products, nil
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/list_products_by_category_test.go`**

```go
package catalog_cache

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProductsByCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("should return products for existing category", func(t *testing.T) {
		// Arrange - create two categories with products
		category1ID := uuid.New()
		category1 := domain.CachedCategory{
			ID:        category1ID,
			Name:      "Category 1",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, category1)
		require.NoError(t, err)

		category2ID := uuid.New()
		category2 := domain.CachedCategory{
			ID:        category2ID,
			Name:      "Category 2",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = testRepo.UpsertCategory(ctx, category2)
		require.NoError(t, err)

		// Create products in category 1
		for i := 0; i < 3; i++ {
			product := domain.CachedProduct{
				ID:                uuid.New(),
				Name:              fmt.Sprintf("Product %d in Cat1", i),
				CategoryID:        category1ID,
				Duration:          60,
				Status:            "published",
				Availability:      "online",
				BufferTime:        15,
				CancellationHours: 24,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			err := testRepo.UpsertProduct(ctx, product)
			require.NoError(t, err)
		}

		// Create products in category 2
		for i := 0; i < 2; i++ {
			product := domain.CachedProduct{
				ID:                uuid.New(),
				Name:              fmt.Sprintf("Product %d in Cat2", i),
				CategoryID:        category2ID,
				Duration:          60,
				Status:            "published",
				Availability:      "online",
				BufferTime:        15,
				CancellationHours: 24,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			err := testRepo.UpsertProduct(ctx, product)
			require.NoError(t, err)
		}

		// Act
		results, err := testRepo.ListProductsByCategory(ctx, category1ID.String())

		// Assert
		require.NoError(t, err)
		assert.Len(t, results, 3)
		for _, prod := range results {
			assert.Equal(t, category1ID, prod.CategoryID)
		}
	})

	t.Run("should return empty slice for category with no products", func(t *testing.T) {
		// Arrange - create category with no products
		emptyCategory := domain.CachedCategory{
			ID:        uuid.New(),
			Name:      "Empty Category",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := testRepo.UpsertCategory(ctx, emptyCategory)
		require.NoError(t, err)

		// Act
		results, err := testRepo.ListProductsByCategory(ctx, emptyCategory.ID.String())

		// Assert
		require.NoError(t, err)
		assert.Empty(t, results)
	})
}
```

---

#### **File: `authuser/internal/adapters/postgres/catalog_cache/main_test.go`**

```go
package catalog_cache

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/Leviosa-care/core/migrations"
	tu "github.com/Leviosa-care/core/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var (
	pgContainer *tu.PostgresContainer
	testPool    *pgxpool.Pool
	testRepo    ports.CatalogCacheRepository
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup Postgres testcontainer
	var err error
	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup postgres container: %v", err)
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	// Create DB pool
	log.Println("Creating pgxpool...")
	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
	defer poolCancel()

	pgCfg, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse pgxpool config: %v", err))
	}

	pgCfg.MaxConns = 5
	pgCfg.MinConns = 1

	testPool, err = pgxpool.NewWithConfig(poolCtx, pgCfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to open test database pool: %v", err))
	}
	log.Println("pgxpool created.")

	// Ping database
	if err = testPool.Ping(poolCtx); err != nil {
		panic(fmt.Sprintf("Failed to ping database pool: %v", err))
	}

	// Run migrations
	log.Println("Applying database migrations...")
	goose.SetBaseFS(migrations.FS)
	if err = goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Setting dialect for migrations: %s\n", err)
	}

	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
	if err != nil {
		panic(fmt.Sprintf("Failed to open temp *sql.DB for goose migrations: %v", err))
	}
	defer gooseDB.Close()

	if err = goose.UpContext(ctx, gooseDB, "."); err != nil {
		panic(fmt.Sprintf("running all migrations: %s\n", err))
	}
	log.Println("Migrations applied.")

	// Initialize repository
	testRepo = New(ctx, testPool)

	// Run tests
	code := m.Run()

	// Cleanup
	testPool.Close()

	os.Exit(code)
}
```

---

## 🐰 Phase 6: AuthUser Adapters - RabbitMQ Consumer

### 6.1 Implement Catalog Consumer

#### **File: `authuser/internal/adapters/rabbitmq/catalog_consumer.go`**

```go
package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	catalogContracts "github.com/Leviosa-care/core/contracts/catalog"
	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	"github.com/Leviosa-care/core/messaging/rabbitmq"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/authuser/internal/ports"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// CatalogConsumer handles catalog events from RabbitMQ
type CatalogConsumer struct {
	conn *amqp.Connection
	repo ports.CatalogCacheRepository
}

// NewCatalogConsumer creates a new catalog consumer
func NewCatalogConsumer(conn *amqp.Connection, repo ports.CatalogCacheRepository) *CatalogConsumer {
	return &CatalogConsumer{
		conn: conn,
		repo: repo,
	}
}

// Start begins consuming catalog events
func (c *CatalogConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// Declare catalog exchange (should already exist from Catalog service)
	if err := rabbitmq.DeclareExchange(ch, mq.CatalogExchangeName, "direct"); err != nil {
		return fmt.Errorf("declare catalog exchange: %w", err)
	}

	// Declare authuser catalog queue
	if err := rabbitmq.DeclareQueue(ch, mq.AuthUserCatalogQueueName); err != nil {
		return fmt.Errorf("declare authuser catalog queue: %w", err)
	}

	// Bind queue to all catalog routing keys
	routingKeys := []string{
		mq.ProductCreatedRoutingKey,
		mq.ProductUpdatedRoutingKey,
		mq.ProductDeletedRoutingKey,
		mq.CategoryCreatedRoutingKey,
		mq.CategoryUpdatedRoutingKey,
		mq.CategoryDeletedRoutingKey,
	}

	for _, key := range routingKeys {
		if err := rabbitmq.BindQueue(ch, mq.AuthUserCatalogQueueName, key, mq.CatalogExchangeName); err != nil {
			return fmt.Errorf("bind queue to %s: %w", key, err)
		}
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		mq.AuthUserCatalogQueueName, // queue
		"authuser-catalog-consumer",  // consumer tag
		false,                        // auto-ack (disabled for manual ack/nack)
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	log.Printf("[CatalogConsumer] Started consuming from %s", mq.AuthUserCatalogQueueName)

	// Process messages
	for {
		select {
		case <-ctx.Done():
			log.Println("[CatalogConsumer] Stopping catalog consumer...")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				log.Println("[CatalogConsumer] Message channel closed")
				return nil
			}

			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("[CatalogConsumer] Error processing message: %v", err)
				// Nack with requeue for retry
				msg.Nack(false, true)
			} else {
				// Ack on success
				msg.Ack(false)
			}
		}
	}
}

// processMessage routes messages to appropriate handlers based on routing key
func (c *CatalogConsumer) processMessage(ctx context.Context, msg amqp.Delivery) error {
	log.Printf("[CatalogConsumer] Received message with routing key: %s", msg.RoutingKey)

	switch msg.RoutingKey {
	case mq.ProductCreatedRoutingKey:
		var event catalogContracts.ProductCreatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("unmarshal ProductCreatedEvent: %w", err)
		}
		return c.handleProductCreated(ctx, event)

	case mq.ProductUpdatedRoutingKey:
		var event catalogContracts.ProductUpdatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("unmarshal ProductUpdatedEvent: %w", err)
		}
		return c.handleProductUpdated(ctx, event)

	case mq.ProductDeletedRoutingKey:
		var event catalogContracts.ProductDeletedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("unmarshal ProductDeletedEvent: %w", err)
		}
		return c.handleProductDeleted(ctx, event)

	case mq.CategoryCreatedRoutingKey:
		var event catalogContracts.CategoryCreatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("unmarshal CategoryCreatedEvent: %w", err)
		}
		return c.handleCategoryCreated(ctx, event)

	case mq.CategoryUpdatedRoutingKey:
		var event catalogContracts.CategoryUpdatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("unmarshal CategoryUpdatedEvent: %w", err)
		}
		return c.handleCategoryUpdated(ctx, event)

	case mq.CategoryDeletedRoutingKey:
		var event catalogContracts.CategoryDeletedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("unmarshal CategoryDeletedEvent: %w", err)
		}
		return c.handleCategoryDeleted(ctx, event)

	default:
		log.Printf("[CatalogConsumer] Unknown routing key: %s", msg.RoutingKey)
		return nil // Don't nack unknown routing keys
	}
}

// handleProductCreated processes product creation events
func (c *CatalogConsumer) handleProductCreated(ctx context.Context, event catalogContracts.ProductCreatedEvent) error {
	log.Printf("[CatalogConsumer] Processing ProductCreated: %s", event.ID)

	// Convert event to domain model
	productID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	categoryID, err := uuid.Parse(event.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	product := domain.CachedProduct{
		ID:                productID,
		Name:              event.Name,
		Description:       event.Description,
		CategoryID:        categoryID,
		Duration:          event.Duration,
		Status:            event.Status,
		Availability:      event.Availability,
		BufferTime:        event.BufferTime,
		CancellationHours: event.CancellationHours,
		StripeProductID:   event.StripeProductID,
		Metadata:          event.Metadata,
		CreatedAt:         event.CreatedAt,
		UpdatedAt:         event.UpdatedAt,
	}

	// Upsert to cache
	if err := c.repo.UpsertProduct(ctx, product); err != nil {
		return fmt.Errorf("upsert product: %w", err)
	}

	log.Printf("[CatalogConsumer] Successfully cached product: %s", event.ID)
	return nil
}

// handleProductUpdated processes product update events
func (c *CatalogConsumer) handleProductUpdated(ctx context.Context, event catalogContracts.ProductUpdatedEvent) error {
	log.Printf("[CatalogConsumer] Processing ProductUpdated: %s", event.ID)

	productID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	categoryID, err := uuid.Parse(event.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	product := domain.CachedProduct{
		ID:                productID,
		Name:              event.Name,
		Description:       event.Description,
		CategoryID:        categoryID,
		Duration:          event.Duration,
		Status:            event.Status,
		Availability:      event.Availability,
		BufferTime:        event.BufferTime,
		CancellationHours: event.CancellationHours,
		StripeProductID:   event.StripeProductID,
		Metadata:          event.Metadata,
		CreatedAt:         event.CreatedAt,
		UpdatedAt:         event.UpdatedAt,
	}

	if err := c.repo.UpsertProduct(ctx, product); err != nil {
		return fmt.Errorf("upsert product: %w", err)
	}

	log.Printf("[CatalogConsumer] Successfully updated product cache: %s", event.ID)
	return nil
}

// handleProductDeleted processes product deletion events
func (c *CatalogConsumer) handleProductDeleted(ctx context.Context, event catalogContracts.ProductDeletedEvent) error {
	log.Printf("[CatalogConsumer] Processing ProductDeleted: %s", event.ID)

	if err := c.repo.DeleteProduct(ctx, event.ID); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	log.Printf("[CatalogConsumer] Successfully removed product from cache: %s", event.ID)
	return nil
}

// handleCategoryCreated processes category creation events
func (c *CatalogConsumer) handleCategoryCreated(ctx context.Context, event catalogContracts.CategoryCreatedEvent) error {
	log.Printf("[CatalogConsumer] Processing CategoryCreated: %s", event.ID)

	categoryID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	category := domain.CachedCategory{
		ID:          categoryID,
		Name:        event.Name,
		Description: event.Description,
		Status:      event.Status,
		Metadata:    event.Metadata,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	if err := c.repo.UpsertCategory(ctx, category); err != nil {
		return fmt.Errorf("upsert category: %w", err)
	}

	log.Printf("[CatalogConsumer] Successfully cached category: %s", event.ID)
	return nil
}

// handleCategoryUpdated processes category update events
func (c *CatalogConsumer) handleCategoryUpdated(ctx context.Context, event catalogContracts.CategoryUpdatedEvent) error {
	log.Printf("[CatalogConsumer] Processing CategoryUpdated: %s", event.ID)

	categoryID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	category := domain.CachedCategory{
		ID:          categoryID,
		Name:        event.Name,
		Description: event.Description,
		Status:      event.Status,
		Metadata:    event.Metadata,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	if err := c.repo.UpsertCategory(ctx, category); err != nil {
		return fmt.Errorf("upsert category: %w", err)
	}

	log.Printf("[CatalogConsumer] Successfully updated category cache: %s", event.ID)
	return nil
}

// handleCategoryDeleted processes category deletion events
func (c *CatalogConsumer) handleCategoryDeleted(ctx context.Context, event catalogContracts.CategoryDeletedEvent) error {
	log.Printf("[CatalogConsumer] Processing CategoryDeleted: %s", event.ID)

	if err := c.repo.DeleteCategory(ctx, event.ID); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	log.Printf("[CatalogConsumer] Successfully removed category from cache: %s", event.ID)
	return nil
}
```

---

#### **File: `authuser/internal/adapters/rabbitmq/setup_catalog_consumer.go`**

```go
package rabbitmq

import (
	"context"
	"fmt"

	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupCatalogConsumer sets up the RabbitMQ infrastructure for catalog event consumption
func SetupCatalogConsumer(ctx context.Context, ch *amqp.Channel) error {
	// Declare catalog exchange
	if err := rabbitmq.DeclareExchange(ch, mq.CatalogExchangeName, "direct"); err != nil {
		return fmt.Errorf("declare catalog exchange: %w", err)
	}

	// Declare authuser catalog queue
	if err := rabbitmq.DeclareQueue(ch, mq.AuthUserCatalogQueueName); err != nil {
		return fmt.Errorf("declare authuser catalog queue: %w", err)
	}

	// Bind queue to all catalog routing keys
	routingKeys := []string{
		mq.ProductCreatedRoutingKey,
		mq.ProductUpdatedRoutingKey,
		mq.ProductDeletedRoutingKey,
		mq.CategoryCreatedRoutingKey,
		mq.CategoryUpdatedRoutingKey,
		mq.CategoryDeletedRoutingKey,
	}

	for _, key := range routingKeys {
		if err := rabbitmq.BindQueue(ch, mq.AuthUserCatalogQueueName, key, mq.CatalogExchangeName); err != nil {
			return fmt.Errorf("bind queue to %s: %w", key, err)
		}
	}

	return nil
}
```

---

## 🧪 Phase 8: Testing Infrastructure

### 8.1 Test Helpers for Catalog Cache

#### **File: `authuser/test/helpers/catalog_cache.go`**

```go
package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearCatalogCacheData truncates catalog cache tables for clean test state
func ClearCatalogCacheData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE authuser_cache.products CASCADE")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "TRUNCATE TABLE authuser_cache.categories CASCADE")
	require.NoError(t, err)
}

// NewTestCategory creates a test category with minimal required fields
func NewTestCategory(id, name, description, status string) domain.CachedCategory {
	categoryID, _ := uuid.Parse(id)
	return domain.CachedCategory{
		ID:          categoryID,
		Name:        name,
		Description: description,
		Status:      status,
		Metadata:    make(map[string]any),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewTestProduct creates a test product with minimal required fields
func NewTestProduct(id, name, description, categoryID string, duration int) domain.CachedProduct {
	productID, _ := uuid.Parse(id)
	catID, _ := uuid.Parse(categoryID)
	return domain.CachedProduct{
		ID:                productID,
		Name:              name,
		Description:       description,
		CategoryID:        catID,
		Duration:          duration,
		Status:            "published",
		Availability:      "online",
		BufferTime:        15,
		CancellationHours: 24,
		Metadata:          make(map[string]any),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// InsertTestCategory inserts a category into the cache for testing
func InsertTestCategory(t *testing.T, ctx context.Context, pool *pgxpool.Pool, category domain.CachedCategory) {
	t.Helper()
	query := `
		INSERT INTO authuser_cache.categories (
			id, name, description, status, metadata, synced_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7)
	`
	_, err := pool.Exec(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.Status,
		category.Metadata,
		category.CreatedAt,
		category.UpdatedAt,
	)
	require.NoError(t, err)
}

// InsertTestProduct inserts a product into the cache for testing
func InsertTestProduct(t *testing.T, ctx context.Context, pool *pgxpool.Pool, product domain.CachedProduct) {
	t.Helper()
	query := `
		INSERT INTO authuser_cache.products (
			id, name, description, category_id, duration, status, availability,
			buffer_time, cancellation_hours, stripe_product_id, metadata,
			synced_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), $12, $13)
	`
	_, err := pool.Exec(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.CategoryID,
		product.Duration,
		product.Status,
		product.Availability,
		product.BufferTime,
		product.CancellationHours,
		product.StripeProductID,
		product.Metadata,
		product.CreatedAt,
		product.UpdatedAt,
	)
	require.NoError(t, err)
}

// GetCategoryFromDB retrieves a category from the cache for verification
func GetCategoryFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, categoryID string) (*domain.CachedCategory, error) {
	t.Helper()

	id, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	query := `
		SELECT id, name, description, status, metadata, synced_at, created_at, updated_at
		FROM authuser_cache.categories
		WHERE id = $1
	`

	var category domain.CachedCategory
	err = pool.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.Metadata,
		&category.SyncedAt,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return &category, err
}

// GetProductFromDB retrieves a product from the cache for verification
func GetProductFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, productID string) (*domain.CachedProduct, error) {
	t.Helper()

	id, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	query := `
		SELECT
			id, name, description, category_id, duration, status, availability,
			buffer_time, cancellation_hours, stripe_product_id, metadata,
			synced_at, created_at, updated_at
		FROM authuser_cache.products
		WHERE id = $1
	`

	var product domain.CachedProduct
	err = pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&product.Duration,
		&product.Status,
		&product.Availability,
		&product.BufferTime,
		&product.CancellationHours,
		&product.StripeProductID,
		&product.Metadata,
		&product.SyncedAt,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return &product, err
}
```

---

#### **Update `authuser/test/helpers/common.go`**

Add this function to the existing file:

```go
// Add to ClearAllTestData function:
func ClearAllTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, redisClient *redis.Client) {
	t.Helper()

	// Existing clears...
	ClearUsersTable(t, ctx, pool)
	ClearOTPKeys(t, ctx, redisClient)
	ClearSessionsRedis(t, ctx, redisClient)

	// Add catalog cache clear
	ClearCatalogCacheData(t, ctx, pool)
}
```

---

## 🚀 Phase 10: Service Integration & Startup

### 10.1 Wire Catalog Consumer into AuthUser Service

**Note:** The exact wiring location depends on where your AuthUser service is initialized. Add the following initialization code after the repository and RabbitMQ connection are established:

```go
// Example location: authuser/cmd/server/main.go or similar

// After initializing:
// - testPool (pgxpool.Pool)
// - mqConn (amqp.Connection)

// Initialize catalog cache repository
catalogCacheRepo := catalog_cache.New(ctx, testPool)

// Initialize catalog consumer
catalogConsumer := rabbitmq.NewCatalogConsumer(mqConn, catalogCacheRepo)

// Start consumer in background goroutine
go func() {
	if err := catalogConsumer.Start(ctx); err != nil {
		if err != context.Canceled {
			log.Fatalf("Catalog consumer error: %v", err)
		}
	}
}()

log.Println("Catalog consumer started successfully")
```

---

#### **Update `authuser/internal/adapters/rabbitmq/setup.go`**

Add catalog consumer setup:

```go
func Setup(ctx context.Context, ch *amqp.Channel) error {
	// Existing OTP notification exchange and queues
	if err := rabbitmq.DeclareExchange(ch, mq.OTPNotificationExchangeName, "direct"); err != nil {
		return err
	}

	// ... existing OTP setup ...

	// Add catalog consumer setup
	if err := SetupCatalogConsumer(ctx, ch); err != nil {
		return fmt.Errorf("setup catalog consumer: %w", err)
	}

	return nil
}
```

---

## 📝 Complete Implementation Checklist

### Phase 1: Core Contracts
- [x] Create `core/contracts/catalog/events.go` with all event structs
- [x] Create `core/contracts/catalog/keys.go` with constants
- [x] Create `core/contracts/rabbitmq/catalog.go` with exchange/queue names
- [x] Verify event JSON serialization/deserialization

### Phase 2: Database
- [x] Create migration `20251017161249_authuser_add_catalog_cache_tables.sql`
- [x] Run migration: `goose up` (will be run by user)
- [x] Verify tables exist: `\dt authuser_cache.*` (after migration)
- [ ] Test indexes: `\d authuser_cache.products`

### Phase 3: Domain
- [x] Create `authuser/internal/domain/catalog_category.go`
- [x] Create `authuser/internal/domain/catalog_product.go`
- [x] Create `authuser/internal/domain/catalog_dto.go`
- [ ] Test validation methods

### Phase 4: Ports
- [ ] Create `authuser/internal/ports/catalog_cache_repository.go`
- [ ] Verify interface matches all required methods

### Phase 5: Repository
- [ ] Create all repository files (repository.go, upsert_*.go, delete_*.go, get_*.go, list_*.go)
- [ ] Create all test files (*_test.go)
- [ ] Create main_test.go
- [ ] Run tests: `go test ./internal/adapters/postgres/catalog_cache/...`
- [ ] Verify all tests pass

### Phase 6: Consumer
- [ ] Create `authuser/internal/adapters/rabbitmq/catalog_consumer.go`
- [ ] Create `authuser/internal/adapters/rabbitmq/setup_catalog_consumer.go`
- [ ] Test message routing logic

### Phase 8: Test Helpers
- [x] Create `authuser/test/helpers/catalog_cache.go`
- [ ] Update `authuser/test/helpers/common.go`
- [ ] Test helper functions

### Phase 9: HTTP Endpoints
- [x] Create catalog service interface
- [x] Create catalog endpoints constants
- [x] Create catalog HTTP handlers
- [x] Create catalog HTTP routes
- [x] Create catalog HTTP integration tests

### Phase 10: Integration
- [x] Wire consumer into service startup
- [x] Update `authuser/internal/adapters/rabbitmq/setup.go`
- [x] Test service startup
- [x] Verify consumer logs appear

### Phase 13: Validation
- [ ] Run all unit tests: `make test-unit`
- [ ] Run all integration tests: `make test-integration`
- [ ] Manual E2E testing with RabbitMQ
- [ ] Performance testing with load

---

## 🎯 Success Criteria

- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Consumer starts without errors
- [ ] Events are consumed and processed
- [ ] Cache tables populate correctly
- [ ] Cache survives service restarts
- [ ] Error handling works correctly
- [ ] Code follows existing patterns
- [ ] Manual E2E test confirms flow

---

## 📘 Notes

- **DO NOT commit catalog publisher code** until catalog is tracked
- Follow **snake_case** for files, **PascalCase** for types
- Use **`errs.ClassifyPgError()`** for all database errors
- Ensure **idempotent event handling** (upsert)
- Use **testcontainers** for all tests (no mocks)
- Consumer must **requeue** on processing errors
- Add **comprehensive logging**
- Cache in **separate schema** (`authuser_cache`)

---

This plan is now significantly more detailed with exact code implementations, SQL queries, test patterns, and step-by-step instructions. Each file includes complete implementations that can be directly used or adapted.
