package eventRepository_test

import (
	"context"
	"database/sql"
	"testing"

	// rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestLoopQuery(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	events := factories.NewBasicEventList()
	var ids []string
	for _, event := range events {
		ids = append(ids, event.ID)
	}
	productsIDs := []string{
		test.GenerateRandomString(16),
		test.GenerateRandomString(16),
		test.GenerateRandomString(16),
	}

	offersIDs := []string{
		test.GenerateRandomString(16),
		test.GenerateRandomString(16),
		test.GenerateRandomString(16),
	}
	tests := []struct {
		name        string
		version     int64
		eventID     string
		columnName  string
		table       string
		values      []string
		expectedErr error
	}{
		{
			name:        "inserting to event_products",
			version:     20250218185527,
			eventID:     events[0].ID,
			columnName:  "product_id",
			table:       "event_products",
			values:      productsIDs,
			expectedErr: nil,
		},
		{
			name:        "inserting to event_offers",
			version:     20250218185533,
			eventID:     events[0].ID,
			columnName:  "offer_id",
			table:       "event_offers",
			values:      offersIDs,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
			defer teardown()
			tx, err := repo.DB.BeginTx(ctx, &sql.TxOptions{})
			if err != nil {
				t.Errorf("starting transaction")
			}
			err = repo.LoopQuery(
				ctx,
				tx,
				tt.columnName,
				tt.table,
				tt.values,
				tt.eventID,
			)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
