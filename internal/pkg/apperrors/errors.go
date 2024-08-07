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
	RoleUnathorized                    = CustomError("Role unauthorized")
	PageParamNotFound                  = CustomError("Page parameter not found")
	RepeatedUser                       = CustomError("Repeated user")
	CannotReportOwnAppreciation        = CustomError("You cannot report your own appreciations")
	RepeatedReport                     = CustomError("You cannot report an appreciation twice")
	UserAlreadyPresent                 = CustomError("user already present")
	RewardAlreadyPresent               = CustomError("reward already present")
	RewardQuotaIsNotSufficient         = CustomError("reward quota is not sufficient")
	InvalidRewardPoint                 = CustomError("invalid reward point")
	SelfRewardError                    = CustomError("user cannot give reward to ourself")
	SelfAppreciationRewardError        = CustomError("user cannot give reward to his posted appreciaiton ")
	InvalidId                          = CustomError("invalid id")
	InternalServerError                = CustomError("internal server error")
	JSONParsingErrorReq                = CustomError("error in parsing request in json")
	JSONParsingErrorResp               = CustomError("error in parsing response in json")
	OutOfRange                         = CustomError("request value is out of range")
	OrganizationNotFound               = CustomError("organization of given id not found")
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
	InvalidPassword                    = CustomError("invalid password")
	InvalidEmail                       = CustomError("invalid email")
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
	NegativeGradePoints                = CustomError("grade points cannot be negative")
	NegativeBadgePoints                = CustomError("badge reward points cannot be negative")
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
	case OrganizationConfigNotFound, OrganizationNotFound, InvalidOrgId, GradeNotFound, AppreciationNotFound, PageParamNotFound, InvalidCoreValueData, InvalidIntranetData:
		return http.StatusNotFound
	case BadRequest, InvalidId, JSONParsingErrorReq, TextFieldBlank, InvalidParentValue, DescFieldBlank, UniqueCoreValue, SelfAppreciationError, CannotReportOwnAppreciation, RepeatedReport, InvalidCoreValueID, InvalidReceiverID, InvalidRewardMultiplier, InvalidRewardQuotaRenewalFrequency, InvalidTimezone, InvalidRewardPoint, InvalidEmail, InvalidPassword, NegativeGradePoints:
		return http.StatusBadRequest
	case InvalidContactEmail, InvalidDomainName, UserAlreadyPresent, RewardAlreadyPresent, RepeatedUser:
		return http.StatusConflict
	case InvalidAuthToken, RoleUnathorized, IntranetValidationFailed:
		return http.StatusUnauthorized
	case RewardQuotaIsNotSufficient:
		return http.StatusUnprocessableEntity
	case OrganizationConfigAlreadyPresent:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
