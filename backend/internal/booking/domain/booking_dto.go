package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingResponse struct {
	ID                 uuid.UUID             `json:"id"`
	AvailabilityID     uuid.UUID             `json:"availability_id"`
	ClientID           *uuid.UUID            `json:"client_id,omitempty"`
	PartnerID          uuid.UUID             `json:"partner_id"`
	RoomID             uuid.UUID             `json:"room_id"`
	ProductID          uuid.UUID             `json:"product_id"`
	SlotStartTime      time.Time             `json:"slot_start_time"`
	SlotEndTime        time.Time             `json:"slot_end_time"`
	Status             BookingStatus         `json:"status"`
	TotalPriceCents    int                   `json:"total_price_cents"`
	Currency           string                `json:"currency"`
	PaymentStatus      PaymentStatus         `json:"payment_status"`
	PaymentIntentID    *string               `json:"payment_intent_id,omitempty"`
	ClientNotes        string                `json:"client_notes,omitempty"`
	PartnerNotes       string                `json:"partner_notes,omitempty"`
	CancellationReason *string               `json:"cancellation_reason,omitempty"`
	CancelledAt        *time.Time            `json:"cancelled_at,omitempty"`
	CompletedAt        *time.Time            `json:"completed_at,omitempty"`
	Availability       *AvailabilityResponse `json:"availability,omitempty"`

	// Booking token for public guest access
	Token string `json:"token,omitempty"`

	// Guest contact fields (populated for guest bookings)
	GuestFirstName string `json:"guest_first_name,omitempty"`
	GuestLastName  string `json:"guest_last_name,omitempty"`
	GuestEmail     string `json:"guest_email,omitempty"`
	GuestPhone     string `json:"guest_phone,omitempty"`
}

// Booking DTOs
type CreateBookingRequest struct {
	AvailabilityID uuid.UUID `json:"availability_id" validate:"required"`
	ClientID       *uuid.UUID `json:"client_id,omitempty"`
	ProductID      uuid.UUID `json:"product_id" validate:"required"`
	SlotStartTime  time.Time `json:"slot_start_time" validate:"required"`
	ClientNotes    string    `json:"client_notes,omitempty" validate:"max=1000"`

	// Guest fields (required when client_id is omitted)
	GuestFirstName string `json:"guest_first_name,omitempty"`
	GuestLastName  string `json:"guest_last_name,omitempty"`
	GuestEmail     string `json:"guest_email,omitempty"`
	GuestPhone     string `json:"guest_phone,omitempty"`
}

type UpdateBookingNotesRequest struct {
	ClientNotes  string `json:"client_notes,omitempty" validate:"max=1000"`
	PartnerNotes string `json:"partner_notes,omitempty" validate:"max=1000"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason" validate:"required,min=1,max=500"`
}

type Transaction struct {
	ID              uuid.UUID     `json:"id"`
	SlotStartTime   string        `json:"slot_start_time"`
	ProductID       uuid.UUID     `json:"product_id"`
	ProductName     string        `json:"product_name"`
	AmountCents     int           `json:"amount_cents"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
}

type EarningsSummary struct {
	CurrentMonthCents int    `json:"current_month_cents"`
	LastMonthCents    int    `json:"last_month_cents"`
	PendingCents      int    `json:"pending_cents"`
	NextPayoutDate    string `json:"next_payout_date"`
	NextPayoutCents   int    `json:"next_payout_cents"`
	Transactions      []Transaction `json:"transactions"`
}

// PartnerBookingResponse is the enriched booking DTO returned by the partner bookings endpoint.
// It includes resolved names for display in the partner agenda UI.
type PartnerBookingResponse struct {
	ID              uuid.UUID     `json:"id"`
	ClientID        *uuid.UUID    `json:"client_id,omitempty"`
	ClientName      string        `json:"client_name"`
	ProductName     string        `json:"product_name"`
	RoomName        string        `json:"room_name"`
	SlotStartTime   time.Time     `json:"slot_start_time"`
	SlotEndTime     time.Time     `json:"slot_end_time"`
	Status          BookingStatus `json:"status"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	TotalPriceCents int           `json:"total_price_cents"`
	Currency        string        `json:"currency"`
	ClientNotes     string        `json:"client_notes,omitempty"`
	PartnerNotes    string        `json:"partner_notes,omitempty"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
}

// AdminBookingResponse is the enriched booking DTO returned by the admin list endpoint.
type AdminBookingResponse struct {
	ID              uuid.UUID     `json:"id"`
	ClientName      string        `json:"client_name"`
	PartnerName     string        `json:"partner_name"`
	ProductName     string        `json:"product_name"`
	RoomName        string        `json:"room_name"`
	SlotStartTime   time.Time     `json:"slot_start_time"`
	SlotEndTime     time.Time     `json:"slot_end_time"`
	Status          BookingStatus `json:"status"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	TotalPriceCents int           `json:"total_price_cents"`
	Currency        string        `json:"currency"`
	CreatedAt       time.Time     `json:"created_at"`
}

// AdminBookingsListResponse is the paginated response for the admin bookings list.
type AdminBookingsListResponse struct {
	Bookings []AdminBookingResponse `json:"bookings"`
	Total    int                    `json:"total"`
	Page     int                    `json:"page"`
	Limit    int                    `json:"limit"`
}
