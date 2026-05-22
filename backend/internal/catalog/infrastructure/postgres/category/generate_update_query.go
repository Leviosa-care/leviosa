package categoryRepository

import (
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func generateUpdateQuery(categoryID string, req *domain.UpdateCategoryRequest) (string, []any, error) {
	sets := []string{}
	args := []any{}
	argCounter := 1

	if req.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, *req.Name)
		argCounter++
	}
	if req.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argCounter))
		args = append(args, *req.Description)
		argCounter++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, *req.Status)
		argCounter++
	}

	if len(sets) == 0 {
		return "", nil, errs.ErrNoFieldsForUpdate
	}

	query := fmt.Sprintf("UPDATE catalog.categories SET %s WHERE id = $%d;",
		strings.Join(sets, ", "), argCounter)
	args = append(args, categoryID)

	return query, args, nil
}
