package availability

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

func (s *AvailabilityService) GetAvailableSlots(ctx context.Context, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	availabilities, err := s.availabilityRepo.GetAvailableSlots(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get available slots: %w", err)
	}

	// Filter out past availabilities and ensure they are truly available
	var availableSlots []*domain.Availability
	now := time.Now()

	for _, availability := range availabilities {
		if availability.IsAvailable() && availability.StartTime.After(now) {
			availableSlots = append(availableSlots, availability)
		}
	}

	return availableSlots, nil
}