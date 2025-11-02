package catalog

import (
	"context"

	authRabbitMQ "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/rabbitmq"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Service manages the catalog cache and consumer
type Service struct {
	cache    *CatalogCache // Use concrete type since we need both interfaces
	consumer *authRabbitMQ.CatalogConsumer
	conn     *amqp.Connection
}

// New creates a new catalog service with the given cache and RabbitMQ connection
func New(ctx context.Context, cache *CatalogCache, conn *amqp.Connection) (*Service, error) {
	consumer := authRabbitMQ.NewCatalogConsumer(conn, cache)

	return &Service{
		cache:    cache,
		consumer: consumer,
		conn:     conn,
	}, nil
}

// StartConsumer starts the catalog event consumer
func (s *Service) StartConsumer(ctx context.Context) error {
	return s.consumer.Start(ctx)
}

// GetCache returns the catalog cache for read access
func (s *Service) GetCache() ports.CatalogCache {
	return s.cache
}

// GetCacheUpdater returns the catalog cache updater for write access
func (s *Service) GetCacheUpdater() ports.CatalogCacheUpdater {
	return s.cache
}
