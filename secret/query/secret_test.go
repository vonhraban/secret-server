package query_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/secret/query"
)

// Set up time travel
type deterministicClock struct {
	now time.Time
}

func (d *deterministicClock) setCurrentTime(time time.Time) {
	d.now = time
}

func (d *deterministicClock) GetCurrentTime() time.Time {
	return d.now
}

// The actual fun
var _ = Describe("Secret", func() {
	// clock := &secret.TimeClock{}

	Context("Given it is 2019-06-15 11:14:00 now", func() {
		clock := &deterministicClock{}
		timeValue, err := time.Parse("2006-01-02 15:04:05", "2019-06-15 11:14:00")
		if err != nil {
			panic(err)
		}
		clock.setCurrentTime(timeValue)

		Context("And given secret 123abc exists with allowed max views of 5 and expiration time of 0 minutes", func() {
			vault := persistence.NewInMemoryVault(clock)
			hash := "49885756-2af3-4f9c-85c6-c4b0d9006e2b"
			secretToStore := &secret.Secret{
				Hash:           hash,
				SecretText:     "123abc",
				RemainingViews: 5,
				CreatedAt:      clock.GetCurrentTime(),
			}

			if err = vault.Store(secretToStore); err != nil {
				panic(err)
			}

			Context("When I retrieve the secret", func() {
				q := query.NewGetSecretQuery(vault, hash)
				storedSecret, err := q.Execute()
				if err != nil {
					panic(err)
				}

				It("should contain secret text 123abc", func() {
					Expect(storedSecret.SecretText).To(Equal("123abc"))
				})

				It("should have 5 remaining views", func() {
					Expect(storedSecret.RemainingViews).To(Equal(5))
				})
			})
		})
	})
})
