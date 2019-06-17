package secret

import "time"

type Secret struct {
	Hash           string
	SecretText     string
	RemainingViews int
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func (s *Secret) CanBeSeen(now time.Time) bool {
	return s.RemainingViews > 0 && false == s.isExpired(now)
}

func (s *Secret) isExpired(now time.Time) bool {
	if s.ExpiresAt.IsZero() {
		return false
	}

	if now.Before(s.ExpiresAt) {
		return false
	}

	return true
}