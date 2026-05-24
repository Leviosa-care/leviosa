package bookingHelpers

import (
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

// ptrUUID returns a pointer to the given uuid.UUID.
func ptrUUID(id uuid.UUID) *uuid.UUID {
	return &id
}

// NewTestBookingEncx creates a basic encrypted booking with random UUIDs
func NewTestBookingEncx(t *testing.T) *domain.BookingEncx {
	t.Helper()
	now := time.Now()
	paymentIntentID := "pi_test_" + uuid.New().String()[:8]

	return &domain.BookingEncx{
		ID:             uuid.New(),
		AvailabilityID: uuid.New(),
		ClientID:       ptrUUID(uuid.New()),
		PartnerID:      uuid.New(),
		RoomID:         uuid.New(),

		// Encrypted fields (mock encryption)
		ProductIDEncrypted:          []byte("encrypted_product_id_data"),
		SlotStartTimeEncrypted:      []byte("encrypted_slot_start_time"),
		SlotEndTimeEncrypted:        []byte("encrypted_slot_end_time"),
		ClientNotesEncrypted:        []byte("encrypted_client_notes"),
		PartnerNotesEncrypted:       []byte("encrypted_partner_notes"),
		CancellationReasonEncrypted: []byte(""),

		// Pricing
		TotalPriceCents: 5000, // €50.00
		Currency:        "EUR",

		// Payment
		PaymentStatus:   domain.PaymentStatusPending,
		PaymentIntentID: &paymentIntentID,

		// Status
		Status:      domain.BookingStatusConfirmed,
		CancelledAt: nil,
		CompletedAt: nil,

		// Timestamps
		CreatedAt: now,
		UpdatedAt: now,

		// Encryption metadata
		DEKEncrypted: []byte("mock_dek_encrypted_data"),
		KeyVersion:   1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "booking-test-key",
			EncryptionTime:   now.Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}

// NewTestBookingEncxWithIDs creates an encrypted booking with specific IDs
func NewTestBookingEncxWithIDs(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncx(t)
	bookingEncx.AvailabilityID = availabilityID
	bookingEncx.ClientID = &clientID
	bookingEncx.PartnerID = partnerID
	bookingEncx.RoomID = roomID
	return bookingEncx
}

// NewTestBookingEncxWithSlot creates an encrypted booking with product and slot times
func NewTestBookingEncxWithSlot(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID, productID uuid.UUID,
	slotStart, slotEnd time.Time,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithIDs(t, availabilityID, clientID, partnerID, roomID)

	// Update encrypted fields with specific product and slot times
	bookingEncx.ProductIDEncrypted = []byte("encrypted_product_" + productID.String()[:8])
	bookingEncx.SlotStartTimeEncrypted = []byte("encrypted_slot_start_" + slotStart.Format(time.RFC3339))
	bookingEncx.SlotEndTimeEncrypted = []byte("encrypted_slot_end_" + slotEnd.Format(time.RFC3339))

	return bookingEncx
}

// NewTestBookingEncxWithStatus creates an encrypted booking with specific status
func NewTestBookingEncxWithStatus(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
	status domain.BookingStatus,
	paymentStatus domain.PaymentStatus,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithIDs(t, availabilityID, clientID, partnerID, roomID)
	bookingEncx.Status = status
	bookingEncx.PaymentStatus = paymentStatus

	// Set appropriate timestamps based on status
	now := time.Now()
	if status == domain.BookingStatusCancelled {
		bookingEncx.CancelledAt = &now
		bookingEncx.CancellationReasonEncrypted = []byte("encrypted_cancellation_reason")
	} else if status == domain.BookingStatusCompleted {
		bookingEncx.CompletedAt = &now
		bookingEncx.PaymentStatus = domain.PaymentStatusPaid
	}

	return bookingEncx
}

// NewCompletedBookingEncx creates a completed encrypted booking
func NewCompletedBookingEncx(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
) *domain.BookingEncx {
	t.Helper()
	return NewTestBookingEncxWithStatus(
		t,
		availabilityID, clientID, partnerID, roomID,
		domain.BookingStatusCompleted,
		domain.PaymentStatusPaid,
	)
}

// NewCancelledBookingEncx creates a cancelled encrypted booking
func NewCancelledBookingEncx(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
	cancellationReason string,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithStatus(
		t,
		availabilityID, clientID, partnerID, roomID,
		domain.BookingStatusCancelled,
		domain.PaymentStatusRefunded,
	)
	bookingEncx.CancellationReasonEncrypted = []byte("encrypted_" + cancellationReason)
	return bookingEncx
}

// NewTestBookingEncxWithPaymentIntent creates an encrypted booking with payment intent
func NewTestBookingEncxWithPaymentIntent(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
	paymentIntentID string,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithIDs(t, availabilityID, clientID, partnerID, roomID)
	bookingEncx.PaymentIntentID = &paymentIntentID
	return bookingEncx
}

// NewUpcomingBookingEncx creates an upcoming encrypted booking (future slot time)
func NewUpcomingBookingEncx(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
	daysInFuture int,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithIDs(t, availabilityID, clientID, partnerID, roomID)

	// Set future slot times
	futureStart := time.Now().AddDate(0, 0, daysInFuture)
	futureEnd := futureStart.Add(1 * time.Hour)

	bookingEncx.SlotStartTimeEncrypted = []byte("encrypted_slot_start_" + futureStart.Format(time.RFC3339))
	bookingEncx.SlotEndTimeEncrypted = []byte("encrypted_slot_end_" + futureEnd.Format(time.RFC3339))

	// Ensure status is confirmed for upcoming bookings
	bookingEncx.Status = domain.BookingStatusConfirmed
	bookingEncx.PaymentStatus = domain.PaymentStatusPaid

	return bookingEncx
}

// NewPastBookingEncx creates a past encrypted booking
func NewPastBookingEncx(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
	daysInPast int,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithIDs(t, availabilityID, clientID, partnerID, roomID)

	// Set past slot times
	pastStart := time.Now().AddDate(0, 0, -daysInPast)
	pastEnd := pastStart.Add(1 * time.Hour)

	bookingEncx.SlotStartTimeEncrypted = []byte("encrypted_slot_start_" + pastStart.Format(time.RFC3339))
	bookingEncx.SlotEndTimeEncrypted = []byte("encrypted_slot_end_" + pastEnd.Format(time.RFC3339))

	// Past bookings are typically completed
	bookingEncx.Status = domain.BookingStatusCompleted
	bookingEncx.PaymentStatus = domain.PaymentStatusPaid
	completedAt := pastEnd.Add(5 * time.Minute)
	bookingEncx.CompletedAt = &completedAt

	return bookingEncx
}

// NewNoShowBookingEncx creates a no-show encrypted booking
func NewNoShowBookingEncx(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
) *domain.BookingEncx {
	t.Helper()
	return NewTestBookingEncxWithStatus(
		t,
		availabilityID, clientID, partnerID, roomID,
		domain.BookingStatusNoShow,
		domain.PaymentStatusPaid,
	)
}

// NewGuestBookingEncx creates an encrypted booking for a guest (no client_id)
func NewGuestBookingEncx(
	t *testing.T,
	firstName, lastName, email, phone string,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncx(t)
	bookingEncx.ClientID = nil
	bookingEncx.GuestFirstNameEncrypted = []byte("encrypted_" + firstName)
	bookingEncx.GuestLastNameEncrypted = []byte("encrypted_" + lastName)
	bookingEncx.GuestEmailEncrypted = []byte("encrypted_" + email)
	bookingEncx.GuestPhoneEncrypted = []byte("encrypted_" + phone)
	return bookingEncx
}

// NewTestBookingEncxWithPrice creates an encrypted booking with specific price
func NewTestBookingEncxWithPrice(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
	priceCents int,
	currency string,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithIDs(t, availabilityID, clientID, partnerID, roomID)
	bookingEncx.TotalPriceCents = priceCents
	bookingEncx.Currency = currency
	return bookingEncx
}

// NewTestBookingEncxWithNotes creates an encrypted booking with specific notes
func NewTestBookingEncxWithNotes(
	t *testing.T,
	availabilityID, clientID, partnerID, roomID uuid.UUID,
	clientNotes, partnerNotes string,
) *domain.BookingEncx {
	t.Helper()
	bookingEncx := NewTestBookingEncxWithIDs(t, availabilityID, clientID, partnerID, roomID)
	bookingEncx.ClientNotesEncrypted = []byte("encrypted_" + clientNotes)
	bookingEncx.PartnerNotesEncrypted = []byte("encrypted_" + partnerNotes)
	return bookingEncx
}
