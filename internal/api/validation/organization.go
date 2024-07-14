package validation

import (

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
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

func OrgValidate(org dto.OrganizationConfig) (err error)  {

	if org.RewardMultiplier <= 0 {
		return apperrors.InvalidRewardMultiplier
	}

	if org.RewardQuotaRenewalFrequency <= 0  {
		return apperrors.InvalidRewardQuotaRenewalFrequency
	}

	if !isTimeZoneValid(org.Timezone) {
		return apperrors.InvalidTimezone
	}

	return
}

func OrgUpdateValidate(org dto.OrganizationConfig) (err error)  {

	if org.Timezone != "" {
		if !isTimeZoneValid(org.Timezone) {
			return apperrors.InvalidTimezone
		}
	}

	return
}

func isTimeZoneValid(tz string) bool {
	_, exists := timeZones[tz]
	return exists
}