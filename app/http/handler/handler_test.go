package handler_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/vonhraban/secret-server/core/log"
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

var _ = Describe("Secret Handler", func() {
	clock := &deterministicClock{}
	vault := persistence.NewInMemoryVault()
	logger := log.NewLogrusLogger("debug")

	var (
		secretHandler        *handler.SecretHandler
		router               *mux.Router
		now                  time.Time
		form                 url.Values
		secretText           string
		expireInMinutes      int
		expireAfterViews     int
		futureExpirationDate time.Time
		pastExpirationDate   time.Time
	)

	BeforeEach(func() {
		secretHandler = handler.NewSecretHandler(vault, clock, logger)
		router = mux.NewRouter()

		var err error
		now, err = time.Parse("2006-01-02 15:04:05", "2019-06-15 11:14:23")
		if err != nil {
			panic(err)
		}
		clock.setCurrentTime(now)

		form = url.Values{}
		secretText = "123abc"
		expireInMinutes = 10
		expireAfterViews = 5

		futureExpirationDate = now.Add(time.Minute * time.Duration(expireInMinutes))
		pastExpirationDate = now.Add(time.Hour * -time.Duration(48))
	})

	Describe("/secret recieves a POST request", func() {
		Context("post request is valid and specified expiration time in minutes", func() {
			It("should save and return the secret", func() {
				// Arrange
				recorder := httptest.NewRecorder()
				h := http.HandlerFunc(secretHandler.Persist)

				form.Add("secret", secretText)
				form.Add("expireAfterViews", strconv.Itoa(expireAfterViews))
				form.Add("expireAfter", strconv.Itoa(expireInMinutes))

				req := httptest.NewRequest("POST", "/v1/secret", strings.NewReader(form.Encode()))
				req.Form = form

				// Action
				h.ServeHTTP(recorder, req)
				var response handler.PersistSecretResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					panic(err)
				}

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(response.SecretText).To(Equal(secretText))
				Expect(response.RemainingViews).To(Equal(expireAfterViews))
				Expect(response.ExpiresAt).To(Equal(futureExpirationDate.Format("2006-01-02 15:04:05")))
			})
		})

		Context("post request is valid and specified expiration time in minutes", func() {
			It("should save and return the secret", func() {
				// Arrange
				expireInMinutes := 0

				recorder := httptest.NewRecorder()
				h := http.HandlerFunc(secretHandler.Persist)

				form.Add("secret", secretText)
				form.Add("expireAfterViews", strconv.Itoa(expireAfterViews))
				form.Add("expireAfter", strconv.Itoa(expireInMinutes))

				req := httptest.NewRequest("POST", "/v1/secret", strings.NewReader(form.Encode()))
				req.Form = form

				// Action
				h.ServeHTTP(recorder, req)
				var response handler.PersistSecretResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					panic(err)
				}

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(response.SecretText).To(Equal(secretText))
				Expect(response.RemainingViews).To(Equal(expireAfterViews))
				Expect(response.ExpiresAt).To(Equal(""))
			})
		})

		Context("post request does not contain secret text", func() {
			It("should give an error saying the secret text is required", func() {
				// Arrange
				recorder := httptest.NewRecorder()
				h := http.HandlerFunc(secretHandler.Persist)

				form.Add("expireAfterViews", strconv.Itoa(expireAfterViews))
				form.Add("expireAfter", strconv.Itoa(expireInMinutes))

				req := httptest.NewRequest("POST", "/v1/secret", strings.NewReader(form.Encode()))
				req.Form = form

				// Action
				h.ServeHTTP(recorder, req)
				var response handler.ErrorResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					panic(err)
				}

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
				Expect(response.Message).To(Equal("Error: secret can not be empty"))
			})
		})

		Context("post request does not contain expire after views", func() {
			It("should give an error saying the expire after views is required", func() {
				// Arrange
				recorder := httptest.NewRecorder()
				h := http.HandlerFunc(secretHandler.Persist)

				form.Add("secret", secretText)
				form.Add("expireAfter", strconv.Itoa(expireInMinutes))

				req := httptest.NewRequest("POST", "/v1/secret", strings.NewReader(form.Encode()))
				req.Form = form

				// Action
				h.ServeHTTP(recorder, req)
				var response handler.ErrorResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					panic(err)
				}

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
				Expect(response.Message).To(Equal("Error: expireAfterViews can not be empty"))
			})
		})

		Context("post request does not contain expire after TTL field", func() {
			It("should give an error saying the expire after views is required", func() {
				// Arrange
				recorder := httptest.NewRecorder()
				h := http.HandlerFunc(secretHandler.Persist)

				form.Add("secret", secretText)
				form.Add("expireAfterViews", strconv.Itoa(expireAfterViews))

				req := httptest.NewRequest("POST", "/v1/secret", strings.NewReader(form.Encode()))
				req.Form = form

				// Action
				h.ServeHTTP(recorder, req)
				var response handler.ErrorResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					panic(err)
				}

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
				Expect(response.Message).To(Equal("Error: expireAfter can not be empty"))
			})
		})

		Context("post request sends accept content header JSON", func() {
			It("should return JSON", func() {
				// Arrange
				recorder := httptest.NewRecorder()
				h := http.HandlerFunc(secretHandler.Persist)

				req := httptest.NewRequest("POST", "/v1/secret", nil)
				req.Header.Add("Accept", "application/json")

				// Action
				h.ServeHTTP(recorder, req)
				var response handler.ErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)

				// Assert
				Expect(err).To(BeNil())
			})
		})

		Context("post request sends accept content header XML", func() {
			It("should return XML", func() {
				// Arrange
				recorder := httptest.NewRecorder()
				h := http.HandlerFunc(secretHandler.Persist)

				req := httptest.NewRequest("POST", "/v1/secret", nil)
				req.Header.Add("Accept", "application/xml")

				// Action
				h.ServeHTTP(recorder, req)
				var response handler.ErrorResponse
				err := xml.NewDecoder(recorder.Body).Decode(&response)

				// Assert
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("/secret/{hash} recieves a GET request", func() {
		Context("a record exists, has more than 0 remaining views and has not expired", func() {
			It("should return the secret", func() {
				// Arrange
				hash := "0a5a98f9-0110-49b1-bd28-4ca10ebae614"

				timeValue, err := time.Parse("2006-01-02 15:04:05", "2019-06-15 11:24:23")
				if err != nil {
					panic(err)
				}

				existingSecret := &secret.Secret{
					Hash:           hash,
					SecretText:     secretText,
					ExpiresAt:      timeValue,
					RemainingViews: expireAfterViews,
				}

				if err = vault.Store(existingSecret); err != nil {
					panic(err)
				}

				router.HandleFunc("/v1/secret/{hash}", secretHandler.View)

				recorder := httptest.NewRecorder()
				url := fmt.Sprintf("/v1/secret/%s", hash)
				req := httptest.NewRequest("GET", url, nil)

				// Action
				router.ServeHTTP(recorder, req)

				// Assert
				var response handler.ViewSecretResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					panic(err)
				}

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(response.SecretText).To(Equal(secretText))
				Expect(response.RemainingViews).To(Equal(expireAfterViews))
				Expect(response.ExpiresAt).To(Equal(futureExpirationDate.Format("2006-01-02 15:04:05")))
			})
		})

		Context("a record exists, has more than 0 remaining views but expired", func() {
			It("should return Not Found", func() {
				// Arrange
				hash := "0a5a98f9-0110-49b1-bd28-4ca10ebae614"

				existingSecret := &secret.Secret{
					Hash:           hash,
					SecretText:     secretText,
					ExpiresAt:      pastExpirationDate,
					RemainingViews: expireAfterViews,
				}

				if err := vault.Store(existingSecret); err != nil {
					panic(err)
				}

				router.HandleFunc("/v1/secret/{hash}", secretHandler.View)

				recorder := httptest.NewRecorder()
				url := fmt.Sprintf("/v1/secret/%s", hash)
				req := httptest.NewRequest("GET", url, nil)

				// Action
				router.ServeHTTP(recorder, req)

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("a record exists, has 0 remaining views and has not yet expired", func() {
			It("should return Not Found", func() {
				// Arrange
				hash := "0a5a98f9-0110-49b1-bd28-4ca10ebae614"
				expireAfterViews := 0

				existingSecret := &secret.Secret{
					Hash:           hash,
					SecretText:     secretText,
					ExpiresAt:      futureExpirationDate,
					RemainingViews: expireAfterViews,
				}

				if err := vault.Store(existingSecret); err != nil {
					panic(err)
				}

				router.HandleFunc("/v1/secret/{hash}", secretHandler.View)

				recorder := httptest.NewRecorder()
				url := fmt.Sprintf("/v1/secret/%s", hash)
				req := httptest.NewRequest("GET", url, nil)

				// Action
				router.ServeHTTP(recorder, req)

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("a record does not exist", func() {
			It("should return Not Found", func() {
				// Arrange
				hash := "0a5a98f9-0110-49b1-bd28-4ca10ebae614"

				router.HandleFunc("/v1/secret/{hash}", secretHandler.View)

				recorder := httptest.NewRecorder()
				url := fmt.Sprintf("/v1/secret/%s", hash)
				req := httptest.NewRequest("GET", url, nil)

				// Action
				router.ServeHTTP(recorder, req)

				// Assert
				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("post request sends accept content header JSON", func() {
			It("should return JSON", func() {
				// Arrange
				hash := "0a5a98f9-0110-49b1-bd28-4ca10ebae614"

				router.HandleFunc("/v1/secret/{hash}", secretHandler.View)

				recorder := httptest.NewRecorder()
				url := fmt.Sprintf("/v1/secret/%s", hash)
				req := httptest.NewRequest("GET", url, nil)
				req.Header.Add("Accept", "application/json")

				// Action
				router.ServeHTTP(recorder, req)

				// Assert
				var response handler.ErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				Expect(err).To(BeNil())
			})
		})

		Context("post request sends accept content header XML", func() {
			It("should return XML", func() {
				// Arrange
				hash := "0a5a98f9-0110-49b1-bd28-4ca10ebae614"

				router.HandleFunc("/v1/secret/{hash}", secretHandler.View)

				recorder := httptest.NewRecorder()
				url := fmt.Sprintf("/v1/secret/%s", hash)
				req := httptest.NewRequest("GET", url, nil)
				req.Header.Add("Accept", "application/xml")

				// Action
				router.ServeHTTP(recorder, req)

				// Assert
				var response handler.ErrorResponse
				err := xml.NewDecoder(recorder.Body).Decode(&response)
				Expect(err).To(BeNil())
			})
		})
	})
})
