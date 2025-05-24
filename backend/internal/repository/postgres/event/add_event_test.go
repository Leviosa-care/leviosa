package eventRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestAddEvent(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	randID := test.GenerateRandomString(16)
	event := factories.NewBasicEvent(nil)
	tests := []struct {
		name            string
		version         int64
		event           *models.Event
		expectedEventID string
		expectedErr     error
	}{
		{
			name:    "nominal case",
			version: 20250218185533,
			event: factories.NewBasicEvent(map[string]any{
				"Day":     3,
				"Month":   11,
				"Year":    4,
				"ID":      randID,
				"PriceID": test.GenerateRandomString(16),
			}),
			expectedEventID: randID,
			expectedErr:     nil,
		},
		{
			name:            "test unique (day, month, year) constraint",
			version:         20250218185533,
			event:           event,
			expectedEventID: "",
			expectedErr:     rp.ErrDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
			defer teardown()
			eventID, err := repo.AddEvent(ctx, tt.event)
			assert.Equal(t, eventID, tt.expectedEventID)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
