package domain

import (
	"time"

	"github.com/google/uuid"
)

// FinancialSummary holds the aggregate KPIs for a date range.
type FinancialSummary struct {
	GrossRevenueCents int `json:"gross_revenue_cents"`
	RefundsCents      int `json:"refunds_cents"`
	NetRevenueCents   int `json:"net_revenue_cents"`
}

// FinancialTransaction is a single row in the financial transactions list.
type FinancialTransaction struct {
	ID              uuid.UUID     `json:"id"`
	SlotStartTime   time.Time     `json:"slot_start_time"`
	ClientName      string        `json:"client_name"`
	PartnerName     string        `json:"partner_name"`
	ProductName     string        `json:"product_name"`
	AmountCents     int           `json:"amount_cents"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	BookingStatus   BookingStatus `json:"booking_status"`
}

// FinancialSummaryResponse is the response shape for GET /admin/bookings/financial-summary.
type FinancialSummaryResponse struct {
	Summary      FinancialSummary      `json:"summary"`
	Transactions []FinancialTransaction `json:"transactions"`
}
