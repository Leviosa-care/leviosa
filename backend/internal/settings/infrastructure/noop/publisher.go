package noop

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/ports"
)

// NoOpPublisher implements the EventPublisher interface but does nothing.
// This is used when event publishing is disabled or in testing scenarios.
type NoOpPublisher struct{}

// NewPublisher creates a new no-op event publisher.
func NewPublisher() ports.EventPublisher {
	return &NoOpPublisher{}
}

// PublishSettingUpdate does nothing and always returns nil.
// This allows the service to work without RabbitMQ dependencies.
func (p *NoOpPublisher) PublishSettingUpdate(ctx context.Context, key string, value any) error {
	// No-op: event publishing is disabled
	return nil
}
