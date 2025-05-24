package voteRepository_test

import (
	"context"
	"testing"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/vote"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	test "github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestFindVotes(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	userID := factories.NewBasicUser(nil).ID
	tests := []struct {
		name         string
		version      int64
		userID       string
		month        int
		year         int
		expectedDays string
		expectedErr  error
	}{
		{
			name:         "no votes in database",
			version:      20240820223653,
			userID:       userID,
			month:        4,
			year:         2025,
			expectedDays: "",
			expectedErr:  rp.ErrNotFound,
		},
		{
			name:         "wrong ID",
			version:      20240820225713,
			userID:       test.GenerateRandomString(16),
			month:        4,
			year:         2025,
			expectedDays: "",
			expectedErr:  rp.ErrNotFound,
		},
		{
			name:         "nominal case",
			version:      20240820225713,
			userID:       userID,
			month:        4,
			year:         2025,
			expectedDays: "23|12|6",
			expectedErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, voteRepository.New)
			defer teardown()
			days, err := repo.FindVotes(ctx, tt.month, tt.year, tt.userID)
			assert.EqualError(t, err, tt.expectedErr)
			assert.Equal(t, days, tt.expectedDays)
		})
	}
}
