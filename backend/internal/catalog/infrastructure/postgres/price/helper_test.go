package priceRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Helper function for clearing tables before each test (optional but recommended) ---
func clearTables(t *testing.T, ctx context.Context, testPool *pgxpool.Pool) {
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)
}
