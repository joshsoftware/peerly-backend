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
	AppreciationNotFound        = CustomError("appreciation not found")
	RoleUnathorized             = CustomError("Role unauthorized")
	PageParamNotFound           = CustomError("Page parameter not found")
	RepeatedUser                = CustomError("Repeated user")
	SelfAppreciationError       = CustomError("user cannot give appreciation to ourself")
	CannotReportOwnAppreciation = CustomError("You cannot report your own appreciations")
	RepeatedReport              = CustomError("You cannot report an appreciation twice")
	InvalidCoreValueID          = CustomError("invalid corevalue id")
	InvalidReceiverID           = CustomError("invalid receiver id")
	UserAlreadyPresent          = CustomError("user already present")
	RewardAlreadyPresent        = CustomError("reward already present")
	RewardQuotaIsNotSufficient  = CustomError("reward quota is not sufficient")
	InvalidRewardPoint          = CustomError("invalid reward point")
	SelfRewardError             = CustomError("user cannot give reward to ourself")
	SelfAppreciationRewardError = CustomError("user cannot give reward to his posted appreciaiton ")
	InvalidId                   = CustomError("invalid id")
	InternalServerError         = CustomError("internal server error")
	JSONParsingErrorReq         = CustomError("error in parsing request in json")
	JSONParsingErrorResp        = CustomError("error in parsing response in json")
	OutOfRange                  = CustomError("request value is out of range")
	OrganizationNotFound        = CustomError("organization of given id not found")
	InvalidContactEmail         = CustomError("contact email is already present")
	InvalidDomainName           = CustomError("domain name is already present")
	InvalidCoreValueData        = CustomError("invalid corevalue data")
	TextFieldBlank              = CustomError("text field cannot be blank")
	DescFieldBlank              = CustomError("description cannot be blank")
	InvalidParentValue          = CustomError("invalid parent core value")
	InvalidOrgId                = CustomError("invalid organisation")
	UniqueCoreValue             = CustomError("choose a unique coreValue name")
	InvalidAuthToken            = CustomError("invalid Auth token")
	IntranetValidationFailed    = CustomError("intranet Validation Failed")
	UserNotFound                = CustomError("user not found")
	InvalidIntranetData         = CustomError("invalid data recieved from intranet")
	GradeNotFound               = CustomError("grade not found")
	BadRequest                  = CustomError("bad request")
	InternalServer              = CustomError("internal Server")
	FailedToCreateDriver        = CustomError("failure to create driver obj")
	MigrationFailure            = CustomError("migrate failure")
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
	case OrganizationNotFound, InvalidOrgId, GradeNotFound, AppreciationNotFound, PageParamNotFound, InvalidIntranetData:
		return http.StatusNotFound
	case BadRequest, InvalidId, JSONParsingErrorReq, TextFieldBlank, InvalidCoreValueData, InvalidParentValue, DescFieldBlank, UniqueCoreValue, IntranetValidationFailed, RepeatedUser, SelfAppreciationError, CannotReportOwnAppreciation, RepeatedReport, InvalidCoreValueID, InvalidReceiverID, InvalidRewardPoint:
		return http.StatusBadRequest
	case InvalidAuthToken, IntranetValidationFailed:
		return http.StatusUnauthorized
	case InvalidContactEmail, InvalidDomainName, UserAlreadyPresent, RewardAlreadyPresent:
		return http.StatusConflict
	case InvalidAuthToken, RoleUnathorized:
		return http.StatusUnauthorized
	case RewardQuotaIsNotSufficient:
		return http.StatusUnprocessableEntity

	default:
		return http.StatusInternalServerError
	}
}
