package metrics

import (
	"context"
	"fmt"
)

// RefreshMetrics manually triggers a refresh of the metrics materialized view
func (s *Service) RefreshMetrics(ctx context.Context) error {
	if err := s.metricsRepo.RefreshMaterializedView(ctx); err != nil {
		return fmt.Errorf("refresh materialized view: %w", err)
	}

	return nil
}
