package eventRepository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestInsertEvent(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	event := factories.NewBasicEvent(map[string]any{})
	tests := []struct {
		name        string
		version     int64
		event       *models.Event
		expectedErr error
	}{
		{
			name:        "nominal case",
			version:     20240820013106,
			event:       event,
			expectedErr: nil,
		},
		{
			name:        "event already exists",
			version:     20240820023513,
			event:       event,
			expectedErr: rp.ErrDatabase,
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
			err = repo.InsertEvent(ctx, tx, tt.event)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
