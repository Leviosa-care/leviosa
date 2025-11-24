package metricsHelpers

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshMetricsMaterializedView manually refreshes the metrics materialized view
// This is needed in tests to ensure metrics are calculated before querying
func RefreshMetricsMaterializedView(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	query := `REFRESH MATERIALIZED VIEW CONCURRENTLY booking.room_daily_metrics`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		t.Fatalf("Failed to refresh metrics materialized view: %v", err)
	}
}

// ClearMetricsMaterializedView removes all data from the metrics view
// Useful for test isolation
func ClearMetricsMaterializedView(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	// First refresh to ensure we're working with current data
	RefreshMetricsMaterializedView(t, ctx, pool)

	// Note: We can't directly DELETE from materialized view
	// Instead, we should clear the underlying tables (availabilities, rooms)
	// and then refresh the view
	query := `DELETE FROM booking.availabilities`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		t.Fatalf("Failed to clear availabilities for metrics reset: %v", err)
	}

	// Refresh to reflect cleared data
	RefreshMetricsMaterializedView(t, ctx, pool)
}

// GetDailyMetricsCount returns the number of daily metrics records in the view
// Useful for verifying data population
func GetDailyMetricsCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	query := `SELECT COUNT(*) FROM booking.room_daily_metrics`

	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to get metrics count: %v", err)
	}

	return count
}
