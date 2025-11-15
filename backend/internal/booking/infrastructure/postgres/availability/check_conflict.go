package availabilityRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) CheckConflict(ctx context.Context, partnerID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s.availabilities
		WHERE partner_id = $1
		AND start_time < $2
		AND end_time > $3
		AND status IN ('available', 'booked', 'blocked')
	`, r.schema)

	args := []interface{}{partnerID, endTime, startTime}

	if excludeID != nil {
		query += " AND id != $4"
		args = append(args, *excludeID)
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, errs.ClassifyPgError("check availability conflict", err)
	}

	return count > 0, nil
}
