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

func TestAddProduct(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	product := factories.NewBasicProduct(nil)

	tests := []struct {
		name        string
		product     *productService.Product
		version     int64
		expectedErr error
	}{
		{
			name:        "product already exists",
			product:     product,
			version:     20250217150632,
			expectedErr: rp.ErrDatabase,
		},
		{
			name:        "nominal case",
			product:     product,
			version:     20250217135721,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, productRepository.New)
			defer teardown()
			err := repo.AddProduct(ctx, tt.product)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
