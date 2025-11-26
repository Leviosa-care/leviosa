package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// MetricsRepository defines the interface for metrics data access
type MetricsRepository interface {
	// GetRoomMetrics retrieves utilization metrics for a room within a date range
	GetRoomMetrics(ctx context.Context, roomID uuid.UUID, startDate, endDate time.Time) ([]*domain.RoomMetrics, error)

	// GetPartnerRoomIDs retrieves all room IDs that a partner has access to
	// The userIDHash parameter should be pre-hashed by the application layer
	GetPartnerRoomIDs(ctx context.Context, userIDHash string) ([]uuid.UUID, error)

	// RefreshMaterializedView refreshes the metrics materialized view
	RefreshMaterializedView(ctx context.Context) error
}
