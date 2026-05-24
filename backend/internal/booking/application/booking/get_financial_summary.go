package booking

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

// financialSummaryMaxFetch is a safety cap on the number of rows fetched from
// the DB before in-memory aggregation. slot_start_time is encrypted at rest and
// cannot be filtered in SQL, so we must pull a bounded set.
const financialSummaryMaxFetch = 10000

func (s *BookingService) GetFinancialSummary(ctx context.Context, from, to time.Time) (*domain.FinancialSummaryResponse, error) {
	// Fetch paid and refunded bookings. payment_status is a plain column so it can
	// be filtered in SQL. slot_start_time is encrypted, so the date range is applied
	// in-memory after decryption. We use CreatedAfter as a SQL pre-filter to narrow
	// the dataset: bookings are rarely created before the slot they cover, so this
	// matches the same pattern as GetAnalyticsSummary without missing real data.
	repoFilter := ports.BookingFilter{
		PaymentStatus: []domain.PaymentStatus{
			domain.PaymentStatusPaid,
			domain.PaymentStatusRefunded,
		},
		CreatedAfter:   &from,
		OrderBy:        "created_at",
		OrderDirection: "desc",
		Limit:          financialSummaryMaxFetch,
	}

	allEncx, err := s.bookingRepo.List(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("list bookings for financial summary: %w", err)
	}

	if len(allEncx) == financialSummaryMaxFetch {
		slog.WarnContext(ctx, "financial summary fetch hit safety cap; results may be incomplete",
			"cap", financialSummaryMaxFetch)
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

	// Filter by slot_start_time within [from, to] in-memory (encrypted field).
	var inRange []*domain.Booking
	for _, b := range decrypted {
		if !b.SlotStartTime.Before(from) && b.SlotStartTime.Before(to) {
			inRange = append(inRange, b)
		}
	}

	// Compute aggregates.
	var grossRevenueCents, refundsCents int
	for _, b := range inRange {
		switch b.PaymentStatus {
		case domain.PaymentStatusPaid:
			grossRevenueCents += b.TotalPriceCents
		case domain.PaymentStatusRefunded:
			refundsCents += b.TotalPriceCents
		}
	}

	// Sort transactions by slot_start_time descending.
	sort.Slice(inRange, func(i, j int) bool {
		return inRange[i].SlotStartTime.After(inRange[j].SlotStartTime)
	})

	// Build enriched transaction rows.
	transactions := make([]domain.FinancialTransaction, 0, len(inRange))
	for _, b := range inRange {
		t := domain.FinancialTransaction{
			ID:            b.ID,
			SlotStartTime: b.SlotStartTime,
			AmountCents:   b.TotalPriceCents,
			PaymentStatus: b.PaymentStatus,
			BookingStatus: b.Status,
			ClientName:    "Utilisateur inconnu",
			PartnerName:   "Praticien inconnu",
			ProductName:   "Produit inconnu",
		}

		if s.authUserClient != nil {
			t.ClientName = s.resolveUserName(ctx, b.ClientID, "Utilisateur inconnu")
			t.PartnerName = s.resolveUserName(ctx, b.PartnerID, "Praticien inconnu")
		}
		if s.productService != nil {
			t.ProductName = s.resolveProductName(ctx, b.ProductID)
		}

		transactions = append(transactions, t)
	}

	return &domain.FinancialSummaryResponse{
		Summary: domain.FinancialSummary{
			GrossRevenueCents: grossRevenueCents,
			RefundsCents:      refundsCents,
			NetRevenueCents:   grossRevenueCents - refundsCents,
		},
		Transactions: transactions,
	}, nil
}
