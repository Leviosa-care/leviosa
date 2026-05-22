package productRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *ProductRepository) GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductRes, error) {
	query := fmt.Sprintf(`
	SELECT
		p.id,
		p.name,
		p.description,
		p.duration,
		p.created_at,
		p.updated_at,
		p.status,
		p.availability,
		p.buffer_time,
		p.cancellation_hours,
		p.metadata,
		c.id,
		c.name,
		c.description
	FROM %s.products p
	JOIN %s.categories c ON p.category_id = c.id
	WHERE p.status = $1
	ORDER BY p.created_at DESC
	`, r.schema, r.schema)

	rows, err := r.pool.Query(ctx, query, domain.Published)
	if err != nil {
		return nil, errs.ClassifyPgError("get all published products", err)
	}
	defer rows.Close()

	var products []*domain.ProductRes

	for rows.Next() {
		var (
			pr       domain.ProductRes
			metaProd []byte
			cat      domain.Category
		)

		err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.Description,
			&pr.Duration,
			&pr.CreatedAt,
			&pr.UpdatedAt,
			&pr.Status,
			&pr.Availability,
			&pr.BufferTime,
			&pr.CancellationHours,
			&metaProd,
			&cat.ID,
			&cat.Name,
			&cat.Description,
		)
		if err != nil {
			return nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan row: %w", err))
		}

		if metaProd != nil {
			if err := json.Unmarshal(metaProd, &pr.Metadata); err != nil {
				return nil, errs.NewInvalidInputErr(fmt.Errorf("failed to unmarshal product metadata: %w", err))
			}
		} else {
			pr.Metadata = make(map[string]any)
		}

		pr.Category = cat
		products = append(products, &pr)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NewDBQueryErr(fmt.Errorf("error during rows iteration: %w", err))
	}

	if len(products) == 0 {
		return []*domain.ProductRes{}, nil
	}

	return products, nil
}
