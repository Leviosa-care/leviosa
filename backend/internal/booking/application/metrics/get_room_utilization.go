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

// GetRoomUtilization retrieves utilization metrics for a specific room
func (s *Service) GetRoomUtilization(
	ctx context.Context,
	roomID uuid.UUID,
	startDate, endDate time.Time,
) (*domain.GetRoomMetricsResponse, error) {
	// Validate date range
	if endDate.Before(startDate) {
		return nil, errs.NewInvalidInputErr(errors.New("end date must be after start date"))
	}

	// Get metrics from repository
	metrics, err := s.metricsRepo.GetRoomMetrics(ctx, roomID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get room metrics: %w", err)
	}

	// Calculate summary statistics
	summary := calculateSummary(metrics)

	// Convert to response DTO
	return &domain.GetRoomMetricsResponse{
		RoomID:       roomID,
		StartDate:    startDate,
		EndDate:      endDate,
		DailyMetrics: convertToDaily(metrics),
		Summary:      summary,
	}, nil
}
