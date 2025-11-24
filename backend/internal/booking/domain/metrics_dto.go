package domain

import (
	"time"

	"github.com/google/uuid"
)

// GetRoomMetricsRequest represents a request for room utilization metrics
type GetRoomMetricsRequest struct {
	RoomID    uuid.UUID `json:"room_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

// GetRoomMetricsResponse contains utilization metrics for a room over a date range
type GetRoomMetricsResponse struct {
	RoomID       uuid.UUID               `json:"room_id"`
	StartDate    time.Time               `json:"start_date"`
	EndDate      time.Time               `json:"end_date"`
	DailyMetrics []DailyMetricsResponse  `json:"daily_metrics"`
	Summary      MetricsSummary          `json:"summary"`
}

// DailyMetricsResponse represents metrics for a single day
type DailyMetricsResponse struct {
	Date               time.Time `json:"date"`
	TotalMinutesOpen   int       `json:"total_minutes_open"`
	TotalMinutesBooked int       `json:"total_minutes_booked"`
	UtilizationPercent float64   `json:"utilization_percent"`
	FragmentationCount int       `json:"fragmentation_count"`
	IdleMinutes        int       `json:"idle_minutes"`
	AverageGapMinutes  int       `json:"average_gap_minutes"`
	EfficiencyScore    float64   `json:"efficiency_score"`
}

// MetricsSummary provides aggregate statistics across a date range
type MetricsSummary struct {
	AverageUtilization  float64 `json:"average_utilization"`
	TotalFragmentation  int     `json:"total_fragmentation"`
	TotalIdleMinutes    int     `json:"total_idle_minutes"`
	AverageEfficiency   float64 `json:"average_efficiency"`
	DaysAnalyzed        int     `json:"days_analyzed"`
}

// GetPartnerMetricsRequest represents a request for partner-wide metrics
type GetPartnerMetricsRequest struct {
	PartnerID uuid.UUID `json:"partner_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

// GetPartnerMetricsResponse contains metrics for all rooms a partner has access to
type GetPartnerMetricsResponse struct {
	PartnerID   uuid.UUID                `json:"partner_id"`
	StartDate   time.Time                `json:"start_date"`
	EndDate     time.Time                `json:"end_date"`
	RoomMetrics []GetRoomMetricsResponse `json:"room_metrics"`
	Summary     MetricsSummary           `json:"summary"`
}
