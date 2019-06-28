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
var _ = Describe("Secret Queries", func() {

	clock := &deterministicClock{}
	vault := persistence.NewInMemoryVault(clock)

	var (
		now                  time.Time
		hash                 string
		secretText           string
		expireInMinutes      int
		expireAfterViews     int
		futureExpirationDate time.Time
		pastExpirationDate   time.Time
	)

	BeforeEach(func() {
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
		pastExpirationDate = now.Add(time.Hour * -time.Duration(48))
	})

	Describe("Secret is queried", func() {
		Context("secret exists, has more than 0 remaining views and not expired", func() {
			It("should return the secret", func() {
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

				// Action
				q := query.NewGetSecretQuery(vault, hash)
				storedSecret, err := q.Execute()
				if err != nil {
					panic(err)
				}

				// Assert
				Expect(storedSecret.Hash).To(Equal(secretToStore.Hash))
				Expect(storedSecret.SecretText).To(Equal(secretToStore.SecretText))
				Expect(storedSecret.RemainingViews).To(Equal(secretToStore.RemainingViews))
				Expect(storedSecret.ExpiresAt).To(Equal(secretToStore.ExpiresAt))
			})
		})

		Context("secret exists, has no (zero) remaining views left and not expired", func() {
			It("should return the secret", func() {
				// Arrange
				expireAfterViews := 0

				secretToStore := &secret.Secret{
					Hash:           hash,
					SecretText:     secretText,
					RemainingViews: expireAfterViews,
					ExpiresAt:      futureExpirationDate,
				}

				if err := vault.Store(secretToStore); err != nil {
					panic(err)
				}

				// Action
				q := query.NewGetSecretQuery(vault, hash)
				storedSecret, err := q.Execute()

				// Assert
				Expect(storedSecret).To(BeNil())
				Expect(err).Should(MatchError(secret.SecretNotFoundError))
			})
		})

		Context("secret exists, has more than 0 remaining views left but expired", func() {
			It("should return the secret", func() {
				// Arrange
				secretToStore := &secret.Secret{
					Hash:           hash,
					SecretText:     secretText,
					RemainingViews: expireAfterViews,
					ExpiresAt:      pastExpirationDate,
				}

				if err := vault.Store(secretToStore); err != nil {
					panic(err)
				}

				// Action
				q := query.NewGetSecretQuery(vault, hash)
				storedSecret, err := q.Execute()

				// Assert
				Expect(storedSecret).To(BeNil())
				Expect(err).Should(MatchError(secret.SecretNotFoundError))
			})
		})
	})
})
