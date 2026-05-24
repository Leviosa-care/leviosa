package booking

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

const earningsMaxFetch = 10000

func (s *BookingService) GetPartnerEarnings(ctx context.Context, partnerID uuid.UUID) (*domain.EarningsSummary, error) {
	filter := ports.BookingFilter{Limit: earningsMaxFetch}
	bookingsEncx, err := s.bookingRepo.GetByPartnerID(ctx, partnerID, filter)
	if err != nil {
		return nil, fmt.Errorf("get partner bookings: %w", err)
	}

	if len(bookingsEncx) == earningsMaxFetch {
		slog.WarnContext(ctx, "partner earnings fetch hit safety cap; results may be incomplete",
			"cap", earningsMaxFetch, "partner_id", partnerID)
	}

	bookings := make([]*domain.Booking, 0, len(bookingsEncx))
	for _, bookingEncx := range bookingsEncx {
		booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt booking %s: %w", bookingEncx.ID, err)
		}
		bookings = append(bookings, booking)
	}

	// Build product name cache to avoid N+1 calls to the catalog service.
	productNames := make(map[uuid.UUID]string)
	if s.productService != nil {
		for _, booking := range bookings {
			if _, ok := productNames[booking.ProductID]; !ok {
				productNames[booking.ProductID] = s.resolveProductName(ctx, booking.ProductID)
			}
		}
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	lastMonthDate := now.AddDate(0, -1, 0)
	lastYear, lastMonth, _ := lastMonthDate.Date()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	summary := &domain.EarningsSummary{
		CurrentMonthCents: 0,
		LastMonthCents:    0,
		PendingCents:      0,
		NextPayoutCents:   0,
		NextPayoutDate:    nextMonday(now).Format(time.RFC3339),
		Transactions:      []domain.Transaction{},
	}

	for _, booking := range bookings {
		bookingYear, bookingMonth, _ := booking.SlotStartTime.Date()

		if booking.PaymentStatus == domain.PaymentStatusPaid {
			if bookingYear == currentYear && bookingMonth == currentMonth {
				summary.CurrentMonthCents += booking.TotalPriceCents
			}
			if bookingYear == lastYear && bookingMonth == lastMonth {
				summary.LastMonthCents += booking.TotalPriceCents
			}
			if booking.SlotStartTime.After(sevenDaysAgo) || booking.SlotStartTime.Equal(sevenDaysAgo) {
				summary.NextPayoutCents += booking.TotalPriceCents
			}
		}

		if booking.PaymentStatus == domain.PaymentStatusPending {
			summary.PendingCents += booking.TotalPriceCents
		}

		productName, ok := productNames[booking.ProductID]
		if !ok {
			productName = "Produit inconnu"
		}

		summary.Transactions = append(summary.Transactions, domain.Transaction{
			ID:            booking.ID,
			SlotStartTime: booking.SlotStartTime.Format(time.RFC3339),
			ProductID:     booking.ProductID,
			ProductName:   productName,
			AmountCents:   booking.TotalPriceCents,
			PaymentStatus: booking.PaymentStatus,
		})
	}

	return summary, nil
}

func nextMonday(t time.Time) time.Time {
	daysUntilMonday := (int(time.Monday) - int(t.Weekday()) + 7) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	return t.AddDate(0, 0, daysUntilMonday)
}
