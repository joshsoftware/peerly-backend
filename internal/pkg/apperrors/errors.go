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
	BadRequest                       = CustomError("Bad request")
	InternalServer                   = CustomError("Failed to write organization db")
	FailedToCreateDriver             = CustomError("failure to create driver obj")
	MigrationFailure                 = CustomError("migrate failure")
	OutOfRange                       = CustomError("request value is out of range")
	OrganizationNotFound             = CustomError("organization  not found")
	InvalidContactEmail              = CustomError("Contact email is already present")
	InvalidDomainName                = CustomError("Domain name is already present")
	InvalidReferenceId               = CustomError("Invalid reference id")
	AttemptExceeded                  = CustomError(" 3 attempts exceeded ")
	InvalidOTP                       = CustomError("invalid otp")
	TimeExceeded                     = CustomError("time exceeded")
	ErrOTPAlreadyExists              = CustomError("otp already exists")
	ErrOTPAttemptsExceeded           = CustomError("attempts exceeded for organization")
	InvalidId                        = CustomError("Invalid id")
	InernalServer                    = CustomError("Failed to write organization db")
	JSONParsingErrorReq              = CustomError("error in parsing request in json")
	ErrRecordNotFound                = CustomError("Database record not found")
	InternalServerError              = CustomError("Internal server error")
	JSONParsingErrorResp             = CustomError("error in parsing response in json")
	InvalidCoreValueData             = CustomError("Invalid corevalue data")
	TextFieldBlank                   = CustomError("Text field cannot be blank")
	DescFieldBlank                   = CustomError("Description cannot be blank")
	InvalidParentValue               = CustomError("Invalid parent core value")
	InvalidOrgId                     = CustomError("Invalid organisation")
	UniqueCoreValue                  = CustomError("Choose a unique coreValue name")
	InvalidAuthToken                 = CustomError("Invalid Auth token")
	IntranetValidationFailed         = CustomError("Intranet Validation Failed")
	UserNotFound                     = CustomError("User not found")
	InvalidIntranetData              = CustomError("Invalid data recieved from intranet")
	GradeNotFound                    = CustomError("Grade not found")
	AppreciationNotFound             = CustomError("appreciation not found")
	RoleUnathorized                  = CustomError("Role unauthorized")
	OrganizationConfigAlreadyPresent = CustomError("organization config already present")
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
	case InternalServerError, JSONParsingErrorResp, InvalidIntranetData:
		return http.StatusInternalServerError
	case OrganizationNotFound, InvalidCoreValueData, InvalidParentValue, InvalidOrgId, GradeNotFound, AppreciationNotFound:
		return http.StatusNotFound
	case BadRequest, InvalidId, JSONParsingErrorReq, TextFieldBlank, DescFieldBlank, UniqueCoreValue, IntranetValidationFailed:
		return http.StatusBadRequest
	case InvalidContactEmail, InvalidDomainName:
		return http.StatusConflict
	case InvalidAuthToken, RoleUnathorized:
		return http.StatusUnauthorized
	case OrganizationConfigAlreadyPresent:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
