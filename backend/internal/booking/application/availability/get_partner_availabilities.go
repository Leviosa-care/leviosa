package availability

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

func (s *AvailabilityService) GetPartnerAvailabilities(ctx context.Context, partnerID uuid.UUID, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	availabilities, err := s.availabilityRepo.GetByPartnerID(ctx, partnerID, filter)
	if err != nil {
		return nil, fmt.Errorf("get partner availabilities: %w", err)
	}

	return availabilities, nil
}