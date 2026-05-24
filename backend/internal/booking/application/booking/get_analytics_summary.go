package booking

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

// analyticsMaxFetch is a safety cap on the number of rows fetched from the DB
// before in-memory aggregation. Encrypted fields cannot be filtered/aggregated in SQL.
const analyticsMaxFetch = 10000

func (s *BookingService) GetAnalyticsSummary(ctx context.Context, months int) (*domain.AnalyticsSummaryResponse, error) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	historyStart := currentMonthStart.AddDate(0, -(months-1), 0)

	// Fetch all paid, confirmed/completed bookings within the history window.
	// We filter on created_at to limit the dataset (status/payment_status are plain columns).
	repoFilter := ports.BookingFilter{
		Status:         []domain.BookingStatus{domain.BookingStatusConfirmed, domain.BookingStatusCompleted},
		PaymentStatus:  []domain.PaymentStatus{domain.PaymentStatusPaid},
		CreatedAfter:   &historyStart,
		OrderBy:        "created_at",
		OrderDirection: "asc",
		Limit:          analyticsMaxFetch,
	}

	bookingsEncx, err := s.bookingRepo.List(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("list bookings for analytics: %w", err)
	}

	if len(bookingsEncx) == analyticsMaxFetch {
		slog.WarnContext(ctx, "analytics fetch hit safety cap; results may be incomplete",
			"cap", analyticsMaxFetch)
	}

	// Decrypt all bookings (slot_start_time and product_id are encrypted at rest).
	var decrypted []*domain.Booking
	for _, encx := range bookingsEncx {
		b, err := domain.DecryptBookingEncx(ctx, s.crypto, encx)
		if err != nil {
			slog.WarnContext(ctx, "skipping booking that failed to decrypt",
				"booking_id", encx.ID,
				"err", err)
			continue
		}
		decrypted = append(decrypted, b)
	}

	// --- Current month KPIs ---
	currentMonthBookings := filterBySlotMonth(decrypted, now.Year(), now.Month())

	var totalRevenue int
	for _, b := range currentMonthBookings {
		totalRevenue += b.TotalPriceCents
	}

	avgBookingValue := 0
	if len(currentMonthBookings) > 0 {
		avgBookingValue = totalRevenue / len(currentMonthBookings)
	}

	newClientsCount, err := s.countNewClients(ctx, currentMonthBookings, currentMonthStart)
	if err != nil {
		return nil, fmt.Errorf("count new clients: %w", err)
	}

	currentMonth := domain.AnalyticsCurrentMonth{
		RevenueCents:         totalRevenue,
		BookingsCount:        len(currentMonthBookings),
		NewClientsCount:      newClientsCount,
		AvgBookingValueCents: avgBookingValue,
	}

	// --- Monthly revenue time series ---
	monthlyRevenue := buildMonthlyRevenue(decrypted, months, currentMonthStart)

	// --- Top products ---
	topProducts := buildTopProducts(ctx, decrypted, s.resolveProductName)

	return &domain.AnalyticsSummaryResponse{
		CurrentMonth:   currentMonth,
		MonthlyRevenue: monthlyRevenue,
		TopProducts:    topProducts,
	}, nil
}

// filterBySlotMonth returns bookings whose slot_start_time falls within the given year and month.
func filterBySlotMonth(bookings []*domain.Booking, year int, month time.Month) []*domain.Booking {
	var out []*domain.Booking
	for _, b := range bookings {
		if b.SlotStartTime.Year() == year && b.SlotStartTime.Month() == month {
			out = append(out, b)
		}
	}
	return out
}

// countNewClients counts distinct clients in currentMonthBookings who have no prior paid booking.
// ClientID is a plaintext field in BookingEncx, so the prior-clients lookup skips decryption.
func (s *BookingService) countNewClients(ctx context.Context, currentMonthBookings []*domain.Booking, currentMonthStart time.Time) (int, error) {
	if len(currentMonthBookings) == 0 {
		return 0, nil
	}

	thisMonthClients := make(map[uuid.UUID]struct{}, len(currentMonthBookings))
	guestBookingsCount := 0
	for _, b := range currentMonthBookings {
		if b.ClientID != nil {
			thisMonthClients[*b.ClientID] = struct{}{}
		} else {
			// Each guest booking is treated as a new client: guests have no persistent
			// identity so we cannot deduplicate them across bookings.
			guestBookingsCount++
		}
	}

	// Query for any paid booking created before this month to identify pre-existing clients.
	priorFilter := ports.BookingFilter{
		Status:        []domain.BookingStatus{domain.BookingStatusConfirmed, domain.BookingStatusCompleted},
		PaymentStatus: []domain.PaymentStatus{domain.PaymentStatusPaid},
		CreatedBefore: &currentMonthStart,
		Limit:         analyticsMaxFetch,
	}
	priorEncx, err := s.bookingRepo.List(ctx, priorFilter)
	if err != nil {
		return 0, fmt.Errorf("list prior bookings: %w", err)
	}

	if len(priorEncx) == analyticsMaxFetch {
		slog.WarnContext(ctx, "prior clients fetch hit safety cap; new_clients_count may be understated",
			"cap", analyticsMaxFetch)
	}

	priorClients := make(map[uuid.UUID]struct{}, len(priorEncx))
	for _, encx := range priorEncx {
		if encx.ClientID != nil {
			priorClients[*encx.ClientID] = struct{}{}
		}
	}

	count := guestBookingsCount
	for clientID := range thisMonthClients {
		if _, existed := priorClients[clientID]; !existed {
			count++
		}
	}
	return count, nil
}

// buildMonthlyRevenue aggregates revenue by calendar month for the last N months.
func buildMonthlyRevenue(bookings []*domain.Booking, months int, currentMonthStart time.Time) []domain.AnalyticsMonthlyRevenue {
	type bucket struct {
		revenue int
		count   int
	}

	buckets := make(map[string]*bucket)

	// Initialise buckets so even months with zero bookings appear.
	for i := 0; i < months; i++ {
		t := currentMonthStart.AddDate(0, -i, 0)
		key := t.Format("2006-01")
		buckets[key] = &bucket{}
	}

	for _, b := range bookings {
		key := b.SlotStartTime.Format("2006-01")
		if _, ok := buckets[key]; ok {
			buckets[key].revenue += b.TotalPriceCents
			buckets[key].count++
		}
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]domain.AnalyticsMonthlyRevenue, 0, len(keys))
	for _, k := range keys {
		result = append(result, domain.AnalyticsMonthlyRevenue{
			Month:         k,
			RevenueCents:  buckets[k].revenue,
			BookingsCount: buckets[k].count,
		})
	}
	return result
}

// buildTopProducts returns the top 5 products by bookings_count with resolved names.
func buildTopProducts(ctx context.Context, bookings []*domain.Booking, resolveName func(context.Context, uuid.UUID) string) []domain.AnalyticsTopProduct {
	type agg struct {
		bookings int
		revenue  int
	}

	productAgg := make(map[uuid.UUID]*agg)
	for _, b := range bookings {
		a, ok := productAgg[b.ProductID]
		if !ok {
			a = &agg{}
			productAgg[b.ProductID] = a
		}
		a.bookings++
		a.revenue += b.TotalPriceCents
	}

	type entry struct {
		id uuid.UUID
		*agg
	}
	entries := make([]entry, 0, len(productAgg))
	for id, a := range productAgg {
		entries = append(entries, entry{id: id, agg: a})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].bookings != entries[j].bookings {
			return entries[i].bookings > entries[j].bookings
		}
		return entries[i].revenue > entries[j].revenue
	})
	if len(entries) > 5 {
		entries = entries[:5]
	}

	result := make([]domain.AnalyticsTopProduct, 0, len(entries))
	for _, e := range entries {
		result = append(result, domain.AnalyticsTopProduct{
			ProductID:     e.id,
			Name:          resolveName(ctx, e.id),
			BookingsCount: e.bookings,
			RevenueCents:  e.revenue,
		})
	}
	return result
}
