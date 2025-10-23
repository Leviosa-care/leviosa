package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type Product struct {
	ID                uuid.UUID        `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description,omitempty"`
	CategoryID        uuid.UUID        `json:"category"`
	Duration          int              `json:"duration"` // in minutes
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
	Status            PublishedStatus  `json:"publishedStatus"`
	Availability      AvailabilityType `json:"availability"`
	BufferTime        int              `json:"bufferTime"`                // in minutes
	CancellationHours int              `json:"cancellationHours"`         // in hours
	StripeProductID   string           `json:"stripeProductId,omitempty"` // ID of the product in Stripe
	Metadata          map[string]any   `json:"metadata,omitempty"`        // Flexible storage for type-specific attributes
}

func (p Product) Valid(ctx context.Context) error {
	var errs errsx.Map
	// Required fields validation
	if p.Name == "" {
		errs.Set("name", "Product name cannot be empty.")
	}
	// Numeric field validation
	if p.Duration <= 0 {
		errs.Set("duration", "Duration must be a positive value.")
	}
	if p.BufferTime < 0 { // Changed to non-negative as per your example
		errs.Set("bufferTime", "Buffer time cannot be negative.")
	}
	if p.CancellationHours < 0 {
		errs.Set("cancellationHours", "Cancellation hours cannot be negative.")
	}
	// Enum field validation using the IsValid() methods
	if !p.Status.IsValid() {
		errs.Set("publishedStatus", fmt.Sprintf("Invalid published status: '%s'. Must be 'published', 'draft', or 'archived'.", p.Status))
	}
	if !p.Availability.IsValid() {
		errs.Set("availability", fmt.Sprintf("Invalid availability type: '%s'. Must be 'online', 'in-person', or 'hybrid'.", p.Availability))
	}

	// TODO: Add more specific metadata validation if needed (e.g., if category "massage"
	// requires a "massage_type" key in Metadata, validate it here).
	// This might involve fetching category metadata too.
	// For example:
	// if p.CategoryID == "some_massage_category_id" {
	//     if _, ok := p.Metadata["massage_type"]; !ok {
	//         errs.Set("metadata", "Massage products require a 'massage_type' in metadata.")
	//     }
	// }
	return errs.AsError()
}
