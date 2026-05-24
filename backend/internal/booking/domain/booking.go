package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// BookingStatus defines the lifecycle status of a booking
type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusNoShow    BookingStatus = "no_show"
)

// PaymentStatus defines the payment state of a booking
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Booking represents a client reservation of partner availability
type Booking struct {
	ID             uuid.UUID  `json:"id"`
	AvailabilityID uuid.UUID  `json:"availability_id"`
	ClientID       *uuid.UUID `json:"client_id,omitempty"`
	PartnerID      uuid.UUID  `json:"partner_id"`
	RoomID         uuid.UUID  `json:"room_id"`

	// Slot information (product and time boundaries)
	ProductID     uuid.UUID `json:"product_id" encx:"encrypt"`               // Link to catalog product
	SlotStartTime time.Time `json:"slot_start_time" encx:"encrypt"`          // Actual slot start (10-min aligned)
	SlotEndTime   time.Time `json:"slot_end_time" encx:"encrypt"`            // Calculated: SlotStartTime + Product.Duration

	// Booking details (encrypted for GDPR compliance)
	ClientNotes  string `json:"client_notes,omitempty" encx:"encrypt"`
	PartnerNotes string `json:"partner_notes,omitempty" encx:"encrypt"`

	// Guest contact fields (encrypted for GDPR compliance, used when ClientID is nil)
	GuestFirstName string `json:"guest_first_name,omitempty" encx:"encrypt"`
	GuestLastName  string `json:"guest_last_name,omitempty" encx:"encrypt"`
	GuestEmail     string `json:"guest_email,omitempty" encx:"encrypt"`
	GuestPhone     string `json:"guest_phone,omitempty" encx:"encrypt"`

	// Pricing information
	TotalPriceCents int    `json:"total_price_cents"`
	Currency        string `json:"currency"`

	// Payment tracking
	PaymentStatus   PaymentStatus `json:"payment_status"`
	PaymentIntentID *string       `json:"payment_intent_id,omitempty"`

	// Booking lifecycle
	Status             BookingStatus `json:"status"`
	CancelledAt        *time.Time    `json:"cancelled_at,omitempty"`
	CancellationReason string        `json:"cancellation_reason,omitempty" encx:"encrypt"`
	CompletedAt        *time.Time    `json:"completed_at,omitempty"`

	// Administrative fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBooking creates a new booking.
// Either clientID must be non-nil, or guest contact fields must be provided.
// For guest bookings: guestFirstName, guestLastName, and at least one of guestEmail or guestPhone are required.
func NewBooking(availabilityID uuid.UUID, clientID *uuid.UUID, partnerID, roomID uuid.UUID, totalPriceCents int, currency string, guestFirstName, guestLastName, guestEmail, guestPhone string) (*Booking, error) {
	if availabilityID == uuid.Nil {
		return nil, ErrInvalidAvailabilityID
	}
	if partnerID == uuid.Nil {
		return nil, ErrInvalidPartnerID
	}
	if roomID == uuid.Nil {
		return nil, ErrInvalidRoomID
	}
	if totalPriceCents < 0 {
		return nil, ErrInvalidBookingPrice
	}

	// Validate identity: either a registered client or guest contact info
	isRegisteredClient := clientID != nil && *clientID != uuid.Nil
	isGuest := guestFirstName != "" && guestLastName != "" && (guestEmail != "" || guestPhone != "")

	if !isRegisteredClient && !isGuest {
		return nil, ErrInvalidClientID
	}
	if isRegisteredClient && isGuest {
		return nil, ErrAmbiguousBookingIdentity
	}

	if currency == "" {
		currency = "EUR" // Default currency
	}

	return &Booking{
		ID:              uuid.New(),
		AvailabilityID:  availabilityID,
		ClientID:        clientID,
		PartnerID:       partnerID,
		RoomID:          roomID,
		TotalPriceCents: totalPriceCents,
		Currency:        currency,
		PaymentStatus:   PaymentStatusPending,
		Status:          BookingStatusConfirmed,
		GuestFirstName:  guestFirstName,
		GuestLastName:   guestLastName,
		GuestEmail:      guestEmail,
		GuestPhone:      guestPhone,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// SetClientNotes sets notes from the client
func (b *Booking) SetClientNotes(notes string) {
	b.ClientNotes = notes
	b.UpdatedAt = time.Now()
}

// SetPartnerNotes sets private notes from the partner
func (b *Booking) SetPartnerNotes(notes string) {
	b.PartnerNotes = notes
	b.UpdatedAt = time.Now()
}

// SetPaymentIntentID sets the Stripe payment intent ID
func (b *Booking) SetPaymentIntentID(paymentIntentID string) {
	b.PaymentIntentID = &paymentIntentID
	b.UpdatedAt = time.Now()
}

// MarkPaymentPaid marks the booking payment as paid
func (b *Booking) MarkPaymentPaid() error {
	if b.PaymentStatus == PaymentStatusRefunded {
		return ErrCannotMarkRefundedAsPaid
	}
	b.PaymentStatus = PaymentStatusPaid
	b.UpdatedAt = time.Now()
	return nil
}

// MarkPaymentFailed marks the booking payment as failed
func (b *Booking) MarkPaymentFailed() {
	b.PaymentStatus = PaymentStatusFailed
	b.UpdatedAt = time.Now()
}

// RefundPayment marks the booking payment as refunded
func (b *Booking) RefundPayment() error {
	if b.PaymentStatus != PaymentStatusPaid {
		return ErrCannotRefundUnpaidBooking
	}
	b.PaymentStatus = PaymentStatusRefunded
	b.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the booking with a reason
func (b *Booking) Cancel(reason string) error {
	if b.Status == BookingStatusCompleted {
		return ErrCannotCancelCompletedBooking
	}
	if b.Status == BookingStatusCancelled {
		return ErrBookingAlreadyCancelled
	}

	now := time.Now()
	b.Status = BookingStatusCancelled
	b.CancelledAt = &now
	b.CancellationReason = reason
	b.UpdatedAt = now
	return nil
}

// Complete marks the booking as completed
func (b *Booking) Complete() error {
	if b.Status == BookingStatusCancelled {
		return ErrCannotCompleteBooking
	}
	if b.Status == BookingStatusCompleted {
		return ErrBookingAlreadyCompleted
	}

	now := time.Now()
	b.Status = BookingStatusCompleted
	b.CompletedAt = &now
	b.UpdatedAt = now
	return nil
}

// MarkNoShow marks the booking as a no-show
func (b *Booking) MarkNoShow() error {
	if b.Status == BookingStatusCancelled {
		return ErrCannotMarkCancelledAsNoShow
	}
	if b.Status == BookingStatusCompleted {
		return ErrCannotMarkCompletedAsNoShow
	}

	b.Status = BookingStatusNoShow
	b.UpdatedAt = time.Now()
	return nil
}

// IsGuestBooking returns true if this booking was made without a registered account.
func (b *Booking) IsGuestBooking() bool {
	return b.ClientID == nil
}

// GuestDisplayName returns the full name of the guest, or empty string if not a guest booking.
func (b *Booking) GuestDisplayName() string {
	parts := make([]string, 0, 2)
	if b.GuestFirstName != "" {
		parts = append(parts, b.GuestFirstName)
	}
	if b.GuestLastName != "" {
		parts = append(parts, b.GuestLastName)
	}
	return strings.Join(parts, " ")
}

// IsActive checks if the booking is in an active state
func (b *Booking) IsActive() bool {
	return b.Status == BookingStatusConfirmed
}

// IsCancellable checks if the booking can be cancelled
func (b *Booking) IsCancellable() bool {
	return b.Status == BookingStatusConfirmed
}

// RequiresPayment checks if the booking still requires payment
func (b *Booking) RequiresPayment() bool {
	return b.PaymentStatus == PaymentStatusPending || b.PaymentStatus == PaymentStatusFailed
}
