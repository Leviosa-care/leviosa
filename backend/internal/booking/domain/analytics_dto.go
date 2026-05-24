package domain

import "github.com/google/uuid"

// AnalyticsCurrentMonth holds KPI figures for the current month.
type AnalyticsCurrentMonth struct {
	RevenueCents       int `json:"revenue_cents"`
	BookingsCount      int `json:"bookings_count"`
	NewClientsCount    int `json:"new_clients_count"`
	AvgBookingValueCents int `json:"avg_booking_value_cents"`
}

// AnalyticsMonthlyRevenue is one data point in the monthly revenue time series.
type AnalyticsMonthlyRevenue struct {
	Month         string `json:"month"`
	RevenueCents  int    `json:"revenue_cents"`
	BookingsCount int    `json:"bookings_count"`
}

// AnalyticsTopProduct is one row in the top-products ranking.
type AnalyticsTopProduct struct {
	ProductID     uuid.UUID `json:"product_id"`
	Name          string    `json:"name"`
	BookingsCount int       `json:"bookings_count"`
	RevenueCents  int       `json:"revenue_cents"`
}

// AnalyticsSummaryResponse is the response shape for GET /admin/analytics/summary.
type AnalyticsSummaryResponse struct {
	CurrentMonth   AnalyticsCurrentMonth    `json:"current_month"`
	MonthlyRevenue []AnalyticsMonthlyRevenue `json:"monthly_revenue"`
	TopProducts    []AnalyticsTopProduct     `json:"top_products"`
}
