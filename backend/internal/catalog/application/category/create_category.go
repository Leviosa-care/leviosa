package category

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *CategoryService) CreateCategory(ctx context.Context, request *domain.CreateCategoryRequest) (string, error) {
	if err := request.Valid(ctx); err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}

	now := time.Now()

	category := &domain.Category{
		ID:          uuid.New(),
		Name:        strings.ToLower(request.Name),
		Description: request.Description,
		Metadata:    request.Metadata,
		Status:      domain.Draft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	categoryID, err := s.repo.AddCategory(ctx, category)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			// This specifically catches the `json.Marshal` error from the repository
			return "", errs.NewInvalidValueErr(fmt.Sprintf("category metadata: %v", err))
		case errors.Is(err, errs.ErrUniqueViolation): // This is the core change: handling the DB unique constraint directly
			return "", errs.NewAlreadyExistsError(err, "category with this name")
		case errors.Is(err, errs.ErrNotNullViolation):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("missing required data for category: %v", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			// Unlikely for category creation, but good to have.
			return "", errs.NewInvalidValueErr(fmt.Sprintf("invalid foreign key for category: %v", err))
		case errors.Is(err, errs.ErrCheckViolation):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("category data failed check constraint: %v", err))
		case errors.Is(err, errs.ErrDBQuery): // Catch general DB query execution errors
			return "", errs.NewQueryFailedErr(fmt.Errorf("repository query failed for category: %w", err))
		case errors.Is(err, errs.ErrDatabase): // More general database issues (e.g., connection)
			return "", errs.NewUnexpectedError(fmt.Errorf("database connection error for category: %w", err))
		case errors.Is(err, errs.ErrContext): // Handle context cancellation/timeout
			return "", errs.NewUnexpectedError(fmt.Errorf("context error during category creation: %w", err))
		default:
			return "", errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during category creation: %w", err))
		}
	}
	return categoryID, nil
}
