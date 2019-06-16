package secret

import "time"

type Secret struct {
	Hash           string `json:"hash"`
	SecretText     string `json:"secretText"`
	RemainingViews int `json:"remainingViews"`
	CreatedAt      time.Time `json:"CreatedAt"`
	ExpiresAt      time.Time `json:"ExpiresAt"`
}
