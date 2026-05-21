package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

func (s *BookingService) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	now := time.Now()
	startOfISOWeek := getStartOfISOWeek(now)

	upcomingFilter := ports.BookingFilter{
		Status: []domain.BookingStatus{domain.BookingStatusConfirmed},
	}

	upcomingBookingsEncx, err := s.bookingRepo.GetUpcoming(ctx, upcomingFilter)
	if err != nil {
		return nil, fmt.Errorf("get upcoming bookings: %w", err)
	}

	pendingFilter := ports.BookingFilter{
		PaymentStatus: []domain.PaymentStatus{domain.PaymentStatusPending},
	}

	pendingBookingsEncx, err := s.bookingRepo.List(ctx, pendingFilter)
	if err != nil {
		return nil, fmt.Errorf("get pending bookings: %w", err)
	}

	thisWeekFilter := ports.BookingFilter{
		CreatedAfter: &startOfISOWeek,
	}

	thisWeekBookingsEncx, err := s.bookingRepo.List(ctx, thisWeekFilter)
	if err != nil {
		return nil, fmt.Errorf("get this week bookings: %w", err)
	}

	publishedProducts, err := s.productService.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get published products: %w", err)
	}

	stats := &domain.DashboardStats{
		BookingsThisWeek:      len(thisWeekBookingsEncx),
		RevenueThisWeek:       0,
		UpcomingBookingsCount: len(upcomingBookingsEncx),
		PendingBookingsCount:  len(pendingBookingsEncx),
		ActiveProductsCount:   len(publishedProducts),
	}

	for _, bookingEncx := range thisWeekBookingsEncx {
		booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt booking %s: %w", bookingEncx.ID, err)
		}
		if booking.PaymentStatus == domain.PaymentStatusPaid {
			stats.RevenueThisWeek += booking.TotalPriceCents
		}
	}

	return stats, nil
}

func getStartOfISOWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	var daysSinceMonday int
	if weekday == time.Sunday {
		daysSinceMonday = 6
	} else {
		daysSinceMonday = int(weekday - time.Monday)
	}
	return time.Date(t.Year(), t.Month(), t.Day()-daysSinceMonday, 0, 0, 0, 0, t.Location())
}
