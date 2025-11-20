package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// CheckPartnerRoomAccess checks if a partner has access to a room at a specific time
// func (s *RoomAllocationService) CheckPartnerRoomAccess(ctx context.Context, partnerID, roomID uuid.UUID, at time.Time) (bool, error) {
func (s *RoomAllocationService) CheckPartnerRoomAccess(ctx context.Context, request *domain.CheckPartnerRoomAccessRequest) (bool, error) {
	if err := request.Valid(ctx); err != nil {
		return false, errs.NewInvalidValueErr(err.Error())
	}

	// Get active allocation for partner and room at the specified time
	allocation, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, request.UserID, request.RoomID, request.At)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return false, nil // No allocation means no access
		}
		return false, fmt.Errorf("check partner room access: %w", err)
	}

	// Check if allocation is active at the specified time
	return allocation.IsActiveAt(request.At), nil
}
