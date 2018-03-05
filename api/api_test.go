package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "redhat/notification-microservice/api"
	"redhat/notification-microservice/domain"
)

type TestEmailServer struct {
	message string
}

func (tes *TestEmailServer) SendEmail(emailData *domain.Email) error {
	return nil
}

var _ = Describe("GET - /status", func() {
	Context("", func() {
		It("should return with OK true", func() {
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/status", nil)
			if err != nil {
				Fail(err.Error())
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(StatusHandler)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				Fail(fmt.Sprintf("handler returned wrong status code: got %d want %d",
					status, http.StatusOK))
			}

			// Check the response body is what we expect.
			expected := `{"OK":true}`
			Expect(rr.Body.String()).To(Equal(expected))
		})
	})
})

var _ = Describe("POST - /api/email", func() {
	Context("", func() {
		It("should return with OK true", func() {

			emailJson := `{"recipients": [ "thirdpartybuilderz@mailinator.com" ], "subject": "Test Order", "body": "This is an awesome email!"}`

			reader := strings.NewReader(emailJson) //Convert string to reader
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", "/api/email", reader)
			if err != nil {
				Fail(err.Error())
			}

			testServer := &TestEmailServer{}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(EmailHandler(testServer))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				Fail(fmt.Sprintf("handler returned wrong status code: got %d want %d",
					status, http.StatusOK))
			}

			// Check the response body is what we expect.
			expected := `{"OK":true}`
			Expect(rr.Body.String()).To(Equal(expected))
		})
	})
})

var _ = Describe("POST - /api/email", func() {
	Context("No recipients given", func() {
		It("should return with bad request", func() {

			emailJson := `{"subject": "Test Order", "body": "This is an awesome email!"}`

			reader := strings.NewReader(emailJson) //Convert string to reader
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", "/api/email", reader)
			if err != nil {
				Fail(err.Error())
			}

			testServer := &TestEmailServer{}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(EmailHandler(testServer))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusBadRequest {
				Fail(fmt.Sprintf("handler returned wrong status code: got %d want %d",
					status, http.StatusOK))
			}

			// Check the response body is what we expect.
			expected := `"(Low) No recipients given for the email"`
			Expect(rr.Body.String()).To(Equal(expected))
		})
	})
})

var _ = Describe("POST - /api/email", func() {
	Context("No email body given", func() {
		It("should return with bad request", func() {

			emailJson := `{"recipients": [ "thirdpartybuilderz@mailinator.com" ], "subject": "Test Order"}`

			reader := strings.NewReader(emailJson) //Convert string to reader
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", "/api/email", reader)
			if err != nil {
				Fail(err.Error())
			}

			testServer := &TestEmailServer{}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(EmailHandler(testServer))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusBadRequest {
				Fail(fmt.Sprintf("handler returned wrong status code: got %d want %d",
					status, http.StatusOK))
			}

			// Check the response body is what we expect.
			expected := `"(Low) No text body given for the email"`
			Expect(rr.Body.String()).To(Equal(expected))
		})
	})
})

var _ = Describe("POST - /api/email", func() {
	Context("No subject given", func() {
		It("should return with OK true", func() {

			emailJson := `{"recipients": [ "thirdpartybuilderz@mailinator.com" ], "body": "This is an awesome email!"}`

			reader := strings.NewReader(emailJson) //Convert string to reader
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", "/api/email", reader)
			if err != nil {
				Fail(err.Error())
			}

			testServer := &TestEmailServer{}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(EmailHandler(testServer))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				Fail(fmt.Sprintf("handler returned wrong status code: got %d want %d",
					status, http.StatusOK))
			}

			// Check the response body is what we expect.
			expected := `{"OK":true}`
			Expect(rr.Body.String()).To(Equal(expected))
		})
	})
})
