package availability

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

type AvailabilityService struct {
	availabilityRepo ports.AvailabilityRepository
	allocationRepo   ports.RoomAllocationRepository
	roomRepo         ports.RoomRepository
	authUserClient   ports.AuthUserClient
}

// New creates a new instance of the availability service
func New(
	availabilityRepo ports.AvailabilityRepository,
	allocationRepo ports.RoomAllocationRepository,
	roomRepo ports.RoomRepository,
	authUserClient ports.AuthUserClient,
) ports.AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
		allocationRepo:   allocationRepo,
		roomRepo:         roomRepo,
		authUserClient:   authUserClient,
	}
}

func (s *AvailabilityService) CreateAvailability(ctx context.Context, partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int) (*domain.Availability, error) {
	// Validate partner exists and is verified
	isValidPartner, err := s.authUserClient.ValidatePartnerExists(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("validate partner: %w", err)
	}
	if !isValidPartner {
		return nil, fmt.Errorf("partner not found or not verified")
	}

	// Verify room exists and is active
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("room not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	if !room.IsActive {
		return nil, fmt.Errorf("cannot create availability for inactive room")
	}

	// Check partner has access to the room at the specified time
	hasAccess, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, partnerID, roomID, startTime)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("partner does not have allocation for this room")
		}
		return nil, fmt.Errorf("check partner room access: %w", err)
	}

	if !hasAccess.IsActiveAt(startTime) || !hasAccess.IsActiveAt(endTime) {
		return nil, fmt.Errorf("partner does not have access to room during specified time")
	}

	// Check for scheduling conflicts
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, partnerID, startTime, endTime, nil)
	if err != nil {
		return nil, fmt.Errorf("check availability conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("availability conflicts with existing schedule")
	}

	// Create domain entity with validation
	availability, err := domain.NewAvailability(partnerID, roomID, startTime, endTime, maxCapacity)
	if err != nil {
		return nil, fmt.Errorf("create availability entity: %w", err)
	}

	// Persist to repository
	if err := s.availabilityRepo.Create(ctx, availability); err != nil {
		return nil, fmt.Errorf("create availability: %w", err)
	}

	return availability, nil
}

func (s *AvailabilityService) CreateRecurringAvailability(ctx context.Context, partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int, pattern domain.RecurrencePattern) (*domain.Availability, error) {
	// Validate partner exists and is verified
	isValidPartner, err := s.authUserClient.ValidatePartnerExists(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("validate partner: %w", err)
	}
	if !isValidPartner {
		return nil, fmt.Errorf("partner not found or not verified")
	}

	// Verify room exists and is active
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("room not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("verify room exists: %w", err)
	}

	if !room.IsActive {
		return nil, fmt.Errorf("cannot create availability for inactive room")
	}

	// Check partner has access to the room
	hasAccess, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, partnerID, roomID, startTime)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("partner does not have allocation for this room")
		}
		return nil, fmt.Errorf("check partner room access: %w", err)
	}

	if !hasAccess.IsActiveAt(startTime) {
		return nil, fmt.Errorf("partner does not have access to room during specified time")
	}

	// Check for scheduling conflicts
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, partnerID, startTime, endTime, nil)
	if err != nil {
		return nil, fmt.Errorf("check availability conflict: %w", err)
	}

	if hasConflict {
		return nil, fmt.Errorf("recurring availability conflicts with existing schedule")
	}

	// Create domain entity with validation
	availability, err := domain.NewRecurringAvailability(partnerID, roomID, startTime, endTime, maxCapacity, pattern)
	if err != nil {
		return nil, fmt.Errorf("create recurring availability entity: %w", err)
	}

	// Persist to repository
	if err := s.availabilityRepo.Create(ctx, availability); err != nil {
		return nil, fmt.Errorf("create recurring availability: %w", err)
	}

	return availability, nil
}

func (s *AvailabilityService) GetAvailability(ctx context.Context, id uuid.UUID) (*domain.Availability, error) {
	availability, err := s.availabilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get availability: %w", err)
	}

	return availability, nil
}

func (s *AvailabilityService) UpdateAvailability(ctx context.Context, id uuid.UUID, startTime, endTime time.Time, serviceType string, priceCents *int, notes string) (*domain.Availability, error) {
	// Get existing availability
	availability, err := s.availabilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get availability for update: %w", err)
	}

	// Check if availability can be updated (not booked)
	if availability.Status == domain.AvailabilityStatusBooked {
		return nil, fmt.Errorf("cannot update booked availability")
	}

	// If time is being changed, check for conflicts
	if !startTime.Equal(availability.StartTime) || !endTime.Equal(availability.EndTime) {
		// Check partner still has room access for new time
		hasAccess, err := s.allocationRepo.GetActiveAllocationForPartnerAndRoom(ctx, availability.PartnerID, availability.RoomID, startTime)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				return nil, fmt.Errorf("partner does not have allocation for this room at new time")
			}
			return nil, fmt.Errorf("check partner room access: %w", err)
		}

		if !hasAccess.IsActiveAt(startTime) || !hasAccess.IsActiveAt(endTime) {
			return nil, fmt.Errorf("partner does not have access to room during new time period")
		}

		// Check for scheduling conflicts (excluding this availability)
		hasConflict, err := s.availabilityRepo.CheckConflict(ctx, availability.PartnerID, startTime, endTime, &id)
		if err != nil {
			return nil, fmt.Errorf("check availability conflict: %w", err)
		}

		if hasConflict {
			return nil, fmt.Errorf("updated time conflicts with existing availability")
		}

		// Update time slot
		if err := availability.UpdateTimeSlot(startTime, endTime); err != nil {
			return nil, fmt.Errorf("update time slot: %w", err)
		}
	}

	// Update service details
	availability.SetServiceDetails(serviceType, priceCents, notes)

	// Persist changes
	if err := s.availabilityRepo.Update(ctx, availability); err != nil {
		return nil, fmt.Errorf("update availability: %w", err)
	}

	return availability, nil
}

func (s *AvailabilityService) CancelAvailability(ctx context.Context, id uuid.UUID) error {
	// Get existing availability
	availability, err := s.availabilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get availability for cancellation: %w", err)
	}

	// Cancel
	availability.Cancel()

	// Persist changes
	if err := s.availabilityRepo.Update(ctx, availability); err != nil {
		return fmt.Errorf("cancel availability: %w", err)
	}

	return nil
}

func (s *AvailabilityService) BlockAvailability(ctx context.Context, id uuid.UUID) error {
	// Get existing availability
	availability, err := s.availabilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.ErrRepositoryNotFound
		}
		return fmt.Errorf("get availability for blocking: %w", err)
	}

	// Block
	availability.Block()

	// Persist changes
	if err := s.availabilityRepo.Update(ctx, availability); err != nil {
		return fmt.Errorf("block availability: %w", err)
	}

	return nil
}

func (s *AvailabilityService) GetPartnerAvailabilities(ctx context.Context, partnerID uuid.UUID, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	availabilities, err := s.availabilityRepo.GetByPartnerID(ctx, partnerID, filter)
	if err != nil {
		return nil, fmt.Errorf("get partner availabilities: %w", err)
	}

	return availabilities, nil
}

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

func (s *AvailabilityService) CheckAvailabilityConflict(ctx context.Context, partnerID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, partnerID, startTime, endTime, excludeID)
	if err != nil {
		return false, fmt.Errorf("check availability conflict: %w", err)
	}

	return hasConflict, nil
}