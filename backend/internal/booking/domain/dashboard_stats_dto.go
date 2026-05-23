package domain

import "time"

type DashboardStats struct {
	BookingsThisWeek      int              `json:"bookings_this_week"`
	RevenueThisWeek       int              `json:"revenue_this_week"`
	UpcomingBookingsCount int              `json:"upcoming_bookings_count"`
	PendingBookingsCount  int              `json:"pending_bookings_count"`
	ActiveProductsCount   int              `json:"active_products_count"`
	RecentBookings        []RecentBooking  `json:"recent_bookings"`
	UpcomingBookings      []UpcomingBooking `json:"upcoming_bookings"`
}

type RecentBooking struct {
	ID           string    `json:"id"`
	ClientName   string    `json:"client_name"`
	ProductName  string    `json:"product_name"`
	PartnerName  string    `json:"partner_name"`
	StartTime    time.Time `json:"start_time"`
	Status       string    `json:"status"`
}

type UpcomingBooking struct {
	ID           string    `json:"id"`
	ClientName   string    `json:"client_name"`
	ProductName  string    `json:"product_name"`
	RoomName     string    `json:"room_name"`
	StartTime    time.Time `json:"start_time"`
	DurationMin  int       `json:"duration_min"`
}
