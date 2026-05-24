package booking

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
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
		RecentBookings:        []domain.RecentBooking{},
		UpcomingBookings:      []domain.UpcomingBooking{},
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

	// Build recent and upcoming booking lists only when optional deps are available.
	if s.authUserClient != nil && s.roomService != nil {
		recentBookings, err := s.buildRecentBookings(ctx)
		if err != nil {
			return nil, fmt.Errorf("build recent bookings: %w", err)
		}
		stats.RecentBookings = recentBookings

		upcomingDashboard, err := s.buildUpcomingBookings(ctx, now, upcomingBookingsEncx)
		if err != nil {
			return nil, fmt.Errorf("build upcoming bookings: %w", err)
		}
		stats.UpcomingBookings = upcomingDashboard
	}

	return stats, nil
}

// buildRecentBookings returns the last 5 completed/confirmed bookings across all partners.
func (s *BookingService) buildRecentBookings(ctx context.Context) ([]domain.RecentBooking, error) {
	filter := ports.BookingFilter{
		Status:         []domain.BookingStatus{domain.BookingStatusCompleted, domain.BookingStatusConfirmed},
		Limit:          5,
		OrderBy:        "created_at",
		OrderDirection: "desc",
	}

	bookingsEncx, err := s.bookingRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list recent bookings: %w", err)
	}

	result := make([]domain.RecentBooking, 0, len(bookingsEncx))
	for _, encx := range bookingsEncx {
		booking, err := domain.DecryptBookingEncx(ctx, s.crypto, encx)
		if err != nil {
			continue
		}

		clientName := "Utilisateur inconnu"
		if booking.IsGuestBooking() {
			clientName = booking.GuestDisplayName()
		} else if s.authUserClient != nil {
			clientName = s.resolveUserName(ctx, *booking.ClientID, "Utilisateur inconnu")
		}
		partnerName := s.resolveUserName(ctx, booking.PartnerID, "Praticien inconnu")
		productName := s.resolveProductName(ctx, booking.ProductID)

		result = append(result, domain.RecentBooking{
			ID:          booking.ID.String(),
			ClientName:  clientName,
			ProductName: productName,
			PartnerName: partnerName,
			StartTime:   booking.SlotStartTime,
			Status:      string(booking.Status),
		})
	}

	return result, nil
}

// buildUpcomingBookings returns the next 5 confirmed bookings from now.
// allEncx is the already-fetched slice of upcoming bookings (avoids a second DB round-trip).
func (s *BookingService) buildUpcomingBookings(ctx context.Context, now time.Time, allEncx []*domain.BookingEncx) ([]domain.UpcomingBooking, error) {
	// Decrypt and collect only future bookings.
	var future []*domain.Booking
	for _, encx := range allEncx {
		b, err := domain.DecryptBookingEncx(ctx, s.crypto, encx)
		if err != nil {
			continue
		}
		if b.SlotStartTime.After(now) {
			future = append(future, b)
		}
	}

	// Sort by slot start time ascending.
	sort.Slice(future, func(i, j int) bool {
		return future[i].SlotStartTime.Before(future[j].SlotStartTime)
	})

	// Take at most 5.
	if len(future) > 5 {
		future = future[:5]
	}

	result := make([]domain.UpcomingBooking, 0, len(future))
	for _, b := range future {
		clientName := "Utilisateur inconnu"
		if b.IsGuestBooking() {
			clientName = b.GuestDisplayName()
		} else {
			clientName = s.resolveUserName(ctx, *b.ClientID, "Utilisateur inconnu")
		}
		productName := s.resolveProductName(ctx, b.ProductID)
		roomName := s.resolveRoomName(ctx, b.RoomID)
		durationMin := int(b.SlotEndTime.Sub(b.SlotStartTime).Minutes())

		result = append(result, domain.UpcomingBooking{
			ID:          b.ID.String(),
			ClientName:  clientName,
			ProductName: productName,
			RoomName:    roomName,
			StartTime:   b.SlotStartTime,
			DurationMin: durationMin,
		})
	}

	return result, nil
}

// resolveUserName looks up a display name for a user, returning the fallback on error or empty result.
func (s *BookingService) resolveUserName(ctx context.Context, userID uuid.UUID, fallback string) string {
	name, err := s.authUserClient.GetUserName(ctx, userID)
	if err != nil || name == "" {
		return fallback
	}
	return name
}

// resolveProductName looks up a product name by ID, returning a fallback on error.
func (s *BookingService) resolveProductName(ctx context.Context, productID uuid.UUID) string {
	product, err := s.productService.GetProductByID(ctx, productID.String())
	if err != nil {
		return "Produit inconnu"
	}
	return product.Name
}

// resolveRoomName looks up a room name by ID, returning a fallback on error.
func (s *BookingService) resolveRoomName(ctx context.Context, roomID uuid.UUID) string {
	room, err := s.roomService.GetRoom(ctx, roomID)
	if err != nil {
		return "Salle inconnue"
	}
	return room.Name
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
