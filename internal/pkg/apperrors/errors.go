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
	UserAlreadyPresent                 = CustomError("User already present")
	RewardAlreadyPresent               = CustomError("Reward already present")
	RewardQuotaIsNotSufficient         = CustomError("Reward quota is not sufficient")
	InvalidRewardPoint                 = CustomError("Invalid reward point")
	SelfRewardError                    = CustomError("User cannot give reward to ourself")
	SelfAppreciationRewardError        = CustomError("User cannot give reward to his posted appreciaiton ")
	InvalidId                          = CustomError("Invalid id")
	InternalServerError                = CustomError("Internal server error")
	JSONParsingErrorReq                = CustomError("Error in parsing request in json")
	JSONParsingErrorResp               = CustomError("Error in parsing response in json")
	OutOfRange                         = CustomError("Request value is out of range")
	OrganizationNotFound               = CustomError("Organization of given id not found")
	InvalidContactEmail                = CustomError("Contact email is already present")
	InvalidDomainName                  = CustomError("Domain name is already present")
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
	AppreciationNotFound               = CustomError("Appreciation not found")
	BadRequest                         = CustomError("Bad request")
	FailedToCreateDriver               = CustomError("Failure to create driver obj")
	MigrationFailure                   = CustomError("Migrate failure")
	SelfAppreciationError              = CustomError("Self-appreciation is not allowed")
	InvalidCoreValueID                 = CustomError("Invalid corevalue id")
	InvalidReceiverID                  = CustomError("Invalid receiver id")
	InvalidPassword                    = CustomError("Invalid password")
	InvalidEmail                       = CustomError("Invalid email")
	OrganizationConfigNotFound         = CustomError("Organizationconfig  not found")
	InvalidReferenceId                 = CustomError("Invalid reference id")
	AttemptExceeded                    = CustomError("3 attempts exceeded ")
	InvalidOTP                         = CustomError("Invalid otp")
	TimeExceeded                       = CustomError("Time exceeded")
	ErrOTPAlreadyExists                = CustomError("Otp already exists")
	ErrOTPAttemptsExceeded             = CustomError("Attempts exceeded for organization")
	InternalServer                     = CustomError("Internal server error")
	ErrRecordNotFound                  = CustomError("Database record not found")
	OrganizationConfigAlreadyPresent   = CustomError("Organization config already present")
	InvalidRewardMultiplier            = CustomError("Reward multiplier should greater than 1")
	InvalidRewardQuotaRenewalFrequency = CustomError("Reward renewal frequency should greater than 1")
	InvalidTimezone                    = CustomError("Enter valid timezone")
	DescriptionLengthExceed            = CustomError("Maximum character length of 500 exceeded")
	InvalidPageSize                    = CustomError("Invalid page size")
	InvalidPage                        = CustomError("Invalid page value")
	NegativeGradePoints                = CustomError("Grade points cannot be negative")
	NegativeBadgePoints                = CustomError("Badge reward points cannot be negative")
	UnauthorizedDeveloper              = CustomError("Unauthorised developer")
	InvalidLoggerLevel                 = CustomError("Invalid Logger Level")
	PreviousQuarterRatingNotAllowed    = CustomError("Reward can be given for current quarter appreciations")
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
	case InvalidLoggerLevel, BadRequest, InvalidId, JSONParsingErrorReq, TextFieldBlank, InvalidParentValue, DescFieldBlank, UniqueCoreValue, SelfAppreciationError, CannotReportOwnAppreciation, RepeatedReport, InvalidCoreValueID, InvalidReceiverID, InvalidRewardMultiplier, InvalidRewardQuotaRenewalFrequency, InvalidTimezone, InvalidRewardPoint, InvalidEmail, InvalidPassword, DescriptionLengthExceed, InvalidPageSize, InvalidPage, NegativeGradePoints, NegativeBadgePoints, PreviousQuarterRatingNotAllowed:
		return http.StatusBadRequest
	case InvalidContactEmail, InvalidDomainName, UserAlreadyPresent, RewardAlreadyPresent, RepeatedUser:
		return http.StatusConflict
	case InvalidAuthToken, RoleUnathorized, IntranetValidationFailed, UnauthorizedDeveloper:
		return http.StatusUnauthorized
	case RewardQuotaIsNotSufficient:
		return http.StatusUnprocessableEntity
	case OrganizationConfigAlreadyPresent:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
