package eventRepository_test

import (
	"context"
	"testing"

	// "github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	test "github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	assert "github.com/hengadev/test-assert"
)

func TestModifyEvent(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	event := factories.NewBasicEvent(nil)
	whereMap := map[string]any{"id": event.ID}

	tests := []struct {
		name        string
		version     int64
		event       *models.Event
		expectedErr error
	}{
		{
			name:        "nil event",
			version:     20240820013106,
			event:       nil,
			expectedErr: rp.ErrValidation,
		},
		{
			name:    "event with prohibited fields for modification",
			version: 20240820023513,
			event: factories.NewBasicEvent(map[string]any{
				"ID": test.GenerateRandomString(16),
			}),
			expectedErr: rp.ErrInternal,
		},
		{
			name:        "nominal case with valid updatable event",
			version:     20240820023513,
			event:       factories.NewModifiableEvent(),
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
			defer teardown()
			err := repo.ModifyEvent(ctx, tt.event, whereMap)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
