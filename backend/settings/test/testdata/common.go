package testdata

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearSettingsTable truncates all settings tables after each test
func ClearSettingsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE settings.plain, settings.encrypted 
		RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

// ClearS3Bucket removes all objects from the test S3 bucket
func ClearS3Bucket(t *testing.T, ctx context.Context) {
	t.Helper()
	// TODO: Implement S3 cleanup when S3 integration is added
	// For now, this is a no-op since S3 functionality is not yet implemented
}

// ClearAllTestData clears all test data from database and S3
func ClearAllTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	ClearSettingsTable(t, ctx, pool)
	ClearS3Bucket(t, ctx)
}