package rabbitmq

const (
	// Exchange
	CatalogExchangeName = "catalog.exchange"

	// Routing keys
	ProductCreatedRoutingKey  = "catalog.product.created"
	ProductUpdatedRoutingKey  = "catalog.product.updated"
	ProductDeletedRoutingKey  = "catalog.product.deleted"
	CategoryCreatedRoutingKey = "catalog.category.created"
	CategoryUpdatedRoutingKey = "catalog.category.updated"
	CategoryDeletedRoutingKey = "catalog.category.deleted"

	// Queue names per consuming service
	AuthUserCatalogQueueName = "authuser.catalog.queue"
	// Add more queues here as other services need catalog data
	// BookingCatalogQueueName  = "booking.catalog.queue"
)