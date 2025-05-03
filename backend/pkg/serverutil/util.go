package serverutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const SIGNINENDPOINT = "signin"
const SIGNUPENDPOINT = "signup"
const SIGNOUTENDPOINT = "signout"

type Validator interface {
	Valid(ctx context.Context) (problems error)
}

var ErrValidStruct = errors.New("Failed to valid the struct")
var ErrDecodeJSON = errors.New("Failed to decode JSON payload")

func NewDecodeJSONErr(err error) error {
	return fmt.Errorf("%w: %w", ErrDecodeJSON, err)
}

func NewValidStructErr(structName string) error {
	return fmt.Errorf("%w: %s", ErrValidStruct, structName)
}

func Decode[T any](body io.ReadCloser) (T, error) {
	var res T
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&res); err != nil {
		return res, NewDecodeJSONErr(err)
	}
	return res, nil
}

func DecodeValid[T Validator](ctx context.Context, body io.ReadCloser) (T, error) {
	var v T
	v, err := Decode[T](body)
	if err != nil {
		return v, err
	}
	if err := v.Valid(ctx); err != nil {
		return v, err
	}
	return v, nil
}

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
