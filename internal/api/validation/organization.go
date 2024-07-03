package validation

import (

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

// map with all the time zones
var timeZones = map[string]bool{
	"NZDT":  true, "IDLE":  true, "NZST":  true, "NZT":   true,
	"AESST": true, "ACSST": true, "CADT":  true, "SADT":  true,
	"AEST":  true, "EAST":  true, "GST":   true, "LIGT":  true,
	"SAST":  true, "CAST":  true, "AWSST": true, "JST":   true,
	"KST":   true, "MHT":   true, "WDT":   true, "MT":    true,
	"AWST":  true, "CCT":   true, "WADT":  true, "WST":   true,
	"JT":    true, "ALMST": true, "WAST":  true, "CXT":   true,
	"ALMT":  true, "MAWT":  true, "IOT":   true, "MVT":   true,
	"TFT":   true, "AFT":   true, "MUT":   true,
	"RET":   true, "SCT":   true, "IT":    true, "EAT":   true,
	"BT":    true, "EETDST":true, "HMT":   true, "BDST":  true,
	"CEST":  true, "CETDST":true, "EET":   true, "FWT":   true,
	"IST":   true, "MEST":  true, "METDST":true, "SST":   true,
	"BST":   true, "CET":   true, "DNT":   true, "FST":   true,
	"MET":   true, "MEWT":  true, "MEZ":   true, "NOR":   true,
	"SET":   true, "SWT":   true, "WETDST":true, "GMT":   true,
	"UT":    true, "UTC":   true, "Z":     true, "ZULU":  true,
	"WET":   true, "WAT":   true, "NDT":   true, "ADT":   true,
	"AWT":   true, "NFT":   true, "NST":   true, "AST":   true,
	"ACST":  true, "ACT":   true, "EDT":   true, "CDT":   true,
	"EST":   true, "CST":   true, "MDT":   true, "MST":   true,
	"PDT":   true, "AKDT":  true, "PST":   true, "YDT":   true,
	"AKST":  true, "HDT":   true, "YST":   true, "AHST":  true,
	"HST":   true, "CAT":   true, "NT":    true, "IDLW":  true,
}

func OrgValidate(org dto.OrganizationConfig) (errorResponse map[string]dto.ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if org.RewardMultiplier <= 0 {
		fieldErrors["reward_multiplier"] = "Please enter reward multiplier greater than 0"
	}

	if org.RewardQuotaRenewalFrequency <= 0  {
		fieldErrors["reward_quota_renewal_frequency"] = "Please enter valid reward renewal frequency"
	}

	if !isTimeZoneValid(org.Timezone) {
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

func OrgUpdateValidate(org dto.OrganizationConfig) (errorResponse map[string]dto.ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if org.Timezone != "" {
		if !isTimeZoneValid(org.Timezone) {
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

func isTimeZoneValid(tz string) bool {
	_, exists := timeZones[tz]
	return exists
}