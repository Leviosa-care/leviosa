package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func ClearPricesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE catalog.prices RESTART IDENTITY CASCADE;")
	require.NoError(t, err)
}

func NewValidPrice() *domain.Price {
	return &domain.Price{
		ID:            uuid.New(),
		ProductID:     uuid.New(),
		StripePriceID: fmt.Sprintf("price_%s", uuid.New().String()[:12]), // Generate unique ID
		Amount:        100,
		Currency:      "EUR",
		Interval:      "month",
		IsActive:      true, // Default to true for tests
		CreatedAt:     time.Now().Truncate(time.Millisecond),
		UpdatedAt:     time.Now().Truncate(time.Millisecond), // Will be set by DB trigger on update, but set here for consistency
	}
}

// NewValidCreatePriceRequest creates a valid CreatePriceRequest for testing
func NewValidCreatePriceRequest() *domain.CreatePriceRequest {
	return &domain.CreatePriceRequest{
		Amount:   2000, // €20.00 in cents
		Currency: "EUR",
		Interval: "month",
		Nickname: StrPtr("Test Price"),
		Metadata: map[string]string{
			"test_key": "test_value",
		},
	}
}

// NewValidUpdatePriceRequest creates a valid UpdatePriceRequest for testing
func NewValidUpdatePriceRequest() *domain.UpdatePriceRequest {
	active := false
	nickname := "Updated Test Price"
	return &domain.UpdatePriceRequest{
		Active:   &active,
		Nickname: &nickname,
		Metadata: map[string]string{
			"updated_key": "updated_value",
		},
	}
}

// GetAllPricesForProduct retrieves all prices for a given product from the database
func GetAllPricesForProduct(t *testing.T, ctx context.Context, productID uuid.UUID, pool *pgxpool.Pool) []*domain.Price {
	t.Helper()

	query := `
		SELECT 
			id, product_id, stripe_price_id, amount, currency, interval, is_active, created_at, updated_at
		FROM 
			catalog.prices
		WHERE 
			product_id = $1
		ORDER BY created_at
	`
	rows, err := pool.Query(ctx, query, productID)
	require.NoError(t, err)
	defer rows.Close()

	var prices []*domain.Price
	for rows.Next() {
		var price domain.Price
		err := rows.Scan(
			&price.ID,
			&price.ProductID,
			&price.StripePriceID,
			&price.Amount,
			&price.Currency,
			&price.Interval,
			&price.IsActive,
			&price.CreatedAt,
			&price.UpdatedAt,
		)
		require.NoError(t, err)
		prices = append(prices, &price)
	}

	require.NoError(t, rows.Err())
	return prices
}

// Helper to insert a price for a product
func InsertPrice(t *testing.T, ctx context.Context, price *domain.Price, pool *pgxpool.Pool) {
	query := `
	INSERT INTO catalog.prices (
		id, product_id, stripe_price_id, amount, currency, interval, is_active, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)`
	_, err := pool.Exec(ctx, query,
		price.ID,
		price.ProductID,
		price.StripePriceID,
		price.Amount, // Example amount
		price.Currency,
		price.Interval, // Example interval
		price.IsActive,
		price.CreatedAt,
		price.UpdatedAt,
	)
	require.NoError(t, err, fmt.Sprintf("Failed to pre-insert price for product %s", price.ProductID))
}

// getPriceByID is a helper function to retrieve a single price record from the database
// by its ID for use in testing. It returns the price struct and an error if not found.
func GetPriceByID(t *testing.T, ctx context.Context, id uuid.UUID, pool *pgxpool.Pool) *domain.Price {
	t.Helper()

	query := `
		SELECT 
			id, product_id, stripe_price_id, amount, currency, interval, is_active, created_at, updated_at
		FROM 
			catalog.prices
		WHERE 
			id = $1
	`
	var price domain.Price
	err := pool.QueryRow(ctx, query, id).Scan(
		&price.ID,
		&price.ProductID,
		&price.StripePriceID,
		&price.Amount,
		&price.Currency,
		&price.Interval,
		&price.IsActive,
		&price.CreatedAt,
		&price.UpdatedAt,
	)
	require.NoError(t, err, fmt.Sprintf("Failed to get price for product %s", price.ProductID))
	return &price
}

// getPriceStatus is a helper function to query the database directly and retrieve
// the is_active status of a given price ID.
func GetPriceStatus(t *testing.T, ctx context.Context, priceID string, testPool *pgxpool.Pool) bool {
	t.Helper()
	var isActive bool
	err := testPool.QueryRow(ctx, "SELECT is_active FROM catalog.prices WHERE id = $1", priceID).Scan(&isActive)
	require.NoError(t, err, "Failed to get price status for ID %s", priceID)
	return isActive
}
