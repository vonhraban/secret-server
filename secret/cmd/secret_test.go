package cmd_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/secret/cmd"
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
var _ = Describe("Secret Commands", func() {

	Describe("Add Secret", func() {

		var (
			vault secret.Vault
			clock *deterministicClock

			now                  time.Time
			hash                 string
			secretText           string
			expireInMinutes      int
			expireAfterViews     int
			futureExpirationDate time.Time
		)

		BeforeEach(func() {
			clock = &deterministicClock{}
			vault = persistence.NewInMemoryVault(clock)

			var err error
			now, err = time.Parse("2006-01-02 15:04:05", "2019-06-15 11:14:23")
			if err != nil {
				panic(err)
			}
			clock.setCurrentTime(now)

			hash = "49885756-2af3-4f9c-85c6-c4b0d9006e2b"
			secretText = "123abc"
			expireInMinutes = 10
			expireAfterViews = 5

			futureExpirationDate = now.Add(time.Minute * time.Duration(expireInMinutes))
		})

		Describe("Command is issued to add a secret", func() {
			Context("it has expires after minutes TTL defined as not zero", func() {
				It("should return the expiration date as now plus TTL", func() {
					// Arrange
					cmd := cmd.NewAddSecretCommand(vault, clock, hash, secretText, expireAfterViews, expireInMinutes)

					// Action
					err := cmd.Execute()
					if err != nil {
						panic(err)
					}

					// Assert
					storedSecret, err := vault.Retrieve(hash)
					Expect(storedSecret).To(Not(BeNil()))
					Expect(err).To(BeNil())

					Expect(storedSecret.Hash).To(Equal(hash))
					Expect(storedSecret.SecretText).To(Equal(secretText))
					Expect(storedSecret.RemainingViews).To(Equal(expireAfterViews))
					Expect(storedSecret.ExpiresAt).To(Equal(futureExpirationDate))
				})
			})
		})

		Context("it has expires after minutes TTL defined as zero", func() {
			It("should return expiration date as zero value", func() {
				// Arrange
				expireInMinutes := 0
				cmd := cmd.NewAddSecretCommand(vault, clock, hash, secretText, expireAfterViews, expireInMinutes)

				// Action
				err := cmd.Execute()
				if err != nil {
					panic(err)
				}

				// Assert
				storedSecret, err := vault.Retrieve(hash)
				Expect(storedSecret).To(Not(BeNil()))
				Expect(err).To(BeNil())

				Expect(storedSecret.Hash).To(Equal(hash))
				Expect(storedSecret.SecretText).To(Equal(secretText))
				Expect(storedSecret.RemainingViews).To(Equal(expireAfterViews))
				Expect(storedSecret.ExpiresAt.IsZero()).To(Equal(true))
			})
		})
	})

	Describe("Deduct Remaining Views Command", func() {

		var (
			clock                *deterministicClock
			vault                secret.Vault
			now                  time.Time
			hash                 string
			secretText           string
			expireInMinutes      int
			expireAfterViews     int
			futureExpirationDate time.Time
		)

		BeforeEach(func() {
			clock = &deterministicClock{}
			vault = persistence.NewInMemoryVault(clock)

			var err error
			now, err = time.Parse("2006-01-02 15:04:05", "2019-06-15 11:14:23")
			if err != nil {
				panic(err)
			}
			clock.setCurrentTime(now)

			hash = "49885756-2af3-4f9c-85c6-c4b0d9006e2b"
			secretText = "123abc"
			expireInMinutes = 10
			expireAfterViews = 5

			futureExpirationDate = now.Add(time.Minute * time.Duration(expireInMinutes))
		})

		Describe("Command is issued to deduct a number of remaining views for a secret", func() {
			Context("the secret in question exists", func() {
				It("should decrease the number of remaining views by 1", func() {
					// Arrange
					secretToStore := &secret.Secret{
						Hash:           hash,
						SecretText:     secretText,
						RemainingViews: expireAfterViews,
						ExpiresAt:      futureExpirationDate,
					}

					if err := vault.Store(secretToStore); err != nil {
						panic(err)
					}

					cmd := cmd.NewDecreaseRemainingViewsCommand(vault, hash)

					// Action
					err := cmd.Execute()
					if err != nil {
						panic(err)
					}

					// Assert
					storedSecret, err := vault.Retrieve(hash)
					Expect(storedSecret).To(Not(BeNil()))
					Expect(err).To(BeNil())

					Expect(storedSecret.Hash).To(Equal(hash))
					Expect(storedSecret.SecretText).To(Equal(secretText))
					// and remaining views should decrease by 1
					Expect(storedSecret.RemainingViews).To(Equal(expireAfterViews - 1))
					Expect(storedSecret.ExpiresAt).To(Equal(futureExpirationDate))
				})
			})

			Context("the secret in question does not exist", func() {
				It("should error", func() {
					// Arrange
					cmd := cmd.NewDecreaseRemainingViewsCommand(vault, hash)

					// Action
					err := cmd.Execute()

					// Assert
					Expect(err).To(Not(BeNil()))
				})
			})
		})
	})
})
