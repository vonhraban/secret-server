package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/core/log"

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
	clock := &deterministicClock{}
	vault := persistence.NewInMemoryVault(clock)
	secretHandler := handler.NewSecretHandler(vault, clock, log.NewLogrusLogger(logrus.New()))

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

	Context("Given there is secret abc123 stored with an expiration date of 2019-06-15 11:24:23 and 1 remaining view under the 0a5a98f9-0110-49b1-bd28-4ca10ebae614 hash", func() {
		timeValue, err := time.Parse("2006-01-02 15:04:05", "2019-06-15 11:24:23")
		if err != nil {
			panic(err)
		}

		existingSecret := &secret.Secret{
			Hash:           "0a5a98f9-0110-49b1-bd28-4ca10ebae614",
			SecretText:     "abc123",
			ExpiresAt:      timeValue,
			RemainingViews: 1,
		}

		if err = vault.Store(existingSecret); err != nil {
			panic(err)
		}

		Context("And a user wants to view secret", func() {
			recorder := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/v1/secret/{hash}", secretHandler.View)

			When("When the request is sent", func() {
				req := httptest.NewRequest("GET", "/v1/secret/0a5a98f9-0110-49b1-bd28-4ca10ebae614", nil)
				r.ServeHTTP(recorder, req)
				Context("Then the secret needs to be returned", func() {
					var response handler.ViewSecretResponse
					// TODO! I should not use domain models here but instead a response object
					json.NewDecoder(recorder.Body).Decode(&response)
					//defer recorder.Body.Close()
					It("And the status code needs to be 200", func() {
						Expect(recorder.Code).To(Equal(http.StatusOK))
					})

					It("And has a secret text of abc123", func() {
						Expect(response.SecretText).To(Equal("abc123"))
					})

					It("And has 0 remaining views", func() {
						Expect(response.RemainingViews).To(Equal(0))
					})

					It("And has an expiration date of 2019-06-15 11:24:23", func() {
						Expect(response.ExpiresAt).To(Equal("2019-06-15 11:24:23"))
					})
				})
			})
		})
	})

	Context("Given there is secret abc123 stored with an expiration date of 2019-06-15 11:24:23 and 0 remaining views under the 0a5a98f9-0110-49b1-bd28-4ca10ebae614 hash", func() {
		timeValue, err := time.Parse("2006-01-02 15:04:05", "2019-06-15 11:24:23")
		if err != nil {
			panic(err)
		}

		existingSecret := &secret.Secret{
			Hash:           "0a5a98f9-0110-49b1-bd28-4ca10ebae614",
			SecretText:     "abc123",
			ExpiresAt:      timeValue,
			RemainingViews: 0,
		}

		if err = vault.Store(existingSecret); err != nil {
			panic(err)
		}

		Context("And a user wants to view secret", func() {
			recorder := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/v1/secret/{hash}", secretHandler.View)

			When("When the request is sent", func() {
				req := httptest.NewRequest("GET", "/v1/secret/0a5a98f9-0110-49b1-bd28-4ca10ebae614", nil)
				r.ServeHTTP(recorder, req)
				Context("Then no secret must be returned since all views are used up", func() {
					It("And the status code needs to be 404", func() {
						Expect(recorder.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})
	})

	Context("Given there is secret abc123 stored with an expiration date of 2019-06-10 11:11:11 and 5 remaining views under the c9d4b534-e2de-43da-ae08-31820a7b83f4 hash", func() {
		timeValue, err := time.Parse("2006-01-02 15:04:05", "2019-06-10 11:11:11")
		if err != nil {
			panic(err)
		}

		existingSecret := &secret.Secret{
			Hash:           "c9d4b534-e2de-43da-ae08-31820a7b83f4",
			SecretText:     "abc123",
			ExpiresAt:      timeValue,
			RemainingViews: 5,
		}

		if err = vault.Store(existingSecret); err != nil {
			panic(err)
		}

		Context("And a user wants to view secret", func() {
			recorder := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/v1/secret/{hash}", secretHandler.View)

			When("When the request is sent", func() {
				req := httptest.NewRequest("GET", "/v1/secret/c9d4b534-e2de-43da-ae08-31820a7b83f4", nil)
				r.ServeHTTP(recorder, req)
				Context("Then no secret must be returned since it is expired", func() {
					It("And the status code needs to be 404", func() {
						Expect(recorder.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})
	})

})
