package metrics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetPartnerUtilization retrieves aggregated metrics for all rooms a partner has access to
func (s *Service) GetPartnerUtilization(
	ctx context.Context,
	partnerID uuid.UUID,
	startDate, endDate time.Time,
) (*domain.GetPartnerMetricsResponse, error) {
	// Validate date range
	if endDate.Before(startDate) {
		return nil, errs.NewInvalidInputErr(errors.New("end date must be after start date"))
	}

	// Get all rooms partner has access to
	roomIDs, err := s.metricsRepo.GetPartnerRoomIDs(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get partner room IDs: %w", err)
	}

	// Get metrics for each room
	roomMetrics := make([]domain.GetRoomMetricsResponse, 0, len(roomIDs))
	allMetrics := []*domain.RoomMetrics{}

	for _, roomID := range roomIDs {
		metrics, err := s.metricsRepo.GetRoomMetrics(ctx, roomID, startDate, endDate)
		if err != nil {
			// Skip rooms with errors (e.g., no data available)
			continue
		}

		if len(metrics) > 0 {
			roomMetrics = append(roomMetrics, domain.GetRoomMetricsResponse{
				RoomID:       roomID,
				StartDate:    startDate,
				EndDate:      endDate,
				DailyMetrics: convertToDaily(metrics),
				Summary:      calculateSummary(metrics),
			})

			allMetrics = append(allMetrics, metrics...)
		}
	}

	// Calculate overall summary across all rooms
	overallSummary := calculateSummary(allMetrics)

	return &domain.GetPartnerMetricsResponse{
		PartnerID:   partnerID,
		StartDate:   startDate,
		EndDate:     endDate,
		RoomMetrics: roomMetrics,
		Summary:     overallSummary,
	}, nil
}
