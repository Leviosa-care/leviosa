package catalog

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CategoryCreatedEvent represents a new category being created
type CategoryCreatedEvent struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"` // published, draft, archived
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// CategoryUpdatedEvent represents a category being updated
type CategoryUpdatedEvent struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// CategoryDeletedEvent represents a category being deleted
type CategoryDeletedEvent struct {
	ID string `json:"id"`
}

// ProductCreatedEvent represents a new product being created
type ProductCreatedEvent struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	CategoryID        string         `json:"categoryId"`
	Duration          int            `json:"duration"`
	Status            string         `json:"status"`
	Availability      string         `json:"availability"` // online, in-person, hybrid
	BufferTime        int            `json:"bufferTime"`
	CancellationHours int            `json:"cancellationHours"`
	StripeProductID   string         `json:"stripeProductId"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

// ProductUpdatedEvent represents a product being updated
type ProductUpdatedEvent struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	CategoryID        string         `json:"categoryId"`
	Duration          int            `json:"duration"`
	Status            string         `json:"status"`
	Availability      string         `json:"availability"`
	BufferTime        int            `json:"bufferTime"`
	CancellationHours int            `json:"cancellationHours"`
	StripeProductID   string         `json:"stripeProductId"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

// ProductDeletedEvent represents a product being deleted
type ProductDeletedEvent struct {
	ID string `json:"id"`
}

// Helper to validate UUID in events
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// Validation methods
func (e CategoryCreatedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid category ID: %s", e.ID)
	}
	if e.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	return nil
}

func (e CategoryUpdatedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid category ID: %s", e.ID)
	}
	if e.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	return nil
}

func (e CategoryDeletedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid category ID: %s", e.ID)
	}
	return nil
}

func (e ProductCreatedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid product ID: %s", e.ID)
	}
	if !isValidUUID(e.CategoryID) {
		return fmt.Errorf("invalid category ID: %s", e.CategoryID)
	}
	if e.Name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	if e.Duration <= 0 {
		return fmt.Errorf("product duration must be positive")
	}
	return nil
}

func (e ProductUpdatedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid product ID: %s", e.ID)
	}
	if !isValidUUID(e.CategoryID) {
		return fmt.Errorf("invalid category ID: %s", e.CategoryID)
	}
	if e.Name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	if e.Duration <= 0 {
		return fmt.Errorf("product duration must be positive")
	}
	return nil
}

func (e ProductDeletedEvent) Validate() error {
	if !isValidUUID(e.ID) {
		return fmt.Errorf("invalid product ID: %s", e.ID)
	}
	return nil
}