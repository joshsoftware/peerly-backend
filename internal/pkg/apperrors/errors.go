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
	BadRequest                         = CustomError("Bad request")
	InternalServer                     = CustomError("Internal server error")
	FailedToCreateDriver               = CustomError("failure to create driver obj")
	MigrationFailure                   = CustomError("migrate failure")
	OutOfRange                         = CustomError("request value is out of range")
	OrganizationConfigNotFound          = CustomError("organizationconfig  not found")
	InvalidContactEmail                = CustomError("Contact email is already present")
	InvalidDomainName                  = CustomError("Domain name is already present")
	InvalidReferenceId                 = CustomError("Invalid reference id")
	AttemptExceeded                    = CustomError(" 3 attempts exceeded ")
	InvalidOTP                         = CustomError("invalid otp")
	TimeExceeded                       = CustomError("time exceeded")
	ErrOTPAlreadyExists                = CustomError("otp already exists")
	ErrOTPAttemptsExceeded             = CustomError("attempts exceeded for organization")
	InvalidId                          = CustomError("Invalid id")
	InernalServer                      = CustomError("Failed to write organization db")
	JSONParsingErrorReq                = CustomError("error in parsing request in json")
	ErrRecordNotFound                  = CustomError("Database record not found")
	InternalServerError                = CustomError("Internal server error")
	JSONParsingErrorResp               = CustomError("error in parsing response in json")
	InvalidCoreValueData               = CustomError("Invalid corevalue data")
	TextFieldBlank                     = CustomError("Text field cannot be blank")
	DescFieldBlank                     = CustomError("Description cannot be blank")
	InvalidParentValue                 = CustomError("Invalid parent core value")
	InvalidOrgId                       = CustomError("Invalid organisation")
	UniqueCoreValue                    = CustomError("Choose a unique coreValue name")
	InvalidAuthToken                   = CustomError("Invalid Auth token")
	IntranetValidationFailed           = CustomError("Intranet Validation Failed")
	UserNotFound                       = CustomError("User not found")
	InvalidIntranetData                = CustomError("Invalid data recieved from intranet")
	GradeNotFound                      = CustomError("Grade not found")
	AppreciationNotFound               = CustomError("appreciation not found")
	RoleUnathorized                    = CustomError("Role unauthorized")
	OrganizationConfigAlreadyPresent   = CustomError("organization config already present")
	PageParamNotFound                  = CustomError("Page parameter not found")
	RepeatedUser                       = CustomError("Repeated user")
	InvalidCoreValueID                 = CustomError("invalid corevalue id")
	InvalidReceiverID                  = CustomError("invalid receiver id")
	SelfAppreciationError              = CustomError("user cannot give appreciation to ourself")
	InvalidRewardMultiplier            = CustomError("reward multiplier should greater than 1")
	InvalidRewardQuotaRenewalFrequency = CustomError("reward renewal frequency should greater than 1")
	InvalidTimezone                    = CustomError("enter valid timezone")
)

// ErrKeyNotSet - Returns error object specific to the key value passed in
func ErrKeyNotSet(key string) (err error) {
	return fmt.Errorf("key not set: %s", key)
}

// GetHTTPStatusCode returns status code according to customerror and default returns InternalServer error
func GetHTTPStatusCode(err error) int {
	switch err {
	case InternalServerError, JSONParsingErrorResp, InvalidIntranetData:
		return http.StatusInternalServerError
	case OrganizationConfigNotFound, InvalidCoreValueData, InvalidParentValue, InvalidOrgId, GradeNotFound, AppreciationNotFound, PageParamNotFound:
		return http.StatusNotFound
	case BadRequest, InvalidId, JSONParsingErrorReq, TextFieldBlank, DescFieldBlank, UniqueCoreValue, IntranetValidationFailed, RepeatedUser, SelfAppreciationError, InvalidCoreValueID, InvalidReceiverID, InvalidRewardMultiplier, InvalidRewardQuotaRenewalFrequency, InvalidTimezone:
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
