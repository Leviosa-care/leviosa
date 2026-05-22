package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c CreateCategoryRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if c.Name == "" {
		errs.Set("name", "category name cannot be empty.")
	}
	if len(c.Name) > 255 {
		errs.Set("name", "category name too long.")
	}
	if c.Description == "" {
		errs.Set("description", "category description cannot be empty.")
	}
	if len(c.Description) > 1000 {
		errs.Set("description", "category description too long.")
	}
	return errs.AsError()
}

type UpdateCategoryRequest struct {
	ID          string         `json:"id"`
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Status      *string        `json:"status,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

func (c UpdateCategoryRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if err := uuid.Validate(c.ID); err != nil {
		errs.Set("ID", "Category ID must be a valid UUID.")
	}
	if c.Name != nil && *c.Name == "" {
		errs.Set("name", "Category name cannot be empty.")
	}
	if c.Status != nil && !PublishedStatus(*c.Status).IsValid() {
		errs.Set("status", "Invalid category status.")
	}
	if c.Name != nil && len(*c.Name) > 255 {
		errs.Set("name", "Category name too long.")
	}
	if c.Description != nil && len(*c.Description) > 1000 {
		errs.Set("description", "Category description too long.")
	}
	return errs.AsError()
}

type CategoryRes struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
