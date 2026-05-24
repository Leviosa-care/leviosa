package bookingHandler

import (
	"net/http"
	"strconv"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

// Helper functions
func parseBookingFilter(r *http.Request) ports.BookingFilter {
	filter := ports.BookingFilter{}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status := domain.BookingStatus(statusStr)
		filter.Status = []domain.BookingStatus{status}
	}

	if paymentStatusStr := r.URL.Query().Get("payment_status"); paymentStatusStr != "" {
		paymentStatus := domain.PaymentStatus(paymentStatusStr)
		filter.PaymentStatus = []domain.PaymentStatus{paymentStatus}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	return filter
}

func convertBookingsToResponses(bookings []*domain.Booking) []domain.BookingResponse {
	var responses []domain.BookingResponse
	for _, booking := range bookings {
		var cancellationReason *string
		if booking.CancellationReason != "" {
			cancellationReason = &booking.CancellationReason
		}

		responses = append(responses, domain.BookingResponse{
			ID:                 booking.ID,
			AvailabilityID:     booking.AvailabilityID,
			ClientID:           booking.ClientID,
			PartnerID:          booking.PartnerID,
			RoomID:             booking.RoomID,
			ProductID:          booking.ProductID,
			SlotStartTime:      booking.SlotStartTime,
			SlotEndTime:        booking.SlotEndTime,
			Status:             booking.Status,
			TotalPriceCents:    booking.TotalPriceCents,
			Currency:           booking.Currency,
			PaymentStatus:      booking.PaymentStatus,
			PaymentIntentID:    booking.PaymentIntentID,
			ClientNotes:        booking.ClientNotes,
			PartnerNotes:       booking.PartnerNotes,
			CancellationReason: cancellationReason,
			CancelledAt:        booking.CancelledAt,
			CompletedAt:        booking.CompletedAt,
			GuestFirstName:     booking.GuestFirstName,
			GuestLastName:      booking.GuestLastName,
			GuestEmail:         booking.GuestEmail,
			GuestPhone:         booking.GuestPhone,
		})
	}
	return responses
}
