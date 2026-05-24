package booking

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

func (s *BookingService) GetPartnerBookingsEnriched(ctx context.Context, partnerID uuid.UUID, filter ports.BookingFilter) ([]domain.PartnerBookingResponse, error) {
	bookings, err := s.GetPartnerBookings(ctx, partnerID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.PartnerBookingResponse, 0, len(bookings))
	for _, b := range bookings {
		resp := domain.PartnerBookingResponse{
			ID:              b.ID,
			ClientID:        b.ClientID,
			ClientName:      "Utilisateur inconnu",
			ProductName:     "Produit inconnu",
			RoomName:        "Salle inconnue",
			SlotStartTime:   b.SlotStartTime,
			SlotEndTime:     b.SlotEndTime,
			Status:          b.Status,
			PaymentStatus:   b.PaymentStatus,
			TotalPriceCents: b.TotalPriceCents,
			Currency:        b.Currency,
			ClientNotes:     b.ClientNotes,
			PartnerNotes:    b.PartnerNotes,
			CompletedAt:     b.CompletedAt,
		}

		if b.IsGuestBooking() {
			resp.ClientName = b.GuestDisplayName()
		} else if s.authUserClient != nil {
			resp.ClientName = s.resolveUserName(ctx, *b.ClientID, "Utilisateur inconnu")
		}
		if s.productService != nil {
			resp.ProductName = s.resolveProductName(ctx, b.ProductID)
		}
		if s.roomService != nil {
			resp.RoomName = s.resolveRoomName(ctx, b.RoomID)
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *BookingService) GetPartnerBookings(ctx context.Context, partnerID uuid.UUID, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookingsEncx, err := s.bookingRepo.GetByPartnerID(ctx, partnerID, filter)
	if err != nil {
		return nil, fmt.Errorf("get partner bookings: %w", err)
	}

	// Decrypt each booking
	bookings := make([]*domain.Booking, 0, len(bookingsEncx))
	for _, bookingEncx := range bookingsEncx {
		booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt booking %s: %w", bookingEncx.ID, err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
