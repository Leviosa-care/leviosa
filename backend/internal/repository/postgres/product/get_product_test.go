package productRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/product"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/product"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestGetProduct(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	product := factories.NewBasicProduct(nil)
	tests := []struct {
		name            string
		productID       string
		version         int64
		expectedErr     error
		expectedProduct *productService.Product
	}{
		{
			name:            "product not in database",
			productID:       test.GenerateRandomString(16),
			version:         20250217135721,
			expectedErr:     rp.ErrNotFound,
			expectedProduct: nil,
		},
		{
			name:            "nominal case",
			productID:       product.ID,
			version:         20250217150632,
			expectedErr:     nil,
			expectedProduct: product,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, productRepository.New)
			defer teardown()
			product, err := repo.GetProduct(ctx, tt.productID)
			assert.EqualError(t, err, tt.expectedErr)
			assert.ReflectEqual(t, product, tt.expectedProduct)
		})
	}
}
