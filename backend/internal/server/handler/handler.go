package handler

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("Requested resource not found.")
	ErrInternal          = errors.New("An internal error occurred:")
	ErrBadRequest        = errors.New("This is bad request.")
	ErrServiceUnvailable = errors.New("Service unavailable.")
	ErrForbidden         = errors.New("Forbidden error")    // authorization error, it is permanent
	ErrUnauthorized      = errors.New("Unauthorized error") // authenticiation error, you need to authenticate, use that in the auth middleware
)

func NewUnauthorizedErr(err error) string {
	return fmt.Sprintf("%s: %s", ErrUnauthorized, err)
}

func NewForbiddenErr(err error) string {
	return fmt.Sprintf("%s: %s", ErrForbidden, err)
}

func NewNotFoundErr(err error) string {
	return fmt.Sprintf("%s: %s", ErrNotFound, err)
}

func NewInternalErr(err error) string {
	return fmt.Sprintf("%s: %s", ErrInternal, err)
}

func NewBadRequestErr(err error) string {
	return fmt.Sprintf("%s: %s", ErrBadRequest, err)
}

func NewServiceUnavailableErr(err error) string {
	return fmt.Sprintf("%s: %s", ErrServiceUnvailable, err)
}

var (
	ErrGet    = errors.New("find")
	ErrCreate = errors.New("create")
	ErrUpdate = errors.New("modify")
	ErrDelete = errors.New("remove")
)

func NewGetErr(domain string, err error) string {
	return fmt.Sprintf("%s %s: %s", ErrGet, domain, err.Error())
}

func NewCreateErr(domain string, err error) string {
	return fmt.Sprintf("%s %s: %s", ErrCreate, domain, err.Error())
}
func NewUpdateErr(domain string, err error) string {
	return fmt.Sprintf("%s %s: %s", ErrUpdate, domain, err.Error())
}

func NewDeleteErr(domain string, err error) string {
	return fmt.Sprintf("%s %s: %s", ErrDelete, domain, err.Error())
}
