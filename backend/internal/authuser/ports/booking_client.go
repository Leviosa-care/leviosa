package ports

import "context"

// BookingClient defines the interface for communicating with the booking service.
// Used by authuser to claim guest bookings after account creation.
type BookingClient interface {
	// ClaimBookings links all guest bookings matching the given email to the
	// specified client. Fire-and-forget semantics: errors are logged but do not
	// propagate to the caller.
	ClaimBookings(ctx context.Context, clientID string, email string)
}
