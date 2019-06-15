package secret

import "time"

type Secret struct {
	ID string
	Token string
	Uses int
	MaxUses int
	CreatedAt time.Time
	ExpiresAt time.Time
}
