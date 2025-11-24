package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
)

// MetricsService defines the interface for metrics business logic
type MetricsService interface {
	// GetRoomUtilization retrieves utilization metrics for a room
	GetRoomUtilization(ctx context.Context, roomID uuid.UUID, startDate, endDate time.Time) (*domain.GetRoomMetricsResponse, error)

	// GetPartnerUtilization retrieves aggregated metrics for all rooms a partner has access to
	GetPartnerUtilization(ctx context.Context, partnerID uuid.UUID, startDate, endDate time.Time) (*domain.GetPartnerMetricsResponse, error)

	// RefreshMetrics manually triggers a refresh of the metrics materialized view
	RefreshMetrics(ctx context.Context) error
}
