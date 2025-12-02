package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// TODO: this is the thing that I should when creating a product
type CreateProductRequest2 struct {
	Product CreateProductRequest `json:"product"`
	Price   CreatePriceRequest   `json:"price"`
}

type CreateProductRequest struct {
	Name              string           `json:"name"`
	Description       string           `json:"description,omitempty"`
	CategoryID        string           `json:"category"`
	Duration          int              `json:"duration"` // in minutes
	Availability      AvailabilityType `json:"availability"`
	BufferTime        int              `json:"bufferTime"`         // in minutes
	CancellationHours int              `json:"cancellationHours"`  // in hours
	Metadata          map[string]any   `json:"metadata,omitempty"` // Flexible storage for type-specific attributes
}

type UpdateProductRequest struct {
	Name              *string        `json:"name"`
	Description       *string        `json:"description,omitempty"`
	CategoryID        *string        `json:"category"`
	Duration          *int           `json:"duration"` // in minutes
	Status            *string        `json:"publishedStatus"`
	Availability      *string        `json:"availability"`
	BufferTime        *int           `json:"bufferTime"`         // in minutes
	CancellationHours *int           `json:"cancellationHours"`  // in hours
	Metadata          map[string]any `json:"metadata,omitempty"` // Flexible storage for type-specific attributes
}

func (u UpdateProductRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if u.Name != nil && *u.Name == "" {
		errs.Set("name", "Product name cannot be empty.")
	}
	if u.Name != nil && len(*u.Name) > 255 {
		errs.Set("name", "Product name too long.")
	}
	if u.Description != nil && len(*u.Description) > 1000 {
		errs.Set("description", "Product description too long.")
	}
	if u.CategoryID != nil {
		if err := uuid.Validate(*u.CategoryID); err != nil {
			errs.Set("categoryID", fmt.Errorf("Product category ID is not valid: %w", err))
		}
	}
	if u.Duration != nil {
		if *u.Duration < 20 {
			errs.Set("duration", "Duration must be at least 20 minutes.")
		} else if *u.Duration%10 != 0 {
			errs.Set("duration", "Duration must be a multiple of 10 minutes (e.g., 20, 30, 40, 60, 90).")
		}
	}
	if u.Status != nil && !PublishedStatus(*u.Status).IsValid() {
		errs.Set("status", "Invalid product status.")
	}
	if u.Availability != nil && !AvailabilityType(*u.Availability).IsValid() {
		errs.Set("availability", "Invalid product availability.")
	}
	if u.BufferTime != nil && *u.BufferTime < 0 { // Changed to non-negative as per your example
		errs.Set("bufferTime", "Buffer time cannot be negative.")
	}
	if u.CancellationHours != nil && *u.CancellationHours < 0 { // Changed to non-negative as per your example
		errs.Set("cancellationHours", "Cancellation hours cannot be negative.")
	}

	return errs.AsError()
}

type ProductRes struct {
	ID                uuid.UUID        `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description,omitempty"`
	Category          Category         `json:"category"`
	Duration          int              `json:"duration"` // in minutes
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
	Status            PublishedStatus  `json:"publishedStatus"`
	Availability      AvailabilityType `json:"availability"`
	BufferTime        int              `json:"bufferTime"`         // in minutes
	CancellationHours int              `json:"cancellationHours"`  // in hours
	Metadata          map[string]any   `json:"metadata,omitempty"` // Flexible storage for type-specific attributes
}
