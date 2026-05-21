package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingResponse struct {
	ID                 uuid.UUID             `json:"id"`
	AvailabilityID     uuid.UUID             `json:"availability_id"`
	ClientID           uuid.UUID             `json:"client_id"`
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
}

// Booking DTOs
type CreateBookingRequest struct {
	AvailabilityID uuid.UUID `json:"availability_id" validate:"required"`
	ClientID       uuid.UUID `json:"client_id" validate:"required"`
	ProductID      uuid.UUID `json:"product_id" validate:"required"`
	SlotStartTime  time.Time `json:"slot_start_time" validate:"required"`
	ClientNotes    string    `json:"client_notes,omitempty" validate:"max=1000"`
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
