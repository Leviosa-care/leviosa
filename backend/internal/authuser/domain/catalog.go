package domain

import (
	"github.com/google/uuid"
)

// CachedCategory represents a simplified category for in-memory caching
// Only essential fields are stored, excluding database-specific fields
type CachedCategory struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"` // Used for filtering - only "published" items are cached
	Metadata    map[string]any   `json:"metadata"`
}

// CachedProduct represents a simplified product for in-memory caching
// Only essential fields are stored, excluding database-specific fields
type CachedProduct struct {
	ID                 uuid.UUID        `json:"id"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	CategoryID         uuid.UUID        `json:"category_id"`
	Duration           int              `json:"duration"`             // Duration in minutes
	Status             string           `json:"status"`               // Used for filtering - only "published" items are cached
	Availability       string           `json:"availability"`         // Availability type
	BufferTime         int              `json:"buffer_time"`          // Buffer time in minutes
	CancellationHours  int              `json:"cancellation_hours"`   // Cancellation notice in hours
	StripeProductID    string           `json:"stripe_product_id"`    // External Stripe product ID
	Metadata           map[string]any   `json:"metadata"`
}

// IsPublished checks if the category has published status
func (c *CachedCategory) IsPublished() bool {
	return c.Status == "published"
}

// IsPublished checks if the product has published status
func (p *CachedProduct) IsPublished() bool {
	return p.Status == "published"
}