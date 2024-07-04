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
	BadRequest               = CustomError("Bad request")
	InternalServer           = CustomError("Failed to write organization db")
	FailedToCreateDriver     = CustomError("failure to create driver obj")
	MigrationFailure         = CustomError("migrate failure")
	InvalidId                = CustomError("Invalid id")
	InternalServerError      = CustomError("Internal server error")
	JSONParsingErrorReq      = CustomError("error in parsing request in json")
	JSONParsingErrorResp     = CustomError("error in parsing response in json")
	OutOfRange               = CustomError("request value is out of range")
	OrganizationNotFound     = CustomError("organization of given id not found")
	InvalidContactEmail      = CustomError("Contact email is already present")
	InvalidDomainName        = CustomError("Domain name is already present")
	InvalidCoreValueData     = CustomError("Invalid corevalue data")
	TextFieldBlank           = CustomError("Text field cannot be blank")
	DescFieldBlank           = CustomError("Description cannot be blank")
	InvalidParentValue       = CustomError("Invalid parent core value")
	InvalidOrgId             = CustomError("Invalid organisation")
	UniqueCoreValue          = CustomError("Choose a unique coreValue name")
	InvalidAuthToken         = CustomError("Invalid Auth token")
	IntranetValidationFailed = CustomError("Intranet Validation Failed")
	UserNotFound             = CustomError("User not found")
	InvalidIntranetData      = CustomError("Invalid data recieved from intranet")
	GradeNotFound            = CustomError("Grade not found")
)

// ErrKeyNotSet - Returns error object specific to the key value passed in
func ErrKeyNotSet(key string) (err error) {
	return fmt.Errorf("key not set: %s", key)
}

func GetHTTPStatusCode(err error) int {
	switch err {
	case InternalServerError, JSONParsingErrorResp, InvalidIntranetData:
		return http.StatusInternalServerError
	case OrganizationNotFound, InvalidCoreValueData, InvalidParentValue, InvalidOrgId, GradeNotFound:
		return http.StatusNotFound
	case InvalidId, JSONParsingErrorReq, TextFieldBlank, DescFieldBlank, UniqueCoreValue, InvalidAuthToken, IntranetValidationFailed:
		return http.StatusBadRequest
	case InvalidContactEmail, InvalidDomainName:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
