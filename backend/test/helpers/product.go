package helpers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearProductsTable cleans the products table for isolation
func ClearProductsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, `TRUNCATE catalog.products CASCADE;`) // Use CASCADE if foreign keys are involved and you want to clear dependents
	require.NoError(t, err)
}

// NewValidProduct is a helper to create a valid product for tests
func NewValidProduct(name string, categoryID uuid.UUID) *domain.Product {
	return &domain.Product{
		ID:                uuid.New(),
		Name:              name,
		Description:       "Test Description",
		CategoryID:        categoryID,
		Duration:          60,
		CreatedAt:         time.Now().Truncate(time.Second),
		Status:            domain.Draft,
		Availability:      domain.Online,
		BufferTime:        15,
		CancellationHours: 24,
		StripeProductID:   "prod_012",
		Metadata:          map[string]any{"color": "red", "size": "M"},
	}
}

// Helper function to insert a product for test setup
func InsertProduct(t *testing.T, ctx context.Context, pool *pgxpool.Pool, product *domain.Product) {
	t.Helper()
	query := `
	INSERT INTO catalog.products (
		id, name, description, category_id, duration,
		created_at, updated_at, status, availability,
		buffer_time, cancellation_hours, stripe_product_id, metadata
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	// Handle NULLable fields for database insertion
	var desc sql.NullString
	if product.Description != "" { // description is `string,omitempty`
		desc = sql.NullString{String: product.Description, Valid: true}
	}

	metadataJSON, err := json.Marshal(product.Metadata) // metadata is `map[string]any,omitempty`
	require.NoError(t, err, "Failed to marshal metadata for product insertion")

	// If metadata is empty map or nil, marshal will return "null". Convert to nil for DB.
	var finalMetadata []byte
	if string(metadataJSON) != "null" {
		finalMetadata = metadataJSON
	}

	_, err = pool.Exec(ctx, query,
		product.ID,
		product.Name,
		desc,
		product.CategoryID,
		product.Duration,
		product.CreatedAt,
		product.UpdatedAt,
		product.Status,       // Will be converted to string by driver
		product.Availability, // Will be converted to string by driver
		product.BufferTime,
		product.CancellationHours,
		product.StripeProductID,
		finalMetadata,
	)
	require.NoError(t, err, fmt.Sprintf("Failed to pre-insert product '%s' for category %s", product.Name, product.CategoryID))
}

// GetProductByID fetches a product from the database by its ID.
// This is a test helper function used to verify the state of the database.
func GetProductByID(t *testing.T, ctx context.Context, productID uuid.UUID, pool *pgxpool.Pool) (*domain.Product, error) {
	t.Helper()

	query := `
		SELECT
			id,
			name,
			description,
			category_id,
			duration,
			created_at,
			updated_at,
			status,
			availability,
			buffer_time,
			cancellation_hours,
			stripe_product_id,
			metadata
		FROM catalog.products
		WHERE id = $1
	`

	var prod domain.Product
	var metadataJSON []byte

	err := pool.QueryRow(ctx, query, productID).Scan(
		&prod.ID,
		&prod.Name,
		&prod.Description,
		&prod.CategoryID,
		&prod.Duration,
		&prod.CreatedAt,
		&prod.UpdatedAt,
		&prod.Status,
		&prod.Availability,
		&prod.BufferTime,
		&prod.CancellationHours,
		&prod.StripeProductID,
		&metadataJSON,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewRepositoryNotFoundErr(err, "product")
		}
		return nil, errs.ClassifyPgError("get products by ID", err)
	}

	if metadataJSON != nil {
		err = json.Unmarshal(metadataJSON, &prod.Metadata)
		require.NoError(t, err, "Failed to unmarshal metadata from database")
	}

	return &prod, nil
}
