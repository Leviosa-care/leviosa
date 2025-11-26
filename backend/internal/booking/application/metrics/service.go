package metrics

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/hengadev/encx"
)

// Service implements ports.MetricsService
type Service struct {
	metricsRepo ports.MetricsRepository
	crypto      encx.CryptoService
}

// New creates a new metrics service
func New(metricsRepo ports.MetricsRepository, crypto encx.CryptoService) ports.MetricsService {
	return &Service{
		metricsRepo: metricsRepo,
		crypto:      crypto,
	}
}
