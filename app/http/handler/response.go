package handler

import (
	"encoding/xml"

	"github.com/vonhraban/secret-server/secret"
)

type secretDefinition struct {
	XMLName        xml.Name `json:"-" xml:"Secret"`
	Hash           string   `json:"hash" xml:"hash"`
	SecretText     string   `json:"secretText" xml:"secretText"`
	RemainingViews int      `json:"remainingViews" xml:"remainingViews"`
	CreatedAt      string   `json:"CreatedAt" xml:"CreatedAt"`
	ExpiresAt      string   `json:"ExpiresAt" xml:"ExpiresAt"`
}

type PersistSecretResponse struct {
	secretDefinition
}

type ErrorResponse struct {
	Message string `json:"message" xml:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func NewErrorResponse(message string) *ErrorResponse{
	return &ErrorResponse{Message: message}
}

func persistSecretResponseFromSecret(secret secret.Secret) *PersistSecretResponse {
	timeExpires := ""
	if !secret.ExpiresAt.IsZero() {
		timeExpires = secret.ExpiresAt.Format("2006-01-02 15:04:05")
	}
	return &PersistSecretResponse{
		secretDefinition{
			Hash:           secret.Hash,
			SecretText:     secret.SecretText,
			RemainingViews: secret.RemainingViews,
			CreatedAt:      secret.CreatedAt.Format("2006-01-02 15:04:05"),
			ExpiresAt:      timeExpires,
		},
	}
}

type ViewSecretResponse struct {
	secretDefinition
}

func viewSecretResponseFromSecret(secret secret.Secret) *ViewSecretResponse {
	timeExpires := ""
	if !secret.ExpiresAt.IsZero() {
		timeExpires = secret.ExpiresAt.Format("2006-01-02 15:04:05")
	}
	return &ViewSecretResponse{
		secretDefinition{
			Hash:           secret.Hash,
			SecretText:     secret.SecretText,
			RemainingViews: secret.RemainingViews,
			CreatedAt:      secret.CreatedAt.Format("2006-01-02 15:04:05"),
			ExpiresAt:      timeExpires,
		},
	}
}
