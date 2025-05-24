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

func TestGetEventByID(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	event := factories.NewBasicEvent(nil)
	fields := []string{"ID", "Title", "Description", "City", "PostalCode", "Address1", "Address2", "PlaceCount", "FreePlace", "EncryptedBeginAt", "EncryptedEndAt"}
	tests := []struct {
		name          string
		version       int64
		eventID       string
		expectedEvent *models.Event
		expectedErr   error
	}{
		{
			name:          "no event in database",
			version:       20240820013106,
			eventID:       test.GenerateRandomString(16),
			expectedEvent: nil,
			expectedErr:   rp.ErrNotFound,
		},
		{
			name:          "nominal case",
			version:       20240820023513,
			eventID:       event.ID,
			expectedEvent: event,
			expectedErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
			defer teardown()
			event, err := repo.GetEventByID(ctx, tt.eventID)
			assert.EqualError(t, err, tt.expectedErr)
			assert.FieldsEqual(t, event, tt.expectedEvent, fields)
		})
	}
}
