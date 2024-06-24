package apperrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

// CustomError represents a custom error type as a string.
type CustomError string

// Error implements the error interface for CustomError.
// It converts the CustomError to a string and returns it.
func (e CustomError) Error() string {
    return string(e)
}

// ErrorStruct - a generic struct you can use to create error messages/logs to be converted
// to JSON or other types of messages/data as you need it
type ErrorStruct struct {
	Message string `json:"message,omitempty"` // Your message to the end user or developer
	Status  int    `json:"status,omitempty"`  // HTTP status code that should go with the message/log (if any)
}

// JSONError - This function writes out an error response with the status
// header passed in
func JSONError(rw http.ResponseWriter, status int, err error) {
	// Create the ErrorStruct object for later use
	errObj := ErrorStruct{
		Message: err.Error(),
		Status:  status,
	}

	errJSON, err := json.Marshal(&errObj)
	if err != nil {
		log.Warn(err, "Error in AppErrors marshalling JSON", err)
	}
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(errJSON)
}

const (
	BadRequest           = CustomError("Bad request")
	InternalServer       = CustomError("Failed to write organization db")
	FailedToCreateDriver = CustomError("failure to create driver obj")
	MigrationFailure     = CustomError("migrate failure")
)

// helper functions
func ErrorResp(rw http.ResponseWriter, err error) {
	// Create the ErrorStruct object for later use
	statusCode := GetHTTPStatusCode(err)
	errObj := ErrorStruct{
		Message: err.Error(),
		Status:  statusCode,
	}

	errJSON, err := json.Marshal(&errObj)
	if err != nil {
		log.Warn(err, "Error in AppErrors marshalling JSON", err)
	}
	rw.WriteHeader(statusCode)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(errJSON)
}

// ErrKeyNotSet - Returns error object specific to the key value passed in
func ErrKeyNotSet(key string) (err error) {
	return fmt.Errorf("key not set: %s", key)
}

func GetHTTPStatusCode(err error) int {
	switch err {
	case InternalServer, FailedToCreateDriver, MigrationFailure:
		return http.StatusInternalServerError
	case BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
