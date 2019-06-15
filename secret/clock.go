package secret

import "time"

type Clock interface {
	GetCurrentTime() time.Time
}

type TimeClock struct{}

func (c *TimeClock) GetCurrentTime() time.Time {
	return time.Now()
}
