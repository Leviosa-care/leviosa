package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ParentID    uuid.UUID  `json:"parent_id" db:"parent_id"`
	ParentType  ParentType `json:"parent_type" db:"parent_type"` // e.g., "product", "category"
	Title       string     `json:"title" db:"title"`
	S3Key       string     `json:"s3_key" db:"s3_key"`
	Size        int64      `json:"size" db:"size"`
	ContentType string     `json:"content_type" db:"content_type"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// ParentType defines allowed values for parent types
type ParentType string

const (
	CategoryType ParentType = "category"
	ProductType  ParentType = "product"
)

// IsValid checks if the AvailabilityType is one of the defined constants.
func (it ParentType) IsValid() bool {
	switch strings.ToLower(string(it)) { // Case-insensitive check
	case strings.ToLower(string(CategoryType)), strings.ToLower(string(ProductType)):
		return true
	}
	return false
}
