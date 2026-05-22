package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type Category struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      PublishedStatus `json:"status"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func (c Category) Valid(ctx context.Context) error {
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
	if !c.Status.IsValid() {
		errs.Set("status", "invalid category status.")
	}
	return errs.AsError()
}

// TODO: create these categories in the database just to try things out
// "massage",
// "wellness",
// "mental coaching",
