package metricsRepository

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetRoomMetrics retrieves metrics from the materialized view for a specific room
func (r *Repository) GetRoomMetrics(ctx context.Context, roomID uuid.UUID, startDate, endDate time.Time) ([]*domain.RoomMetrics, error) {
	query := `
		SELECT
			room_id,
			date,
			total_minutes_open,
			total_minutes_booked,
			utilization_percent,
			fragmentation_count,
			idle_minutes,
			average_gap_minutes,
			created_at,
			updated_at
		FROM booking.room_daily_metrics
		WHERE room_id = $1
			AND date >= $2
			AND date <= $3
		ORDER BY date ASC
	`

	rows, err := r.pool.Query(ctx, query, roomID, startDate, endDate)
	if err != nil {
		return nil, errs.ClassifyPgError("query room metrics", err)
	}
	defer rows.Close()

	metrics := []*domain.RoomMetrics{}

	for rows.Next() {
		metric := &domain.RoomMetrics{}

		err := rows.Scan(
			&metric.RoomID,
			&metric.Date,
			&metric.TotalMinutesOpen,
			&metric.TotalMinutesBooked,
			&metric.UtilizationPercent,
			&metric.FragmentationCount,
			&metric.IdleMinutes,
			&metric.AverageGapMinutes,
			&metric.CreatedAt,
			&metric.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan room metric", err)
		}

		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room metrics", err)
	}

	return metrics, nil
}
