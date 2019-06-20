package handler

import (
	"errors"
	"net/http"
	"strconv"
)

type persistSecretRequest struct {
	secret           string
	expireAfterViews int
	expireAfter      int
}

func validate(r *http.Request) error {
	if r.FormValue("secret") == "" {
		return errors.New("Secret is empty")
	}

	if r.FormValue("expireAfterViews") == "" {
		return errors.New("Expire after views is empty")
	}

	if r.FormValue("expireAfter") == "" {
		return errors.New("Expire after is empty")
	}

	if _, err := strconv.Atoi(r.FormValue("expireAfter")); err != nil {
		return errors.New("Expire after is not a valid date")
	}

	return nil
}

// TODO! Rename to BuildPersist...
func persistSecretRequestFromHTTPRequest(r *http.Request) (*persistSecretRequest, error) {
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
