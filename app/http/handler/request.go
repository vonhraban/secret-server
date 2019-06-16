package handler

import (
	"net/http"
	"strconv"
)

type persistSecretRequest struct {
	secret           string
	expireAfterViews int
	expireAfter      int
}

func persistSecretRequestFromHTTPRequest(r *http.Request) (*persistSecretRequest, error) {
	r.ParseForm()
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
