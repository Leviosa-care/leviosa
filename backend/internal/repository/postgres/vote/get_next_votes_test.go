package voteRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/vote"
	"github.com/hengadev/leviosa/internal/repository/postgres/vote"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestGetNextVotes(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	tests := []struct {
		name                   string
		version                int64
		month                  int
		year                   int
		expectedAvailableVotes []*vote.AvailableVote
		expectedErr            error
	}{
		{
			name:                   "no available votes in db",
			version:                20240821105317,
			month:                  4,
			year:                   2025,
			expectedAvailableVotes: nil,
			expectedErr:            nil,
			// expectedErr:            rp.ErrDatabase,
		},
		{
			name:                   "month too small",
			version:                20240821110737,
			month:                  8,
			year:                   2025,
			expectedAvailableVotes: nil,
			expectedErr:            nil,
		},
		{
			name:                   "month too large",
			version:                20240821110737,
			month:                  17,
			year:                   2025,
			expectedAvailableVotes: nil,
			expectedErr:            nil,
		},
		{
			name:                   "year too small",
			version:                20240821110737,
			month:                  17,
			year:                   1998,
			expectedAvailableVotes: nil,
			expectedErr:            nil,
		},
		{
			name:                   "nominal case",
			version:                20240821110737,
			month:                  3,
			year:                   2025,
			expectedAvailableVotes: factories.NewAvailableVotesList(),
			expectedErr:            nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, voteRepository.New)
			defer teardown()
			days, err := repo.GetNextVotes(ctx, tt.month, tt.year)
			assert.EqualError(t, err, tt.expectedErr)
			for i, day := range days {
				assert.ReflectEqual(t, day, tt.expectedAvailableVotes[i])
			}
		})
	}
}

func TestParseDays(t *testing.T) {
	tests := []struct {
		days         string
		expectedDays []int
		wantErr      bool
		name         string
	}{
		{days: "", expectedDays: nil, wantErr: true, name: "no days provided"},
		{days: "12|24|28", expectedDays: []int{12, 24, 28}, wantErr: false, name: "nominal case"},
		{days: "1", expectedDays: []int{1}, wantErr: false, name: "just one day"},
		{days: "rerger", expectedDays: nil, wantErr: true, name: "wrong days format"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			days, err := voteRepository.ExportedParseDays(tt.days)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.ReflectEqual(t, days, tt.expectedDays)
		})
	}
}
