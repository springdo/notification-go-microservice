package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//RespondWithError - standard method for returning JSON error responses from application
func RespondWithError(w http.ResponseWriter, appError AppError) {
	fmt.Println(appError.Error())
	RespondWithJSON(w, appError.GetHTTPStatusCode(), appError.Error())
}

//RespondWithJSON - standard method for returning JSON responses from application
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if payload == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.WriteHeader(code)
	w.Write(response)
}
