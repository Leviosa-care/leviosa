package availability

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AvailabilityService) GetAvailableSlots(ctx context.Context, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	availabilitiesEncx, err := s.availabilityRepo.GetAvailableSlots(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get available slots: %w", err)
	}

	// Filter out past availabilities and ensure they are truly available
	var availableSlots []*domain.Availability
	now := time.Now()

	for _, availabilityEncx := range availabilitiesEncx {
		availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availabilityEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("availability", err)
		}
		if availability.IsAvailable() && availabilityEncx.StartTime.After(now) {
			availableSlots = append(availableSlots, availability)
		}
	}

	return availableSlots, nil
}
