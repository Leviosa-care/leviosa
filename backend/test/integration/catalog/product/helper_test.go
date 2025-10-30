package product_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// clearTables is a helper function to truncate relevant tables after each test
// to ensure test isolation.
func clearTables(t *testing.T, ctx context.Context) {
	t.Helper()
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE catalog.prices, catalog.products, catalog.categories RESTART IDENTITY CASCADE;")
	require.NoError(t, err)
}
