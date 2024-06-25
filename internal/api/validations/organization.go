package validations

import (
	"fmt"
	"regexp"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

const (
	emailRegexString  = `^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`
	domainRegexString = `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`
)

var Hi5QuotaRenewalFrequency = map[string]int{
	"week":    1,
	"month":   1,
	"quarter": 1,
	"year":    1,
}

var emailRegex = regexp.MustCompile(emailRegexString)
var domainRegex = regexp.MustCompile(domainRegexString)

func OrgValidate(org dto.Organization) (errorResponse map[string]dto.ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if org.Name == "" {
		fieldErrors["name"] = "Can't be blank"
	}

	if !emailRegex.MatchString(org.ContactEmail) {
		fieldErrors["email"] = "Please enter a valid email"
	}

	if !domainRegex.MatchString(org.DomainName) {
		fieldErrors["domain_name"] = "Please enter valid domain"
	}

	if org.SubscriptionValidUpto.IsZero() || org.SubscriptionValidUpto.Before(time.Now()) {
		fieldErrors["subscription_valid_upto"] = "Please enter subscription valid upto date"
	}

	if org.Hi5Limit == 0 {
		fieldErrors["hi5_limit"] = "Please enter hi5 limit greater than 0"
	}

	if org.Hi5QuotaRenewalFrequency == "" || Hi5QuotaRenewalFrequency[org.Hi5QuotaRenewalFrequency] == 0 {
		fieldErrors["hi5_quota_renewal_frequency"] = "Please enter valid hi5 renewal frequency"
	}

	if !checkValidTimezone(org.Timezone) {
		fieldErrors["timezone"] = "Please enter valid timezone"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = map[string]dto.ErrorResponse{
		"error": {
			Error: dto.ErrorObject{
				Code:          "invalid_data",
				MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
				Fields:        fieldErrors,
			},
		},
	}

	return
}

func OrgUpdateValidate(org dto.Organization) (errorResponse map[string]dto.ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	fmt.Println("in validation in update")

	if org.ID <= 0 {
		fieldErrors["id"] = "Please enter valid id"
	}

	if org.ContactEmail != "" && !emailRegex.MatchString(org.ContactEmail) {
		fieldErrors["email"] = "Please enter a valid email"
	}

	if org.DomainName != "" && !domainRegex.MatchString(org.DomainName) {
		fieldErrors["domain_name"] = "Please enter valid domain"
	}

	if !org.SubscriptionValidUpto.IsZero() && org.SubscriptionValidUpto.Before(time.Now()) {
		fieldErrors["subscription_valid_upto"] = "Please enter subscription valid upto date"
	}

	if org.Hi5QuotaRenewalFrequency != "" && Hi5QuotaRenewalFrequency[org.Hi5QuotaRenewalFrequency] == 0 {
		fieldErrors["hi5_quota_renewal_frequency"] = "Please enter valid hi5 renewal frequency"
	}

	if org.Timezone != "" {
		if !checkValidTimezone(org.Timezone) {
			fieldErrors["timezone"] = "Please enter valid timezone"
		}
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = map[string]dto.ErrorResponse{
		"error": {
			Error: dto.ErrorObject{
				Code:          "invalid_data",
				MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
				Fields:        fieldErrors,
			},
		},
	}

	return
}

func OTPInfoValidate(otp dto.OTP) (errorResponse map[string]dto.ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)
	fmt.Println("length: ", len(otp.OTPCode))
	if len(otp.OTPCode) != 6 {
		fieldErrors["otp_code"] = "enter 6 digit valid otp code"
	}

	if otp.OrgId <= 0 {
		fieldErrors["id"] = "Please enter valid organization id"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = map[string]dto.ErrorResponse{
		"error": {
			Error: dto.ErrorObject{
				Code:          "invalid_data",
				MessageObject: dto.MessageObject{Message: "Please provide valid otp data"},
				Fields:        fieldErrors,
			},
		},
	}

	return
}

func checkValidTimezone(timezone string) bool {
	if timezone == "UTC" {
		return true
	}
	// if timezone != "UTC"{
	// 	return false
	// }
	return false
}
