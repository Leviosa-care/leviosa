package product_test

import (
	"context"
	"testing"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
)

// clearTables is a helper function to truncate relevant tables after each test
// to ensure test isolation.
func clearTables(t *testing.T, ctx context.Context) {
	t.Helper()
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)
	td.ClearImagesTable(t, ctx, testPool)
	tu.ClearAuthData(t, ctx, authCtx)
}
