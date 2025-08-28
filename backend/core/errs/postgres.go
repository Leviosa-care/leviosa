package errs

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ClassifyPgError maps specific PgErrors to your sentinel errors
func ClassifyPgError(operation string, err error) error {
	if err == nil {
		return nil
	}

	// Handle pgx.ErrNoRows specifically
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%s: %w", operation, ErrRepositoryNotFound)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		var repoErr error
		switch pgErr.Code {
		// Constraint violations (23xxx)
		case "23505": // Unique violation
			repoErr = fmt.Errorf("%w: %w", ErrUniqueViolation, pgErr)
		case "23503": // Foreign key violation
			repoErr = fmt.Errorf("%w: %w", ErrForeignKeyViolation, pgErr)
		case "23502": // Not null violation
			repoErr = fmt.Errorf("%w: %w", ErrNotNullViolation, pgErr)
		case "23514": // Check violation
			repoErr = fmt.Errorf("%w: %w", ErrCheckViolation, pgErr)

		// Data exceptions (22xxx)
		case "22P02", "22001": // invalid_text_representation, string_data_right_truncation
			repoErr = fmt.Errorf("%w: %w", ErrInvalidInput, pgErr)
		case "22003": // numeric_value_out_of_range
			repoErr = fmt.Errorf("%w: %w", ErrInvalidInput, pgErr)
		case "22023": // invalid_parameter_value
			repoErr = fmt.Errorf("%w: %w", ErrInvalidInput, pgErr)
		case "22025": // invalid_escape_sequence
			repoErr = fmt.Errorf("%w: %w", ErrInvalidInput, pgErr)

		// Connection exceptions (08xxx)
		case "08000", "08003", "08006": // connection_exception, connection_does_not_exist, connection_failure
			repoErr = fmt.Errorf("%w: %w", ErrConnectionFailure, pgErr)
		case "08P01": // protocol_violation
			repoErr = fmt.Errorf("%w: %w", ErrConnectionFailure, pgErr)

		// Transaction rollback (40xxx)
		case "40001": // serialization_failure
			repoErr = fmt.Errorf("%w: %w", ErrTransactionFailure, pgErr)
		case "40P01": // deadlock_detected
			repoErr = fmt.Errorf("%w: %w", ErrDeadlock, pgErr)

		// Syntax and access errors (42xxx)
		case "42601": // syntax_error
			repoErr = fmt.Errorf("%w: %w", ErrInternal, pgErr)
		case "42P01": // undefined_table
			repoErr = fmt.Errorf("%w: %w", ErrInternal, pgErr)
		case "42703": // undefined_column
			repoErr = fmt.Errorf("%w: %w", ErrInternal, pgErr)
		case "42501": // insufficient_privilege
			repoErr = fmt.Errorf("%w: %w", ErrPermissionDenied, pgErr)
		case "42P06": // duplicate_schema
			repoErr = fmt.Errorf("%w: %w", ErrInternal, pgErr)

		// Resource exhaustion (53xxx)
		case "53100": // disk_full
			repoErr = fmt.Errorf("%w: %w", ErrResourceExhausted, pgErr)
		case "53200": // out_of_memory
			repoErr = fmt.Errorf("%w: %w", ErrResourceExhausted, pgErr)
		case "53300": // too_many_connections
			repoErr = fmt.Errorf("%w: %w", ErrTooManyConnections, pgErr)

		// Query cancellation (57xxx)
		case "57014": // query_canceled
			repoErr = fmt.Errorf("%w: %w", ErrQueryCancelled, pgErr)
		case "57P01": // admin_shutdown
			repoErr = fmt.Errorf("%w: %w", ErrConnectionFailure, pgErr)

		// Authentication errors (28xxx)
		case "28000": // invalid_authorization_specification
			repoErr = fmt.Errorf("%w: %w", ErrPermissionDenied, pgErr)
		case "28P01": // invalid_password
			repoErr = fmt.Errorf("%w: %w", ErrPermissionDenied, pgErr)

		default:
			// For any other specific PgError, wrap it with a general database error
			repoErr = fmt.Errorf("%w: %w", ErrDatabase, pgErr)
		}
		// Always wrap the classified error with a clear message indicating the operation.
		return fmt.Errorf("%s: %w", operation, repoErr)
	}

	// If it's not a pgconn.PgError, it's a general database error or something else
	// (e.g., a context error or a network issue).
	return fmt.Errorf("%s: %w: %w", operation, ErrDBQuery, err)
}
