package metrics

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

// Service implements ports.MetricsService
type Service struct {
	metricsRepo ports.MetricsRepository
}

// New creates a new metrics service
func New(metricsRepo ports.MetricsRepository) ports.MetricsService {
	return &Service{
		metricsRepo: metricsRepo,
	}
}
