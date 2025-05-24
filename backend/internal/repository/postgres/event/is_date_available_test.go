package eventRepository_test

import (
	"context"
	"testing"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestIsDateAvailable(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	event := factories.NewBasicEvent(nil)
	tests := []struct {
		name        string
		version     int64
		day         int
		month       int
		year        int
		expectedErr error
	}{
		{
			name:        "not in database",
			version:     20240820013106,
			day:         event.Day,
			month:       event.Month,
			year:        event.Year,
			expectedErr: nil,
		},
		{
			name:        "exists in database",
			version:     20240820023513,
			day:         event.Day,
			month:       event.Month,
			year:        event.Year,
			expectedErr: rp.ErrValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
			defer teardown()
			err := repo.IsDateAvailable(ctx, tt.day, tt.month, tt.year)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
