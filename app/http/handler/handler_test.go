package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
)

var _ = Describe("Secret Handler", func() {
	vault := persistence.NewInMemoryVault()
	clock := &secret.TimeClock{}
	secretHandler := &handler.SecretHandler{
		Vault: vault,
		Clock: clock,
	}

	Context("Given a user wants to persist a new secret", func() {
		// TODO! Factory
		recorder := httptest.NewRecorder() // TODO! Factory
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
							var response secret.Secret
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

							It("And has an expiration date of not later than 10 minutes from now", func() {
								//panic("Not implemented")
							})
						})

					})
				})
			})
		})
	})
})
