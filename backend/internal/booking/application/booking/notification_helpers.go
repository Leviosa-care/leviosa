package booking

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

// buildNotificationData creates a BookingNotificationData from a booking.
// The notification service implementation is responsible for fetching
// additional details (user emails, room names, etc.) based on the IDs.
func (s *BookingService) buildNotificationData(booking *domain.Booking, productName string) ports.BookingNotificationData {
	return ports.BookingNotificationData{
		// Required identifiers
		BookingID:  booking.ID,
		ClientID:   booking.ClientID,
		PartnerID:  booking.PartnerID,
		RoomID:     booking.RoomID,
		ProductID:  booking.ProductID,

		// Required timing
		SlotStartTime: booking.SlotStartTime,
		SlotEndTime:   booking.SlotEndTime,

		// Required payment details
		TotalPriceCents: booking.TotalPriceCents,
		Currency:        booking.Currency,

		// Optional pre-populated details
		ProductName: productName,

		// Cancellation details (populated by caller if applicable)
		CancellationReason: booking.CancellationReason,
		CancelledAt:        booking.CancelledAt,
	}
}
