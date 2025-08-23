package errs

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrRepositoryNotFound   = errors.New("record not found")
	ErrRepositoryNotCreated = errors.New("record not created")
	ErrRepositoryNotUpdated = errors.New("record not updated")
	ErrRepositoryNotDeleted = errors.New("record not deleted")
	ErrDatabase             = errors.New("general database error")    // Broad DB error
	ErrInternal             = errors.New("repository internal error") // For non-DB related issues within repo
	ErrContext              = errors.New("context related error")
	ErrValidation           = errors.New("input validation failed within repository")   // Better for issues *before* DB interaction if needed
	ErrInvalidInput         = errors.New("invalid input data for repository operation") // For issues like bad JSON marshalling
	ErrNoFieldsForUpdate    = errors.New("no fields provided for update")               // Define this ONCE

	// PostgreSQL specific errors
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
	ErrNotNullViolation    = errors.New("not null constraint violation")
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrCheckViolation      = errors.New("check constraint violation")

	// Wrapper for query execution problems
	ErrDBQuery = errors.New("database query execution error")

	// Error for external storage operations (like S3)
	ErrExternalStorage = errors.New("external storage operation failed")
)

func NewDBQueryErr(err error) error {
	return fmt.Errorf("%w: %w", ErrDBQuery, err)
}

// NewExternalStorageErr wraps an error from an external storage system
func NewExternalStorageErr(err error, operation, key string) error {
	return fmt.Errorf("%s %w for key '%s': %w", operation, ErrExternalStorage, key, err)
}

func NewValidationErr(err error, domainName string) error {
	return fmt.Errorf("%s %w: %w", domainName, ErrValidation, err)
}

func NewContextErr(err error) error {
	return fmt.Errorf("%w: %w", ErrContext, err)
}

func NewInternalErr(err error) error {
	return fmt.Errorf("%w: %w", ErrInternal, err)
}

// NewInvalidInputErr specifically for input issues like metadata marshalling
func NewInvalidInputErr(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidInput, err)
}

func NewRepositoryNotFoundErr(err error, domainName string) error {
	return fmt.Errorf("%s %w: %w", domainName, ErrRepositoryNotFound, err)
}

func NewRepositoryNotCreatedErr(err error, domainName string) error {
	return fmt.Errorf("%s %w: %w", domainName, ErrRepositoryNotCreated, err)
}

func NewRepositoryNotUpdatedErr(err error, domainName string) error {
	return fmt.Errorf("%s %w: %w", domainName, ErrRepositoryNotUpdated, err)
}

func NewRepositoryNotDeletedErr(err error, domainName string) error {
	return fmt.Errorf("%s %w: %w", domainName, ErrRepositoryNotDeleted, err)
}

func NewDatabaseErr(err error) error {
	return fmt.Errorf("%w: %w", ErrDatabase, err)
}

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
