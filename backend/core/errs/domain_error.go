package errs

import (
	"errors"
	"fmt"
)

var (
	// Generic errors that can be wrapped by more specific ones
	ErrQueryFailed      = errors.New("database query execution failed")
	ErrUnexpectedError  = errors.New("unexpected error occurred")
	ErrDomainNotFound   = errors.New("resssource not found")
	ErrDomainNotUpdated = errors.New("resssource not updated")
	ErrDomainNotCreated = errors.New("resssource not created")
	ErrDomainNotDeleted = errors.New("resssource not deleted")
	ErrExternalService  = errors.New("external service failure")
	ErrConflict         = errors.New("conflict with existing data") // General conflict

	// Specific business/validation errors
	ErrAlreadyExists       = errors.New("ressource already exists")
	ErrInvalidValue        = errors.New("invalid value")
	ErrValueMismatch       = errors.New("value mismatch")
	ErrAccountLocked       = errors.New("account is locked")
	ErrExpiredToken        = errors.New("token has expired")
	ErrRateLimit           = errors.New("rate limit exceeded")
	ErrCategoryHasProducts = errors.New("category has associated products")
	ErrUnauthorized        = errors.New("unauthorized action") // This might be a global domain error

	// Data serialization/deserialization errors (less common to originate here, more in handler/repo)
	ErrMarshalJSON   = errors.New("json marshalling")
	ErrUnmarshalJSON = errors.New("json unmarshalling")
	ErrNotEncrypted  = errors.New("not encrypted")
	ErrNotDecrypted  = errors.New("not decrypted")

	// Input/Format specific errors
	ErrParsing = errors.New("parsing error")
	ErrFormat  = errors.New("format error")
)

func NewConflictErr(err error) error {
	return fmt.Errorf("%w: %w", ErrConflict, err)
}

func NewCategoryHasProductsErr() error {
	return fmt.Errorf("%w: %w", ErrConflict, ErrCategoryHasProducts)
}

func NewExternalServiceErr(err error, message string) error {
	return fmt.Errorf("%s: %w: %w", message, ErrExternalService, err)
}

// Wrapper functions for domain errors
func NewAlreadyExistsError(err error, resourceName string) error {
	return fmt.Errorf("%w: %s already exists: %w", ErrAlreadyExists, resourceName, err)
}

func NewParsingError(domain string, err error) error {
	return fmt.Errorf("%w: parsing %s: %w", ErrParsing, domain, err)
}

func NewFormatError(domain string, err error) error {
	return fmt.Errorf("%w: format error for %s: %w", ErrFormat, domain, err)
}

func NewRateLimitErr(err error, domainName string) error {
	return fmt.Errorf("%s %w: %w", domainName, ErrRateLimit, err)
}

func NewNotEncryptedErr(resource string, err error) error {
	return fmt.Errorf("%s %w: %w", resource, ErrNotEncrypted, err)
}

func NewNotDecryptedErr(resource string, err error) error {
	return fmt.Errorf("%s %w: %w", resource, ErrNotDecrypted, err)
}

func NewInvalidValueErr(description string) error {
	return fmt.Errorf("%w: %s", ErrInvalidValue, description)
}

func NewLockedAccountErr(err error, name string) error {
	return fmt.Errorf("[%s] - %w: %w", name, ErrAccountLocked, err)
}

func NewValueMismatchErr(storedValue, providedValue any) error {
	return fmt.Errorf("%w: stored value=%v, provided value=%v", ErrValueMismatch, storedValue, providedValue)
}

func NewExpiredTokenErr(name string, err error) error {
	return fmt.Errorf("%w: %s token: %w", ErrExpiredToken, name, err)
}

func NewNotFoundErr(err error, resourceName string) error {
	return fmt.Errorf("%w: %s: %w", ErrDomainNotFound, resourceName, err)
}

func NewNotCreatedErr(err error, resourceName string) error {
	return fmt.Errorf("%w: %s: %w", ErrDomainNotCreated, resourceName, err)
}

func NewNotDeletedErr(err error, resourceName string) error {
	return fmt.Errorf("%w: %s: %w", ErrDomainNotDeleted, resourceName, err)
}

func NewNotUpdatedErr(err error, resourceName string) error {
	return fmt.Errorf("%w: %s: %w", ErrDomainNotUpdated, resourceName, err)
}

func NewJSONMarshalErr(err error) error {
	return fmt.Errorf("%w: %w", ErrMarshalJSON, err)
}

func NewJSONUnmarshalErr(err error) error {
	return fmt.Errorf("%w: %w", ErrUnmarshalJSON, err)
}

func NewQueryFailedErr(err error) error {
	return fmt.Errorf("%w: %w", ErrQueryFailed, err)
}

func NewUnexpectedError(err error) error {
	return fmt.Errorf("%w: %w", ErrUnexpectedError, err)
}
