package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// BookingNotificationService defines the interface for sending booking-related notifications.
// This interface is implemented by the notification module and injected into the booking service.
type BookingNotificationService interface {
	// SendBookingConfirmation sends a confirmation notification to the client after a booking is created.
	SendBookingConfirmation(ctx context.Context, data BookingNotificationData) error

	// SendBookingCancellation sends a cancellation notification to both client and partner.
	SendBookingCancellation(ctx context.Context, data BookingNotificationData) error

	// SendBookingReminder sends a reminder notification to the client before their appointment.
	SendBookingReminder(ctx context.Context, data BookingNotificationData) error

	// SendPaymentConfirmation sends a payment confirmation notification to the client.
	SendPaymentConfirmation(ctx context.Context, data BookingNotificationData) error

	// SendPaymentFailed sends a notification when payment processing fails.
	SendPaymentFailed(ctx context.Context, data BookingNotificationData) error
}

// BookingNotificationData contains information needed to send booking notifications.
// Required fields (IDs) are always populated by the booking service.
// Optional fields may be populated if available, otherwise the notification
// service implementation is responsible for fetching additional details.
type BookingNotificationData struct {
	// Required: Booking identifiers
	BookingID  uuid.UUID
	ClientID   uuid.UUID
	PartnerID  uuid.UUID
	RoomID     uuid.UUID
	ProductID  uuid.UUID

	// Required: Appointment timing
	SlotStartTime time.Time
	SlotEndTime   time.Time

	// Required: Payment details
	TotalPriceCents int
	Currency        string

	// Optional: Pre-populated details (if available)
	// The notification service implementation may fetch these if not provided
	ClientEmail  string
	ClientName   string
	ClientPhone  string
	PartnerEmail string
	PartnerName  string
	ProductName  string
	RoomName     string
	BuildingName string
	Address      string

	// Cancellation details (only for cancellation notifications)
	CancellationReason string
	CancelledAt        *time.Time
}
