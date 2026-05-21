package domain

type DashboardStats struct {
	BookingsThisWeek      int `json:"bookings_this_week"`
	RevenueThisWeek       int `json:"revenue_this_week"`
	UpcomingBookingsCount int `json:"upcoming_bookings_count"`
	PendingBookingsCount  int `json:"pending_bookings_count"`
	ActiveProductsCount   int `json:"active_products_count"`
}
