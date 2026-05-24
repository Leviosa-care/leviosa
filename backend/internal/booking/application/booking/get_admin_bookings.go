package booking

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

// adminBookingsMaxFetch is a safety cap on the number of rows fetched from the
// DB before in-memory time-range filtering. slot_start_time is encrypted at
// rest and cannot be filtered in SQL, so we must pull a bounded set. Raise
// this constant if the platform's booking volume grows past it.
const adminBookingsMaxFetch = 5000

func (s *BookingService) GetAdminBookings(ctx context.Context, filter ports.AdminBookingsFilter) (*domain.AdminBookingsListResponse, error) {
	// Default pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}

	// Build repository filter from admin filter
	repoFilter := ports.BookingFilter{
		OrderBy:        "created_at",
		OrderDirection: "desc",
		Limit:          adminBookingsMaxFetch,
	}

	if filter.Status != nil {
		repoFilter.Status = []domain.BookingStatus{*filter.Status}
	}
	if filter.PartnerID != nil {
		repoFilter.PartnerID = filter.PartnerID
	}

	// Fetch all matching bookings (we need to decrypt to apply time-range
	// filters since slot_start_time is encrypted at rest).
	allEncx, err := s.bookingRepo.List(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("list bookings for admin: %w", err)
	}

	if len(allEncx) == adminBookingsMaxFetch {
		slog.WarnContext(ctx, "admin bookings fetch hit safety cap; results may be incomplete",
			"cap", adminBookingsMaxFetch)
	}

	// Decrypt all bookings
	var decrypted []*domain.Booking
	for _, encx := range allEncx {
		b, err := domain.DecryptBookingEncx(ctx, s.crypto, encx)
		if err != nil {
			slog.WarnContext(ctx, "skipping booking that failed to decrypt",
				"booking_id", encx.ID,
				"err", err)
			continue
		}
		decrypted = append(decrypted, b)
	}

	// Apply time-range filters in-memory (encrypted fields)
	if filter.From != nil {
		from := *filter.From
		filtered := decrypted[:0]
		for _, b := range decrypted {
			if !b.SlotStartTime.Before(from) {
				filtered = append(filtered, b)
			}
		}
		decrypted = filtered
	}
	if filter.To != nil {
		to := *filter.To
		filtered := decrypted[:0]
		for _, b := range decrypted {
			if !b.SlotEndTime.After(to) {
				filtered = append(filtered, b)
			}
		}
		decrypted = filtered
	}

	// Sort by slot_start_time descending
	sort.Slice(decrypted, func(i, j int) bool {
		return decrypted[i].SlotStartTime.After(decrypted[j].SlotStartTime)
	})

	total := len(decrypted)

	// Apply pagination
	start := (filter.Page - 1) * filter.Limit
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if end > total {
		end = total
	}
	page := decrypted[start:end]

	// Enrich with names
	bookings := make([]domain.AdminBookingResponse, 0, len(page))
	for _, b := range page {
		bookings = append(bookings, s.enrichBooking(ctx, b))
	}

	return &domain.AdminBookingsListResponse{
		Bookings: bookings,
		Total:    total,
		Page:     filter.Page,
		Limit:    filter.Limit,
	}, nil
}

// enrichBooking populates the enriched fields of an AdminBookingResponse.
func (s *BookingService) enrichBooking(ctx context.Context, b *domain.Booking) domain.AdminBookingResponse {
	resp := domain.AdminBookingResponse{
		ID:              b.ID,
		ClientName:      "Utilisateur inconnu",
		PartnerName:     "Praticien inconnu",
		ProductName:     "Produit inconnu",
		RoomName:        "Salle inconnue",
		SlotStartTime:   b.SlotStartTime,
		SlotEndTime:     b.SlotEndTime,
		Status:          b.Status,
		PaymentStatus:   b.PaymentStatus,
		TotalPriceCents: b.TotalPriceCents,
		Currency:        b.Currency,
		CreatedAt:       b.CreatedAt,
	}

	if s.authUserClient != nil {
		if b.IsGuestBooking() {
			resp.ClientName = b.GuestDisplayName()
		} else {
			resp.ClientName = s.resolveUserName(ctx, *b.ClientID, "Utilisateur inconnu")
		}
		resp.PartnerName = s.resolveUserName(ctx, b.PartnerID, "Praticien inconnu")
	} else if b.IsGuestBooking() {
		resp.ClientName = b.GuestDisplayName()
	}

	if s.productService != nil {
		resp.ProductName = s.resolveProductName(ctx, b.ProductID)
	}

	if s.roomService != nil {
		resp.RoomName = s.resolveRoomName(ctx, b.RoomID)
	}

	return resp
}
