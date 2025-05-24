package eventRepository_test

// import (
// 	"context"
// 	"testing"
//
// 	"github.com/hengadev/leviosa/internal/domain/event"
// 	"github.com/hengadev/leviosa/internal/sqlite"
// 	"github.com/hengadev/leviosa/internal/sqlite/event"
// 	"github.com/hengadev/leviosa/tests/assert"
// )

// NOTE: I should get the past event for those the user actually registered for

// func TestGetEventForUser(t *testing.T) {
// 	t.Setenv("TEST_MIGRATION_PATH", "../migrations/tests")
//
// 	tests := []struct {
// 		event           *event.Event
// 		expectedEventID string
// 		wantErr         bool
// 		version         int64
// 		name            string
// 	}{
// 		{},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			ctx := context.Background()
// 			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, eventRepository.New)
// 			defer teardown()
// 			eventUser, err := repo.AddEvent(ctx, tt.event)
// 			assert.Equal(t, err != nil, tt.wantErr)
// 			_ = eventUser
// 		})
// 	}
// }
