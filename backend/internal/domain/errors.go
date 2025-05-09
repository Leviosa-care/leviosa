package domain

import (
	"errors"
	"fmt"
)

var (
	ErrQueryFailed    = errors.New("database query execution failed")
	ErrUnexpectedType = errors.New("unexpected error type")
	ErrNotFound       = errors.New("resssource not found")
	ErrNotUpdated     = errors.New("resssource not updated")
	ErrNotCreated     = errors.New("resssource not created")
	ErrExpiredToken   = errors.New("expired token")
	ErrNotDeleted     = errors.New("resssource not deleted")
	ErrValueMismatch  = errors.New("value mismatch")
	ErrAccountLocked  = errors.New("locked account")
	ErrMarshalJSON    = errors.New("json marshalling")
	ErrUnmarshalJSON  = errors.New("json unmarshalling")
	ErrInvalidValue   = errors.New("invalid value")
	ErrNotEncrypted   = errors.New("not encrypted")
	ErrNotDecrypted   = errors.New("not decrypted")
	ErrRateLimit      = errors.New("rate limit error")
	ErrParsing        = errors.New("parsing error")
	ErrFormat         = errors.New("format error")
)

func NewParsingError(domain string, err error) error {
	return fmt.Errorf("%w: parsing %s: %w", ErrParsing, domain, err)
}

func NewFormatError(domain string, err error) error {
	return fmt.Errorf("%w: parsing %s: %w", ErrFormat, domain, err)
}

func NewRateLimitErr(err error, domainName string) error {
	return fmt.Errorf("%s %w: %w", domainName, ErrRateLimit, err)
}

func NewNotEncryptedErr(ressource string, err error) error {
	return fmt.Errorf("%s %w: %w", ressource, ErrNotEncrypted, err)
}

func NewNotDecryptedErr(ressource string, err error) error {
	return fmt.Errorf("%s %w: %w", ressource, ErrNotDecrypted, err)
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
	return fmt.Errorf("%w: %w - %s", ErrExpiredToken, err, name)
}

func NewNotFoundErr(err error) error {
	return fmt.Errorf("%w: %w", ErrNotFound, err)
}

func NewNotCreatedErr(err error) error {
	return fmt.Errorf("%w: %w", ErrNotCreated, err)
}

func NewNotDeletedErr(err error) error {
	return fmt.Errorf("%w: %w", ErrNotDeleted, err)
}

func NewNotUpdatedErr(err error) error {
	return fmt.Errorf("%w: %w", ErrNotUpdated, err)
}

func NewJSONMarshalErr(err error) error {
	return fmt.Errorf("%w: %w", ErrMarshalJSON, err)
}

func NewJSONUnmarshalErr(err error) error {
	return fmt.Errorf("%w: %w", ErrMarshalJSON, err)
}

func NewQueryFailedErr(err error) error {
	return fmt.Errorf("%w: %w", ErrQueryFailed, err)
}

func NewUnexpectTypeErr(err error) error {
	return fmt.Errorf("%w: %w", ErrUnexpectedType, err)
}
