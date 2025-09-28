package domain

import (
	"time"

	"github.com/google/uuid"
)

// Building DTOs
type CreateBuildingRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Address     string `json:"address" validate:"required,min=1,max=500"`
	City        string `json:"city" validate:"required,min=1,max=100"`
	PostalCode  string `json:"postal_code" validate:"required,min=1,max=20"`
	Country     string `json:"country" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Phone       string `json:"phone,omitempty" validate:"max=20"`
	Email       string `json:"email,omitempty" validate:"email,max=255"`
}

type UpdateBuildingRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Address     string `json:"address" validate:"required,min=1,max=500"`
	City        string `json:"city" validate:"required,min=1,max=100"`
	PostalCode  string `json:"postal_code" validate:"required,min=1,max=20"`
	Country     string `json:"country" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"max=1000"`
}

type UpdateBuildingContactRequest struct {
	Phone string `json:"phone,omitempty" validate:"max=20"`
	Email string `json:"email,omitempty" validate:"email,max=255"`
}

type BuildingResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	Country     string    `json:"country"`
	Description string    `json:"description,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Room DTOs
type CreateRoomRequest struct {
	BuildingID  uuid.UUID `json:"building_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description,omitempty" validate:"max=1000"`
	Capacity    int       `json:"capacity" validate:"required,min=1,max=50"`
	Equipment   []string  `json:"equipment,omitempty"`
	PriceCents  *int      `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Currency    string    `json:"currency,omitempty" validate:"omitempty,len=3"`
}

type UpdateRoomRequest struct {
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description,omitempty" validate:"max=1000"`
	Capacity    int       `json:"capacity" validate:"required,min=1,max=50"`
	Equipment   []string  `json:"equipment,omitempty"`
	PriceCents  *int      `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Currency    string    `json:"currency,omitempty" validate:"omitempty,len=3"`
}

type RoomResponse struct {
	ID          uuid.UUID        `json:"id"`
	BuildingID  uuid.UUID        `json:"building_id"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Capacity    int              `json:"capacity"`
	Equipment   []string         `json:"equipment,omitempty"`
	PriceCents  *int             `json:"price_cents,omitempty"`
	Currency    string           `json:"currency,omitempty"`
	IsActive    bool             `json:"is_active"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Building    *BuildingResponse `json:"building,omitempty"`
}

// Room Allocation DTOs
type CreateSharedAllocationRequest struct {
	RoomID    uuid.UUID `json:"room_id" validate:"required"`
	PartnerID uuid.UUID `json:"partner_id" validate:"required"`
}

type CreateDedicatedAllocationRequest struct {
	RoomID    uuid.UUID  `json:"room_id" validate:"required"`
	PartnerID uuid.UUID  `json:"partner_id" validate:"required"`
	StartDate *time.Time `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type UpdateDedicatedAllocationRequest struct {
	StartDate *time.Time `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type RoomAllocationResponse struct {
	ID             uuid.UUID        `json:"id"`
	RoomID         uuid.UUID        `json:"room_id"`
	PartnerID      uuid.UUID        `json:"partner_id"`
	AllocationType AllocationType   `json:"allocation_type"`
	StartDate      *time.Time       `json:"start_date,omitempty"`
	EndDate        *time.Time       `json:"end_date,omitempty"`
	IsActive       bool             `json:"is_active"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Room           *RoomResponse    `json:"room,omitempty"`
}

// Availability DTOs
type CreateAvailabilityRequest struct {
	PartnerID    uuid.UUID  `json:"partner_id" validate:"required"`
	RoomID       uuid.UUID  `json:"room_id" validate:"required"`
	StartTime    time.Time  `json:"start_time" validate:"required"`
	EndTime      time.Time  `json:"end_time" validate:"required"`
	MaxCapacity  int        `json:"max_capacity" validate:"required,min=1,max=50"`
	ServiceType  string     `json:"service_type,omitempty" validate:"max=255"`
	PriceCents   *int       `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Notes        string     `json:"notes,omitempty" validate:"max=1000"`
}

type CreateRecurringAvailabilityRequest struct {
	PartnerID    uuid.UUID         `json:"partner_id" validate:"required"`
	RoomID       uuid.UUID         `json:"room_id" validate:"required"`
	StartTime    time.Time         `json:"start_time" validate:"required"`
	EndTime      time.Time         `json:"end_time" validate:"required"`
	MaxCapacity  int               `json:"max_capacity" validate:"required,min=1,max=50"`
	Pattern      RecurrencePattern `json:"pattern" validate:"required"`
	ServiceType  string            `json:"service_type,omitempty" validate:"max=255"`
	PriceCents   *int              `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Notes        string            `json:"notes,omitempty" validate:"max=1000"`
}

type UpdateAvailabilityRequest struct {
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	ServiceType string    `json:"service_type,omitempty" validate:"max=255"`
	PriceCents  *int      `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Notes       string    `json:"notes,omitempty" validate:"max=1000"`
}

type AvailabilityResponse struct {
	ID               uuid.UUID            `json:"id"`
	PartnerID        uuid.UUID            `json:"partner_id"`
	RoomID           uuid.UUID            `json:"room_id"`
	StartTime        time.Time            `json:"start_time"`
	EndTime          time.Time            `json:"end_time"`
	MaxCapacity      int                  `json:"max_capacity"`
	CurrentBookings  int                  `json:"current_bookings"`
	Status           AvailabilityStatus   `json:"status"`
	ServiceType      string               `json:"service_type,omitempty"`
	PriceCents       *int                 `json:"price_cents,omitempty"`
	Notes            string               `json:"notes,omitempty"`
	RecurrenceRule   *RecurrencePattern   `json:"recurrence_rule,omitempty"`
	ParentID         *uuid.UUID           `json:"parent_id,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	Room             *RoomResponse        `json:"room,omitempty"`
}

// Booking DTOs
type CreateBookingRequest struct {
	AvailabilityID uuid.UUID `json:"availability_id" validate:"required"`
	ClientID       uuid.UUID `json:"client_id" validate:"required"`
	ClientNotes    string    `json:"client_notes,omitempty" validate:"max=1000"`
}

type UpdateBookingNotesRequest struct {
	ClientNotes  string `json:"client_notes,omitempty" validate:"max=1000"`
	PartnerNotes string `json:"partner_notes,omitempty" validate:"max=1000"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason" validate:"required,min=1,max=500"`
}

type ProcessPaymentRequest struct {
	PaymentIntentID string `json:"payment_intent_id" validate:"required"`
}

type BookingResponse struct {
	ID               uuid.UUID             `json:"id"`
	AvailabilityID   uuid.UUID             `json:"availability_id"`
	ClientID         uuid.UUID             `json:"client_id"`
	PartnerID        uuid.UUID             `json:"partner_id"`
	RoomID           uuid.UUID             `json:"room_id"`
	Status           BookingStatus         `json:"status"`
	TotalPriceCents  int                   `json:"total_price_cents"`
	Currency         string                `json:"currency"`
	PaymentStatus    PaymentStatus         `json:"payment_status"`
	PaymentIntentID  *string               `json:"payment_intent_id,omitempty"`
	ClientNotes      string                `json:"client_notes,omitempty"`
	PartnerNotes     string                `json:"partner_notes,omitempty"`
	CancellationReason *string             `json:"cancellation_reason,omitempty"`
	CancelledAt      *time.Time            `json:"cancelled_at,omitempty"`
	CompletedAt      *time.Time            `json:"completed_at,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	Availability     *AvailabilityResponse `json:"availability,omitempty"`
	Room             *RoomResponse         `json:"room,omitempty"`
}