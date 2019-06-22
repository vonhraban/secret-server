package handler

import (
	"fmt"
	"net/http"
	"strconv"
)

type persistSecretRequest struct {
	secret           string
	expireAfterViews int
	expireAfter      int
}

type EmptyValueError struct {
	Field string
}

func NewEmptyValueError(field string) *EmptyValueError{
	return &EmptyValueError{Field: field}
}

func (e *EmptyValueError) Error() string {
	return fmt.Sprintf("Error: %s can not be empty", e.Field)
}

func validate(r *http.Request) error {
	if r.FormValue("secret") == "" {
		return NewEmptyValueError("secret")
	}

	if r.FormValue("expireAfterViews") == "" {
		return NewEmptyValueError("expireAfterViews")
	}

	if r.FormValue("expireAfter") == "" {
		return NewEmptyValueError("expireAfter")
	}

	return nil
}

func buildPersistSecretRequestFromHTTPRequest(r *http.Request) (*persistSecretRequest, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := validate(r); err != nil {
		return nil, err
	}

	secret := r.FormValue("secret")
	expireAfterViews, err := strconv.Atoi(r.FormValue("expireAfterViews"))
	if err != nil {
		return nil, err
	}

	expireAfter, err := strconv.Atoi(r.FormValue("expireAfter"))
	if err != nil {
		return nil, err
	}

	return &persistSecretRequest{
		secret:           secret,
		expireAfterViews: expireAfterViews,
		expireAfter:      expireAfter,
	}, nil
}
