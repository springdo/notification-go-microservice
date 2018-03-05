package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "redhat/notification-microservice/api"
)

var _ = Describe("RespondWithJSON", func() {
	Context("With a valid structure", func() {
		It("should return with OK with valid test struct", func() {

			type testStruct struct {
				Key   string
				Value string
			}

			testVar := testStruct{
				Key:   "TestKey",
				Value: "TestVal",
			}

			rr := httptest.NewRecorder()

			RespondWithJSON(rr, 200, testVar)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				Fail(fmt.Sprintf("handler returned wrong status code: got %d want %d",
					status, http.StatusOK))
			}

			// Check the response body is what we expect.
			expected := `{"Key":"TestKey","Value":"TestVal"}`
			Expect(rr.Body.String()).To(Equal(expected))
		})
	})
})

var _ = Describe("RespondWithJSON", func() {
	Context("With no structure", func() {
		It("should return with OK with valid test struct", func() {
			rr := httptest.NewRecorder()

			RespondWithJSON(rr, 200, nil)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusInternalServerError {
				Fail(fmt.Sprintf("handler returned wrong status code: got %d want %d",
					status, http.StatusOK))
			}

			// Check the response body is what we expect.
			expected := ``
			Expect(rr.Body.String()).To(Equal(expected))
		})
	})
})
