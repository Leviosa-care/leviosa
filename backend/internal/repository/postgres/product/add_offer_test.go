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

func TestAddOffer(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	offer := factories.NewBasicOffer(nil)
	tests := []struct {
		name        string
		offer       *productService.Offer
		version     int64
		expectedErr error
	}{
		{
			name:        "offer already exists",
			offer:       offer,
			version:     20250217165531,
			expectedErr: rp.ErrDatabase,
		},
		{
			name:        "nominal case",
			offer:       offer,
			version:     20250217155108,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, productRepository.New)
			defer teardown()
			err := repo.AddOffer(ctx, tt.offer)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
