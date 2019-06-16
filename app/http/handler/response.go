package handler

import (
	"github.com/vonhraban/secret-server/secret"
)

type secretDefinition struct {
	Hash           string `json:"hash"`
	SecretText     string `json:"secretText"`
	RemainingViews int    `json:"remainingViews"`
	CreatedAt      string `json:"CreatedAt"`
	ExpiresAt      string `json:"ExpiresAt"`
}

type PersistSecretResponse struct {
	secretDefinition
}

func persistSecretResponseFromSecret(secret *secret.Secret) *PersistSecretResponse {
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

func viewSecretResponseFromSecret(secret *secret.Secret) *ViewSecretResponse {
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
