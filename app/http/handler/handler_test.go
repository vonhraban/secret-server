package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vonhraban/secret-server/secret_server/handler"
)

var _ = Describe("All handlers", func() {
	Describe("Hello world", func() {
		Context("GET Request sent to the handler", func() {
			It("should return hello world", func() {
				recorder := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/", nil)
				h := http.HandlerFunc(HelloWorldHandler)
				h.ServeHTTP(recorder, req)
				Expect(strings.TrimSpace(recorder.Body.String())).To(Equal("Hello World"))
			})
		})
	})

	Describe("Hello name", func() {
		Context("Invalid payload POSTed", func() {
			It("should shout", func() {
				recorder := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("rubbish")))
				h := http.HandlerFunc(HelloNameHandler)
				h.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				//Expect(strings.TrimSpace(recorder.Body.String())).To(Equal("Hello World"))
			})
		})

		Context("Valid payload POSTed", func() {
			It("should greet by name", func() {
				recorder := httptest.NewRecorder()
				request := struct {
					Name string `json:"name"`
				}{Name: "John"}
				marshalledReq, _ := json.Marshal(request)
				req := httptest.NewRequest("POST", "/", bytes.NewBuffer(marshalledReq))
				h := http.HandlerFunc(HelloNameHandler)
				h.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(strings.TrimSpace(recorder.Body.String())).To(Equal("Hello John"))
			})
		})
	})
})
