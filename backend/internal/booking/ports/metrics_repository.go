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
	GetPartnerRoomIDs(ctx context.Context, partnerID uuid.UUID) ([]uuid.UUID, error)

	// RefreshMaterializedView refreshes the metrics materialized view
	RefreshMaterializedView(ctx context.Context) error
}
