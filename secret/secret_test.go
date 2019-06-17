package secret_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
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

		Context("When a secret 123abc is added with allowed max views of 5 and expiration time of 9 minutes", func() {
			vault := persistence.NewInMemoryVault(clock)
			hash := "cfeb626e-f945-47f1-9ec3-1a066273c733"
			cmd := secret.NewAddSecretCommand(vault, clock, hash, "123abc", 5, 9)
			err := cmd.Execute()
			if err != nil {
				panic(err)
			}
			Context("Then this secret should be stored", func() {
				storedSecret, err := vault.Retrieve(hash)
				if err != nil {
					panic(err)
				}
				It("has a 5 remaining views", func() {
					Expect(storedSecret.RemainingViews).To(Equal(5))
				})
				It("has the time created set to 2019-06-15 11:14:00", func() {
					expectedTime, err := time.Parse("2006-01-02 15:04:05", "2019-06-15 11:14:00")
					if err != nil {
						panic(err)
					}
					Expect(storedSecret.CreatedAt).To(Equal(expectedTime))
				})

				It("has the time expires set to 2019-06-15 12:14:00", func() {
					expectedTime, err := time.Parse("2006-01-02 15:04:05", "2019-06-15 11:23:00")
					if err != nil {
						panic(err)
					}
					Expect(storedSecret.ExpiresAt).To(Equal(expectedTime))
				})
			})
		})

		Context("When a secret 123abc is added with allowed max views of 5 and expiration time of 0 minutes", func() {
			vault := persistence.NewInMemoryVault(clock)
			hash := "1212cf75-2fb5-4df2-a730-1fb9fc63b93b"
			cmd := secret.NewAddSecretCommand(vault, clock, hash, "123abc", 5, 0)
			err := cmd.Execute()
			if err != nil {
				panic(err)
			}
			Context("Then this secret should be stored", func() {
				storedSecret, err := vault.Retrieve(hash)
				if err != nil {
					panic(err)
				}

				It("has no expiration time", func() {
					Expect(storedSecret.ExpiresAt.IsZero()).To(Equal(true))
				})
			})
		})

		Context("And given secret 123abc exists with allowed max views of 5 and expiration time of 0 minutes", func() {
			vault := persistence.NewInMemoryVault(clock)
			hash := "49885756-2af3-4f9c-85c6-c4b0d9006e2b"
			secretToStore := &secret.Secret{
				Hash:           hash,
				SecretText:     "123abc",
				RemainingViews: 5,
				CreatedAt:      clock.GetCurrentTime(),
			}
			_, err := vault.Store(secretToStore)
			if err != nil {
				panic(err)
			}

			Context("When I decrease the available views for this secret", func() {
				cmd := secret.NewDecreaseRemainingViewsCommand(vault, hash)
				if err := cmd.Execute(); err != nil {
					panic(err)
				}

				It("should now have only 4 remaining views", func() {
					storedSecret, err := vault.Retrieve(hash)
					if err != nil {
						panic(err)
					}
					Expect(storedSecret.RemainingViews).To(Equal(4))
				})
			})
		})

		Context("And given secret 123abc exists with allowed max views of 5 and expiration time of 0 minutes", func() {
			vault := persistence.NewInMemoryVault(clock)
			hash := "49885756-2af3-4f9c-85c6-c4b0d9006e2b"
			secretToStore := &secret.Secret{
				Hash:           hash,
				SecretText:     "123abc",
				RemainingViews: 5,
				CreatedAt:      clock.GetCurrentTime(),
			}
			_, err := vault.Store(secretToStore)
			if err != nil {
				panic(err)
			}

			Context("When I retrieve the secret", func() {
				cmd := secret.NewGetSecretQuery(vault, hash)
				storedSecret, err := cmd.Execute()
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
