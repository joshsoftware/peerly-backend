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
	OrganizationConfigNotFound         = CustomError("organizationconfig  not found")
	InvalidReferenceId                 = CustomError("Invalid reference id")
	AttemptExceeded                    = CustomError(" 3 attempts exceeded ")
	InvalidOTP                         = CustomError("invalid otp")
	TimeExceeded                       = CustomError("time exceeded")
	ErrOTPAlreadyExists                = CustomError("otp already exists")
	ErrOTPAttemptsExceeded             = CustomError("attempts exceeded for organization")
	ErrRecordNotFound                  = CustomError("Database record not found")
	OrganizationConfigAlreadyPresent   = CustomError("organization config already present")
	InvalidRewardMultiplier            = CustomError("reward multiplier should greater than 1")
	InvalidRewardQuotaRenewalFrequency = CustomError("reward renewal frequency should greater than 1")
	InvalidTimezone                    = CustomError("enter valid timezone")
	InvalidId                          = CustomError("invalid id")
	InternalServerError                = CustomError("internal server error")
	JSONParsingErrorReq                = CustomError("error in parsing request in json")
	JSONParsingErrorResp               = CustomError("error in parsing response in json")
	OutOfRange                         = CustomError("request value is out of range")
	OrganizationNotFound               = CustomError("organization of given id not found")
	RoleUnathorized                    = CustomError("Role unauthorized")
	PageParamNotFound                  = CustomError("Page parameter not found")
	RepeatedUser                       = CustomError("Repeated user")
	InvalidContactEmail                = CustomError("contact email is already present")
	InvalidDomainName                  = CustomError("domain name is already present")
	InvalidCoreValueData               = CustomError("invalid corevalue data")
	TextFieldBlank                     = CustomError("text field cannot be blank")
	DescFieldBlank                     = CustomError("description cannot be blank")
	InvalidParentValue                 = CustomError("invalid parent core value")
	InvalidOrgId                       = CustomError("invalid organisation")
	UniqueCoreValue                    = CustomError("choose a unique coreValue name")
	InvalidAuthToken                   = CustomError("invalid Auth token")
	IntranetValidationFailed           = CustomError("intranet Validation Failed")
	UserNotFound                       = CustomError("user not found")
	InvalidIntranetData                = CustomError("invalid data recieved from intranet")
	GradeNotFound                      = CustomError("grade not found")
	AppreciationNotFound               = CustomError("appreciation not found")
	BadRequest                         = CustomError("bad request")
	InternalServer                     = CustomError("internal Server")
	FailedToCreateDriver               = CustomError("failure to create driver obj")
	MigrationFailure                   = CustomError("migrate failure")
	SelfAppreciationError              = CustomError("self-appreciation is not allowed")
	InvalidCoreValueID                 = CustomError("invalid corevalue id")
	InvalidReceiverID                  = CustomError("invalid receiver id")
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
	case OrganizationConfigNotFound, OrganizationNotFound, InvalidCoreValueData, InvalidOrgId, AppreciationNotFound, PageParamNotFound:
		return http.StatusNotFound
	case BadRequest, InvalidId, JSONParsingErrorReq, TextFieldBlank, InvalidParentValue, DescFieldBlank, UniqueCoreValue, InvalidIntranetData, IntranetValidationFailed, RepeatedUser, SelfAppreciationError, InvalidCoreValueID, InvalidReceiverID, InvalidRewardMultiplier, InvalidRewardQuotaRenewalFrequency, InvalidTimezone, GradeNotFound:
		return http.StatusBadRequest
	case InvalidAuthToken, RoleUnathorized:
		return http.StatusUnauthorized
	case InvalidContactEmail, InvalidDomainName:
		return http.StatusConflict
	case OrganizationConfigAlreadyPresent:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
