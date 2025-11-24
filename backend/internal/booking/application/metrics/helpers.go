package metrics

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
)

// calculateSummary computes aggregate statistics across multiple metrics
func calculateSummary(metrics []*domain.RoomMetrics) domain.MetricsSummary {
	if len(metrics) == 0 {
		return domain.MetricsSummary{
			AverageUtilization: 0,
			TotalFragmentation: 0,
			TotalIdleMinutes:   0,
			AverageEfficiency:  0,
			DaysAnalyzed:       0,
		}
	}

	var totalUtilization float64
	var totalFragmentation int
	var totalIdleMinutes int
	var totalEfficiency float64

	for _, metric := range metrics {
		totalUtilization += metric.UtilizationPercent
		totalFragmentation += metric.FragmentationCount
		totalIdleMinutes += metric.IdleMinutes
		totalEfficiency += metric.EfficiencyScore()
	}

	daysAnalyzed := len(metrics)

	return domain.MetricsSummary{
		AverageUtilization: totalUtilization / float64(daysAnalyzed),
		TotalFragmentation: totalFragmentation,
		TotalIdleMinutes:   totalIdleMinutes,
		AverageEfficiency:  totalEfficiency / float64(daysAnalyzed),
		DaysAnalyzed:       daysAnalyzed,
	}
}

// convertToDaily converts domain metrics to DTO response format
func convertToDaily(metrics []*domain.RoomMetrics) []domain.DailyMetricsResponse {
	dailyMetrics := make([]domain.DailyMetricsResponse, 0, len(metrics))

	for _, metric := range metrics {
		dailyMetrics = append(dailyMetrics, domain.DailyMetricsResponse{
			Date:               metric.Date,
			TotalMinutesOpen:   metric.TotalMinutesOpen,
			TotalMinutesBooked: metric.TotalMinutesBooked,
			UtilizationPercent: metric.UtilizationPercent,
			FragmentationCount: metric.FragmentationCount,
			IdleMinutes:        metric.IdleMinutes,
			AverageGapMinutes:  metric.AverageGapMinutes,
			EfficiencyScore:    metric.EfficiencyScore(),
		})
	}

	return dailyMetrics
}
