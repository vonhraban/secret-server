package secret

import "time"

type Secret struct {
	Hash           string
	SecretText     string
	RemainingViews int
	CreatedAt      time.Time
	ExpiresAt      time.Time
}
