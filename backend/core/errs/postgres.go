package errs

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// ClassifyPgError maps specific PgErrors to your sentinel errors
func ClassifyPgError(operation string, err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		var repoErr error
		switch pgErr.Code {
		case "23505": // Unique violation
			repoErr = fmt.Errorf("%w: %w", ErrUniqueViolation, pgErr)
		case "23503": // Foreign key violation
			repoErr = fmt.Errorf("%w: %w", ErrForeignKeyViolation, pgErr)
		case "23502": // Not null violation
			repoErr = fmt.Errorf("%w: %w", ErrNotNullViolation, pgErr)
		case "23514": // Check violation
			repoErr = fmt.Errorf("%w: %w", ErrCheckViolation, pgErr)
		case "22P02", "22001": // invalid_text_representation, string_data_right_truncation
			// These are often due to malformed input data
			repoErr = fmt.Errorf("%w: %w", ErrInvalidInput, pgErr)
		case "42601", "42P01": // syntax_error, undefined_table
			// These indicate a bug in our application's SQL query
			repoErr = fmt.Errorf("%w: %w", ErrInternal, pgErr)
		default:
			// For any other specific PgError, wrap it with a general database error.
			repoErr = fmt.Errorf("%w: %w", ErrDatabase, pgErr)
		}
		// Always wrap the classified error with a clear message indicating the operation.
		return fmt.Errorf("%s: %w", operation, repoErr)
	}

	// If it's not a pgconn.PgError, it's a general database error or something else
	// (e.g., a context error or a network issue).
	return fmt.Errorf("%s: %w: %w", operation, ErrDBQuery, err)
}

