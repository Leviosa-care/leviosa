package availability

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *AvailabilityService) CheckAvailabilityConflict(ctx context.Context, partnerID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, partnerID, startTime, endTime, excludeID)
	if err != nil {
		return false, fmt.Errorf("check availability conflict: %w", err)
	}

	return hasConflict, nil
}
