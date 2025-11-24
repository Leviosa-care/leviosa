package metricsRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// RefreshMaterializedView manually refreshes the metrics materialized view
func (r *Repository) RefreshMaterializedView(ctx context.Context) error {
	query := `SELECT booking.refresh_room_metrics()`

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return errs.ClassifyPgError("refresh materialized view", err)
	}

	return nil
}
