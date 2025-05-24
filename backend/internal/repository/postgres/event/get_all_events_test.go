package eventRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	test "github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	assert "github.com/hengadev/test-assert"
)

func TestGetAllEvents(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	factories_events := factories.NewBasicEventList()
	fields := []string{"ID", "Title", "Description", "City", "PostalCode", "Address1", "Address2", "PlaceCount", "FreePlace", "EncryptedBeginAt", "EncryptedEndAt"}

	tests := []struct {
		name           string
		version        int64
		expectedEvents []*models.Event
		expectedErr    error
	}{
		{
			name:           "no event in database",
			version:        20240820013106,
			expectedEvents: nil,
			expectedErr:    nil,
		},
		{
			name:           "nominal case",
			version:        20240820103230,
			expectedEvents: factories_events,
			expectedErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
			defer teardown()
			events, err := repo.GetAllEvents(ctx)
			assert.EqualError(t, err, tt.expectedErr)
			for idx, _ := range events {
				assert.FieldsEqual(t, events[idx], tt.expectedEvents[idx], fields)
			}
		})
	}
}
