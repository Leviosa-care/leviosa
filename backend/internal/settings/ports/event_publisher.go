package ports

import "context"

// EventPublisher defines the interface for publishing setting update events.
// Implementations can use RabbitMQ, other message brokers, or no-op for testing.
type EventPublisher interface {
	// PublishSettingUpdate publishes a setting update event with the given key and value.
	// The key identifies which setting changed (e.g., "company_name", "otp_duration").
	// The value is the new setting value (can be string, int, etc.).
	PublishSettingUpdate(ctx context.Context, key string, value any) error
}
