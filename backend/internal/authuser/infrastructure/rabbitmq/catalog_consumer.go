package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	catalogContracts "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/catalog"
	rabbitmqContracts "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	coreRabbitmq "github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// CatalogConsumer handles catalog events from RabbitMQ
type CatalogConsumer struct {
	conn  *amqp.Connection
	cache ports.CatalogCacheUpdater
}

// NewCatalogConsumer creates a new catalog event consumer
func NewCatalogConsumer(conn *amqp.Connection, cache ports.CatalogCacheUpdater) *CatalogConsumer {
	return &CatalogConsumer{
		conn:  conn,
		cache: cache,
	}
}

// Start begins consuming catalog event messages
func (c *CatalogConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare catalog exchange
	if err := coreRabbitmq.DeclareExchange(ch, rabbitmqContracts.CatalogExchangeName, "topic"); err != nil {
		return err
	}

	// Declare authuser catalog queue
	if err := coreRabbitmq.DeclareQueue(ch, rabbitmqContracts.AuthUserCatalogQueueName); err != nil {
		return err
	}

	// Bind queue to exchange for all catalog events
	routingKeys := []string{
		rabbitmqContracts.CategoryCreatedRoutingKey,
		rabbitmqContracts.CategoryUpdatedRoutingKey,
		rabbitmqContracts.CategoryDeletedRoutingKey,
		rabbitmqContracts.ProductCreatedRoutingKey,
		rabbitmqContracts.ProductUpdatedRoutingKey,
		rabbitmqContracts.ProductDeletedRoutingKey,
	}

	for _, routingKey := range routingKeys {
		if err := coreRabbitmq.BindQueue(ch, rabbitmqContracts.AuthUserCatalogQueueName, routingKey, rabbitmqContracts.CatalogExchangeName); err != nil {
			return err
		}
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		rabbitmqContracts.AuthUserCatalogQueueName, // queue
		"authuser-catalog-consumer",                // consumer
		false,                                      // auto-ack
		false,                                      // exclusive
		false,                                      // no-local
		false,                                      // no-wait
		nil,                                        // args
	)
	if err != nil {
		return err
	}

	log.Printf("Starting authuser catalog consumer...")

	// Process messages
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping authuser catalog consumer...")
			return ctx.Err()
		case msg := <-msgs:
			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("Error processing catalog message: %v", err)
				// Nack the message so it can be retried
				msg.Nack(false, true)
			} else {
				// Ack the message
				msg.Ack(false)
			}
		}
	}
}

// processMessage handles individual catalog event messages
func (c *CatalogConsumer) processMessage(ctx context.Context, msg amqp.Delivery) error {
	log.Printf("Received catalog event: routingKey=%s", msg.RoutingKey)

	switch msg.RoutingKey {
	case rabbitmqContracts.CategoryCreatedRoutingKey:
		return c.processCategoryCreated(ctx, msg.Body)
	case rabbitmqContracts.CategoryUpdatedRoutingKey:
		return c.processCategoryUpdated(ctx, msg.Body)
	case rabbitmqContracts.CategoryDeletedRoutingKey:
		return c.processCategoryDeleted(ctx, msg.Body)
	case rabbitmqContracts.ProductCreatedRoutingKey:
		return c.processProductCreated(ctx, msg.Body)
	case rabbitmqContracts.ProductUpdatedRoutingKey:
		return c.processProductUpdated(ctx, msg.Body)
	case rabbitmqContracts.ProductDeletedRoutingKey:
		return c.processProductDeleted(ctx, msg.Body)
	default:
		log.Printf("Ignoring catalog event with routing key: %s", msg.RoutingKey)
	}

	return nil
}

// processCategoryCreated handles category creation events
func (c *CatalogConsumer) processCategoryCreated(ctx context.Context, data []byte) error {
	var event catalogContracts.CategoryCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CategoryCreatedEvent: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid CategoryCreatedEvent: %w", err)
	}

	// Parse UUID from string
	categoryID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %w", err)
	}

	category := &domain.CachedCategory{
		ID:          categoryID,
		Name:        event.Name,
		Description: event.Description,
		Status:      event.Status,
		Metadata:    event.Metadata,
	}

	if err := c.cache.UpsertCategory(ctx, category); err != nil {
		return fmt.Errorf("failed to upsert category %s: %w", event.ID, err)
	}

	log.Printf("Category created/updated in cache: id=%s, name=%s", event.ID, event.Name)

	return nil
}

// processCategoryUpdated handles category update events
func (c *CatalogConsumer) processCategoryUpdated(ctx context.Context, data []byte) error {
	var event catalogContracts.CategoryUpdatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CategoryUpdatedEvent: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid CategoryUpdatedEvent: %w", err)
	}

	// Parse UUID from string
	categoryID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %w", err)
	}

	category := &domain.CachedCategory{
		ID:          categoryID,
		Name:        event.Name,
		Description: event.Description,
		Status:      event.Status,
		Metadata:    event.Metadata,
	}

	if err := c.cache.UpsertCategory(ctx, category); err != nil {
		return fmt.Errorf("failed to upsert category %s: %w", event.ID, err)
	}

	log.Printf("Category updated in cache: id=%s, name=%s", event.ID, event.Name)

	return nil
}

// processCategoryDeleted handles category deletion events
func (c *CatalogConsumer) processCategoryDeleted(ctx context.Context, data []byte) error {
	var event catalogContracts.CategoryDeletedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CategoryDeletedEvent: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid CategoryDeletedEvent: %w", err)
	}

	// Parse UUID from string
	categoryID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %w", err)
	}

	if err := c.cache.DeleteCategory(ctx, categoryID); err != nil {
		// Delete operations on in-memory cache should not fail
		log.Printf("Error deleting category from cache: id=%s, error=%v", event.ID, err)
	} else {
		log.Printf("Category deleted from cache: id=%s", event.ID)
	}

	return nil
}

// processProductCreated handles product creation events
func (c *CatalogConsumer) processProductCreated(ctx context.Context, data []byte) error {
	var event catalogContracts.ProductCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal ProductCreatedEvent: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid ProductCreatedEvent: %w", err)
	}

	// Parse UUIDs from strings
	productID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid product ID format: %w", err)
	}

	categoryID, err := uuid.Parse(event.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %w", err)
	}

	product := &domain.CachedProduct{
		ID:                productID,
		Name:              event.Name,
		Description:       event.Description,
		CategoryID:        categoryID,
		Duration:          event.Duration,
		Status:            event.Status,
		Availability:      event.Availability,
		BufferTime:        event.BufferTime,
		CancellationHours: event.CancellationHours,
		StripeProductID:   event.StripeProductID,
		Metadata:          event.Metadata,
	}

	if err := c.cache.UpsertProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to upsert product %s: %w", event.ID, err)
	}

	log.Printf("Product created/updated in cache: id=%s, name=%s, categoryId=%s",
		event.ID, event.Name, event.CategoryID)

	return nil
}

// processProductUpdated handles product update events
func (c *CatalogConsumer) processProductUpdated(ctx context.Context, data []byte) error {
	var event catalogContracts.ProductUpdatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal ProductUpdatedEvent: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid ProductUpdatedEvent: %w", err)
	}

	// Parse UUIDs from strings
	productID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid product ID format: %w", err)
	}

	categoryID, err := uuid.Parse(event.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %w", err)
	}

	product := &domain.CachedProduct{
		ID:                productID,
		Name:              event.Name,
		Description:       event.Description,
		CategoryID:        categoryID,
		Duration:          event.Duration,
		Status:            event.Status,
		Availability:      event.Availability,
		BufferTime:        event.BufferTime,
		CancellationHours: event.CancellationHours,
		StripeProductID:   event.StripeProductID,
		Metadata:          event.Metadata,
	}

	if err := c.cache.UpsertProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to upsert product %s: %w", event.ID, err)
	}

	log.Printf("Product updated in cache: id=%s, name=%s, categoryId=%s",
		event.ID, event.Name, event.CategoryID)

	return nil
}

// processProductDeleted handles product deletion events
func (c *CatalogConsumer) processProductDeleted(ctx context.Context, data []byte) error {
	var event catalogContracts.ProductDeletedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal ProductDeletedEvent: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid ProductDeletedEvent: %w", err)
	}

	// Parse UUID from string
	productID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid product ID format: %w", err)
	}

	if err := c.cache.DeleteProduct(ctx, productID); err != nil {
		// Delete operations on in-memory cache should not fail
		log.Printf("Error deleting product from cache: id=%s, error=%v", event.ID, err)
	} else {
		log.Printf("Product deleted from cache: id=%s", event.ID)
	}

	return nil
}
