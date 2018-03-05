package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redhat/notification-microservice/domain"
)

// StatusHandler will return the status of the app
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, AppStatus{OK: true})
}

// EmailHandler used to handle emails
func EmailHandler(emailServer domain.EmailServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			appErr := NewAppError(err.Error(), VALIDATION, Low)
			RespondWithError(w, appErr)
			return
		}

		var emailData *domain.Email
		if err := json.Unmarshal(bodyBytes, &emailData); err != nil || emailData == nil {
			appErr := NewAppError(err.Error(), VALIDATION, Low)
			RespondWithError(w, appErr)
			return
		}

		if len(emailData.Recipients) == 0 {
			appErr := NewAppError("No recipients given for the email", VALIDATION, Low)
			RespondWithError(w, appErr)
			return
		}

		if len(emailData.Body) == 0 {
			appErr := NewAppError("No text body given for the email", VALIDATION, Low)
			RespondWithError(w, appErr)
			return
		}

		if err := emailServer.SendEmail(emailData); err != nil {
			appErr := NewAppError(err.Error(), UNKNOWN, Low)
			RespondWithError(w, appErr)
			return
		}
		// send back to sender
		RespondWithJSON(w, http.StatusOK, AppStatus{OK: true})
	}
}

// JiraHandler used to handle emails
func JiraHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			appErr := NewAppError(err.Error(), VALIDATION, Low)
			RespondWithError(w, appErr)
			return
		}

		var jiraData domain.JiraTicket
		if err := json.Unmarshal(bodyBytes, &jiraData); err != nil {
			appErr := NewAppError(err.Error(), VALIDATION, Low)
			RespondWithError(w, appErr)
			return
		}

		jiraResp, err := jiraData.Send()
		if err != nil {
			appErr := NewAppError(err.Error(), UNKNOWN, Low)
			RespondWithError(w, appErr)
			return
		}

		// send back to sender
		RespondWithJSON(w, http.StatusOK, jiraResp)
	}
}
