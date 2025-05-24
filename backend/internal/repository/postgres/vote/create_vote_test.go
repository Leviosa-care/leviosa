package voteRepository_test

import (
	"context"
	"testing"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres/vote"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestCreateVote(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	userID := factories.NewBasicUser(nil).ID
	tests := []struct {
		name        string
		version     int64
		userID      string
		days        string
		month       int
		year        int
		expectedErr error
	}{
		{
			name:        "vote already exists but for other days",
			version:     20240820225713,
			userID:      userID,
			days:        "5|18|26",
			month:       4,
			year:        2025,
			expectedErr: rp.ErrDatabase,
		},
		{
			name:        "check non empty days constraint for days",
			version:     20240820225713,
			userID:      userID,
			days:        "",
			month:       5,
			year:        2025,
			expectedErr: rp.ErrDatabase,
		},
		{
			name:        "month to small",
			version:     20240820223653,
			userID:      userID,
			days:        "23|12|6",
			month:       0,
			year:        2025,
			expectedErr: rp.ErrDatabase,
		},
		{
			name:        "month too large",
			version:     20240820223653,
			userID:      userID,
			days:        "23|12|6",
			month:       16,
			year:        2025,
			expectedErr: rp.ErrDatabase,
		},
		{
			name:        "year too small",
			version:     20240820223653,
			userID:      userID,
			days:        "23|12|6",
			month:       4,
			year:        1998,
			expectedErr: rp.ErrDatabase,
		},
		{
			name:        "nominal case",
			version:     20240820223653,
			userID:      userID,
			days:        "23|12|6",
			month:       4,
			year:        2025,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, voteRepository.New)
			defer teardown()
			err := repo.CreateVote(ctx, tt.userID, tt.days, tt.month, tt.year)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
