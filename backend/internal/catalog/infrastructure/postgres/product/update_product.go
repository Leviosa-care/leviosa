package productRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *ProductRepository) UpdateProduct(ctx context.Context, productID uuid.UUID, product *domain.UpdateProductRequest) error {
	query, args, err := generateUpdateQuery(r.schema, productID.String(), product)
	if err != nil {
		// If generateUpdateQuery returns "no fields provided", handle it or let it propagate.
		// If it's a marshalling error, it's an invalid input error from the repo's perspective.
		if errors.Is(err, errs.ErrNoFieldsForUpdate) { // Check for the specific error string (or make it a custom error)
			return errs.NewInvalidInputErr(err)
		}
		return errs.NewInternalErr(fmt.Errorf("failed to generate update query: %w", err))
	}

	commandTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errs.ClassifyPgError("update product", err)
	}

	if commandTag.RowsAffected() == 0 {
		// If 0 rows were affected, it means the product with the given ID was not found.
		return errs.NewRepositoryNotFoundErr(nil, "product")
	}

	return nil
}
