package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/vonhraban/secret-server/persistence"
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

var _ = Describe("Secret Handler", func() {
	vault := persistence.NewInMemoryVault()
	clock := &deterministicClock{}
	secretHandler := handler.NewSecretHandler(vault, clock)

	Context("Given it is 2019-06-15 11:14:23", func() {
		timeValue, err := time.Parse("2006-01-02 15:04:05", "2019-06-15 11:14:23")
		if err != nil {
			panic(err)
		}
		clock.setCurrentTime(timeValue)

		Context("Given a user wants to persist a new secret", func() {
			recorder := httptest.NewRecorder()
			h := http.HandlerFunc(secretHandler.Persist)
			form := url.Values{}

			Context("That has a secret text of abc123", func() {
				form.Add("secret", "abc123")

				Context("And that is expiring after 5 views", func() {
					form.Add("expireAfterViews", "5")

					Context("And that is expiring after 10 minutes", func() {
						form.Add("expireAfter", "10")

						When("When the request is sent", func() {
							req := httptest.NewRequest("POST", "/v1/secret", strings.NewReader(form.Encode()))
							req.Form = form
							h.ServeHTTP(recorder, req)
							Context("Then the secret needs to be returned", func() {
								var response handler.PersistSecretResponse
								// TODO! I should not use domain models here but instead a response object
								json.NewDecoder(recorder.Body).Decode(&response)
								//defer recorder.Body.Close()
								It("And the status code needs to be 200", func() {
									Expect(recorder.Code).To(Equal(http.StatusOK))
								})

								It("And has a secret text of abc123", func() {
									Expect(response.SecretText).To(Equal("abc123"))
								})

								It("And has an expiration after 5 views", func() {
									Expect(response.RemainingViews).To(Equal(5))
								})

								It("And has an expiration date of 2019-06-15 11:24:23", func() {
									Expect(response.ExpiresAt).To(Equal("2019-06-15 11:24:23"))
								})
							})

						})
					})
				})
			})
		})
	})
})
