package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/google/uuid"
)

// BuildingService defines the interface for building business logic
type BuildingService interface {
	// CreateBuilding creates a new building with validation
	CreateBuilding(ctx context.Context, name, address, city, postalCode, country string) (*domain.Building, error)

	// GetBuilding retrieves a building by ID
	GetBuilding(ctx context.Context, id uuid.UUID) (*domain.Building, error)

	// UpdateBuilding updates building details with validation
	UpdateBuilding(ctx context.Context, id uuid.UUID, name, address, city, postalCode, country string) (*domain.Building, error)

	// SetBuildingContactInfo sets optional contact information
	SetBuildingContactInfo(ctx context.Context, id uuid.UUID, description, phone, email string) (*domain.Building, error)

	// DeactivateBuilding marks a building as inactive
	DeactivateBuilding(ctx context.Context, id uuid.UUID) error

	// ActivateBuilding marks a building as active
	ActivateBuilding(ctx context.Context, id uuid.UUID) error

	// ListBuildings retrieves buildings with filtering
	ListBuildings(ctx context.Context, filter BuildingFilter) ([]*domain.Building, error)

	// GetBuildingCount returns total count with filtering
	GetBuildingCount(ctx context.Context, filter BuildingFilter) (int, error)
}

// RoomService defines the interface for room business logic
type RoomService interface {
	// CreateRoom creates a new room with validation
	CreateRoom(ctx context.Context, buildingID uuid.UUID, name string, capacity int) (*domain.Room, error)

	// GetRoom retrieves a room by ID
	GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error)

	// UpdateRoom updates room details with validation
	UpdateRoom(ctx context.Context, id uuid.UUID, name, description, roomNumber string, capacity int) (*domain.Room, error)

	// SetRoomEquipment updates the room's equipment list
	SetRoomEquipment(ctx context.Context, id uuid.UUID, equipment []string) (*domain.Room, error)

	// SetRoomRate sets the hourly rate for the room
	SetRoomRate(ctx context.Context, id uuid.UUID, rateCents int) (*domain.Room, error)

	// ClearRoomRate removes the hourly rate
	ClearRoomRate(ctx context.Context, id uuid.UUID) (*domain.Room, error)

	// DeactivateRoom marks a room as inactive
	DeactivateRoom(ctx context.Context, id uuid.UUID) error

	// ActivateRoom marks a room as active
	ActivateRoom(ctx context.Context, id uuid.UUID) error

	// ListRooms retrieves rooms with filtering
	ListRooms(ctx context.Context, filter RoomFilter) ([]*domain.Room, error)

	// GetRoomsByBuilding retrieves all rooms in a building
	GetRoomsByBuilding(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.Room, error)
}

// RoomAllocationService defines the interface for room allocation business logic
type RoomAllocationService interface {
	// CreateSharedAllocation creates a shared room allocation
	CreateSharedAllocation(ctx context.Context, roomID, partnerID uuid.UUID) (*domain.RoomAllocation, error)

	// CreateDedicatedAllocation creates a dedicated room allocation with time bounds
	CreateDedicatedAllocation(ctx context.Context, roomID, partnerID uuid.UUID, startDate, endDate *time.Time) (*domain.RoomAllocation, error)

	// GetAllocation retrieves an allocation by ID
	GetAllocation(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error)

	// UpdateDedicatedPeriod updates the time period for a dedicated allocation
	UpdateDedicatedPeriod(ctx context.Context, id uuid.UUID, startDate, endDate *time.Time) (*domain.RoomAllocation, error)

	// DeactivateAllocation marks an allocation as inactive
	DeactivateAllocation(ctx context.Context, id uuid.UUID) error

	// GetPartnerAllocations retrieves all allocations for a partner
	GetPartnerAllocations(ctx context.Context, partnerID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error)

	// GetRoomAllocations retrieves all allocations for a room
	GetRoomAllocations(ctx context.Context, roomID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error)

	// CheckPartnerRoomAccess verifies if a partner has access to a room at a specific time
	CheckPartnerRoomAccess(ctx context.Context, partnerID, roomID uuid.UUID, at time.Time) (bool, error)
}

// AvailabilityService defines the interface for availability business logic
type AvailabilityService interface {
	// CreateAvailability creates a new availability slot
	CreateAvailability(ctx context.Context, partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int) (*domain.Availability, error)

	// CreateRecurringAvailability creates a recurring availability slot
	CreateRecurringAvailability(ctx context.Context, partnerID, roomID uuid.UUID, startTime, endTime time.Time, maxCapacity int, pattern domain.RecurrencePattern) (*domain.Availability, error)

	// GetAvailability retrieves an availability by ID
	GetAvailability(ctx context.Context, id uuid.UUID) (*domain.Availability, error)

	// UpdateAvailability updates availability details
	UpdateAvailability(ctx context.Context, id uuid.UUID, startTime, endTime time.Time, serviceType string, priceCents *int, notes string) (*domain.Availability, error)

	// CancelAvailability cancels an availability slot
	CancelAvailability(ctx context.Context, id uuid.UUID) error

	// BlockAvailability blocks an availability slot
	BlockAvailability(ctx context.Context, id uuid.UUID) error

	// GetPartnerAvailabilities retrieves availabilities for a partner
	GetPartnerAvailabilities(ctx context.Context, partnerID uuid.UUID, filter AvailabilityFilter) ([]*domain.Availability, error)

	// GetAvailableSlots retrieves available slots for booking
	GetAvailableSlots(ctx context.Context, filter AvailabilityFilter) ([]*domain.Availability, error)

	// CheckAvailabilityConflict checks for scheduling conflicts
	CheckAvailabilityConflict(ctx context.Context, partnerID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error)
}

// BookingService defines the interface for booking business logic
type BookingService interface {
	// CreateBooking creates a new booking
	CreateBooking(ctx context.Context, availabilityID, clientID uuid.UUID, clientNotes string) (*domain.Booking, error)

	// GetBooking retrieves a booking by ID
	GetBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// UpdateBookingNotes updates client or partner notes
	UpdateBookingNotes(ctx context.Context, id uuid.UUID, clientNotes, partnerNotes string) (*domain.Booking, error)

	// CancelBooking cancels a booking with reason
	CancelBooking(ctx context.Context, id uuid.UUID, reason string) (*domain.Booking, error)

	// CompleteBooking marks a booking as completed
	CompleteBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// MarkNoShow marks a booking as no-show
	MarkNoShow(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// ProcessPayment handles payment processing
	ProcessPayment(ctx context.Context, id uuid.UUID, paymentIntentID string) (*domain.Booking, error)

	// RefundBooking processes a refund
	RefundBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error)

	// GetClientBookings retrieves bookings for a client
	GetClientBookings(ctx context.Context, clientID uuid.UUID, filter BookingFilter) ([]*domain.Booking, error)

	// GetPartnerBookings retrieves bookings for a partner
	GetPartnerBookings(ctx context.Context, partnerID uuid.UUID, filter BookingFilter) ([]*domain.Booking, error)

	// GetUpcomingBookings retrieves upcoming confirmed bookings
	GetUpcomingBookings(ctx context.Context, filter BookingFilter) ([]*domain.Booking, error)
}