package apperrors

import (
	"fmt"
	"net/http"

)

// CustomError represents a custom error type as a string.
type CustomError string

// Error implements the error interface for CustomError.
// It converts the CustomError to a string and returns it.
func (e CustomError) Error() string {
	return string(e)
}

const (
	BadRequest           = CustomError("Bad request")
	InternalServer       = CustomError("Internal Server")
	FailedToCreateDriver = CustomError("failure to create driver obj")
	MigrationFailure     = CustomError("migrate failure")
)

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
