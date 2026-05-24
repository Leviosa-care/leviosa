package bookingadapter

import (
	"context"

	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

type checker struct {
	svc bookingPorts.BookingService
}

// New returns a BookingChecker backed by the booking service.
func New(svc bookingPorts.BookingService) *checker {
	return &checker{svc: svc}
}

// HasBookingRelationship checks whether the client has at least one booking with the partner.
func (c *checker) HasBookingRelationship(ctx context.Context, partnerID, clientID uuid.UUID) (bool, error) {
	cid := clientID
	bookings, err := c.svc.GetPartnerBookings(ctx, partnerID, bookingPorts.BookingFilter{
		ClientID: &cid,
		Limit:    1,
	})
	if err != nil {
		return false, err
	}
	return len(bookings) > 0, nil
}
