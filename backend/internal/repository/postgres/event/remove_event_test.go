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

func TestRemoveEvent(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	event := factories.NewBasicEvent(nil)

	tests := []struct {
		name        string
		version     int64
		eventID     string
		expectedErr error
	}{
		{
			name:        "no event in database",
			version:     20240820013106,
			eventID:     event.ID,
			expectedErr: rp.ErrNotDeleted,
		},
		{
			name:        "ID do not exist in database",
			version:     20240820023513,
			eventID:     test.GenerateRandomString(16),
			expectedErr: rp.ErrNotDeleted,
		},
		{
			name:        "nominal case",
			version:     20240820023513,
			eventID:     event.ID,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
			defer teardown()
			err := repo.RemoveEvent(ctx, tt.eventID)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
