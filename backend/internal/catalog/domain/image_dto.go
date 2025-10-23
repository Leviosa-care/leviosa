package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type CreateImageRequest struct {
	ParentID   string     `json:"parent_id"`
	ParentType ParentType `json:"parent_type"`
	Title      string     `json:"title"`
	IsActive   *bool      `json:"is_active,omitempty"` // in minutes
}

func (r CreateImageRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if _, err := uuid.Parse(r.ParentID); err != nil {
		errs.Set("parent ID", "parent ID should be a valid UUID.")
	}
	if !r.ParentType.IsValid() {
		errs.Set("parent type", "invalid parent type.")
	}
	if r.Title == "" {
		errs.Set("title", "title cannot be empty.")
	}
	if len(r.Title) > 255 {
		errs.Set("title", "title too long.")
	}
	return errs.AsError()
}

type ImageModifierRequest struct {
	ImageID    string `json:"image_id"`
	ParentID   string `json:"parent_id"`
	ParentType string `json:"parent_type"`
}

func (i ImageModifierRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if err := uuid.Validate(i.ImageID); err != nil {
		errs.Set("image ID", "invalid value, must be a valid UUID.")
	}
	if err := uuid.Validate(i.ParentID); err != nil {

		errs.Set("parent ID", "invalid value, must be a valid UUID.")
	}
	if !ParentType(i.ParentType).IsValid() {
		errs.Set("parent type", "invalid parent type value. Must be 'category' or 'product'")
	}
	return errs.AsError()
}
