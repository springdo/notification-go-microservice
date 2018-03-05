package api

import (
	"fmt"
	"net/http"
)

// AppError - hold all the information regarding an error for an application
type AppError struct {
	Message  string
	Type     Type
	Severity Severity
}

// NewAppError - creates a new instance of an application error
func NewAppError(message string, errorType Type, errorSeverity Severity) AppError {
	return AppError{
		Message:  message,
		Type:     errorType,
		Severity: errorSeverity,
	}
}

func (appErr *AppError) Error() string {
	return fmt.Sprintf("(%s) %s", severityNames[appErr.Severity], appErr.Message)
}

// GetHTTPStatusCode - returns the status code applicable for the type of an application error
func (appErr *AppError) GetHTTPStatusCode() int {
	switch appErr.Type {
	case VALIDATION:
		return http.StatusBadRequest
	case AUTHORISATION:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// Type - enum for the type of application error, more will be added
type Type uint

const (
	//VALIDATION - for invalid input
	VALIDATION Type = iota
	//AUTHORISATION - for unauthorised input
	AUTHORISATION
	//UNKNOWN - for everything else
	UNKNOWN
)

// Severity - enum for the severity level of application error
type Severity uint

const (
	//Low - severity level
	Low Severity = iota
	//Medium - severity level
	Medium
	//High - severity level
	High
)

var severityNames = [...]string{"Low", "Medium", "High"}

func (severity Severity) String() string {
	return severityNames[severity]
}
