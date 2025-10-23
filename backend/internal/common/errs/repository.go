package errs

import (
	"errors"
	"fmt"
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

	// Connection and infrastructure errors
	ErrConnectionFailure  = errors.New("database connection failure")
	ErrTooManyConnections = errors.New("too many database connections")
	ErrQueryCancelled     = errors.New("query was cancelled")
	ErrTransactionFailure = errors.New("transaction failure")
	ErrDeadlock           = errors.New("database deadlock detected")
	ErrPermissionDenied   = errors.New("insufficient database permissions")
	ErrResourceExhausted  = errors.New("database resources exhausted")

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
