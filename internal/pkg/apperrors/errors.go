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

// Custome errors with errormessage
const (
	InvalidId                = CustomError("invalid id")
	InternalServerError      = CustomError("internal server error")
	JSONParsingErrorReq      = CustomError("error in parsing request in json")
	JSONParsingErrorResp     = CustomError("error in parsing response in json")
	OutOfRange               = CustomError("request value is out of range")
	OrganizationNotFound     = CustomError("organization of given id not found")
	RoleUnathorized          = CustomError("Role unauthorized")
	PageParamNotFound        = CustomError("Page parameter not found")
	RepeatedUser             = CustomError("Repeated user")
	InvalidContactEmail      = CustomError("contact email is already present")
	InvalidDomainName        = CustomError("domain name is already present")
	InvalidCoreValueData     = CustomError("invalid corevalue data")
	TextFieldBlank           = CustomError("text field cannot be blank")
	DescFieldBlank           = CustomError("description cannot be blank")
	InvalidParentValue       = CustomError("invalid parent core value")
	InvalidOrgId             = CustomError("invalid organisation")
	UniqueCoreValue          = CustomError("choose a unique coreValue name")
	InvalidAuthToken         = CustomError("invalid Auth token")
	IntranetValidationFailed = CustomError("intranet Validation Failed")
	UserNotFound             = CustomError("user not found")
	InvalidIntranetData      = CustomError("invalid data recieved from intranet")
	GradeNotFound            = CustomError("grade not found")
	BadRequest               = CustomError("bad request")
	InternalServer           = CustomError("internal Server")
	FailedToCreateDriver     = CustomError("failure to create driver obj")
	MigrationFailure         = CustomError("migrate failure")
)

// ErrKeyNotSet - Returns error object specific to the key value passed in
func ErrKeyNotSet(key string) (err error) {
	return fmt.Errorf("key not set: %s", key)
}

// GetHTTPStatusCode returns status code according to customerror and default returns InternalServer error
func GetHTTPStatusCode(err error) int {
	switch err {
	case InternalServerError, JSONParsingErrorResp:
		return http.StatusInternalServerError
	case OrganizationNotFound, InvalidOrgId, GradeNotFound, PageParamNotFound, InvalidIntranetData:
		return http.StatusNotFound
	case InvalidId, JSONParsingErrorReq, TextFieldBlank, InvalidCoreValueData, InvalidParentValue, DescFieldBlank, UniqueCoreValue, RepeatedUser:
		return http.StatusBadRequest
	case InvalidAuthToken, IntranetValidationFailed, RoleUnathorized:
		return http.StatusUnauthorized
	case InvalidContactEmail, InvalidDomainName:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
