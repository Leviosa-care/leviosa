package catalogHelpers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// TestProduct represents minimal product data for testing
type TestProduct struct {
	ID           uuid.UUID
	CategoryID   uuid.UUID
	Name         string
	Description  string
	Duration     int
	BufferTime   int
	IsPublished  bool
	PriceCents   int
	StripePriceID string
}

// InsertTestProduct inserts a product directly into the database for testing
func InsertTestProduct(t *testing.T, ctx context.Context, pool *pgxpool.Pool, product TestProduct) {
	query := `
		INSERT INTO catalog.products (
			id, category_id, name, description,
			duration, buffer_time, is_published,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	now := time.Now()
	_, err := pool.Exec(ctx, query,
		product.ID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Duration,
		product.BufferTime,
		product.IsPublished,
		now,
		now,
	)
	require.NoError(t, err, "Failed to insert test product")

	// Insert a price for the product
	priceQuery := `
		INSERT INTO catalog.prices (
			id, product_id, amount, currency,
			is_active, stripe_price_id, stripe_product_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	priceID := uuid.New()
	_, err = pool.Exec(ctx, priceQuery,
		priceID,
		product.ID,
		product.PriceCents,
		"USD",
		true,
		product.StripePriceID,
		"", // stripe_product_id
		now,
		now,
	)
	require.NoError(t, err, "Failed to insert test price")
}

// CreateDefaultTestProducts creates standard test products for booking tests
func CreateDefaultTestProducts(t *testing.T, ctx context.Context, pool *pgxpool.Pool, categoryID uuid.UUID) []TestProduct {
	products := []TestProduct{
		{
			ID:            uuid.MustParse("10000000-0000-0000-0000-000000000001"),
			CategoryID:    categoryID,
			Name:          "60-Minute Massage",
			Description:   "Standard 60-minute massage session",
			Duration:      60,
			BufferTime:    15,
			IsPublished:   true,
			PriceCents:    8000,
			StripePriceID: "price_test_60min",
		},
		{
			ID:            uuid.MustParse("10000000-0000-0000-0000-000000000002"),
			CategoryID:    categoryID,
			Name:          "90-Minute Massage",
			Description:   "Extended 90-minute massage session",
			Duration:      90,
			BufferTime:    15,
			IsPublished:   true,
			PriceCents:    12000,
			StripePriceID: "price_test_90min",
		},
		{
			ID:            uuid.MustParse("10000000-0000-0000-0000-000000000003"),
			CategoryID:    categoryID,
			Name:          "30-Minute Consultation",
			Description:   "Initial consultation session",
			Duration:      30,
			BufferTime:    10,
			IsPublished:   true,
			PriceCents:    4000,
			StripePriceID: "price_test_30min",
		},
	}

	for _, product := range products {
		InsertTestProduct(t, ctx, pool, product)
	}

	return products
}

// ClearProductsTable removes all products from the database
func ClearProductsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	// Delete in order to respect foreign keys
	_, err := pool.Exec(ctx, "DELETE FROM catalog.prices")
	require.NoError(t, err, "Failed to clear prices table")

	_, err = pool.Exec(ctx, "DELETE FROM catalog.products")
	require.NoError(t, err, "Failed to clear products table")
}

// CreateTestCategory creates a test category for products
func CreateTestCategory(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	categoryID := uuid.New()
	query := `
		INSERT INTO catalog.categories (
			id, name, description, is_published, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	now := time.Now()
	_, err := pool.Exec(ctx, query,
		categoryID,
		"Massage Therapy",
		"Various massage therapy services",
		true,
		now,
		now,
	)
	require.NoError(t, err, "Failed to create test category")

	return categoryID
}
