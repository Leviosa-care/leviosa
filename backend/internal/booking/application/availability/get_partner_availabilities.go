package availability

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *AvailabilityService) GetPartnerAvailabilities(ctx context.Context, partnerID uuid.UUID, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	availabilitiesEncx, err := s.availabilityRepo.GetByPartnerID(ctx, partnerID, filter)
	if err != nil {
		return nil, fmt.Errorf("get partner availabilities: %w", err)
	}

	var availabilities []*domain.Availability
	for _, availavailabilityEncx := range availabilitiesEncx {
		availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availavailabilityEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("availability", err)
		}
		availabilities = append(availabilities, availability)
	}

	return availabilities, nil
}
