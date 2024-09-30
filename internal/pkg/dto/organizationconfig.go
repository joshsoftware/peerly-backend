package dto

import (
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
)

// map with all the time zones
var timeZones = map[string]bool{
	"NZDT": true, "IDLE": true, "NZST": true, "NZT": true,
	"AESST": true, "ACSST": true, "CADT": true, "SADT": true,
	"AEST": true, "EAST": true, "GST": true, "LIGT": true,
	"SAST": true, "CAST": true, "AWSST": true, "JST": true,
	"KST": true, "MHT": true, "WDT": true, "MT": true,
	"AWST": true, "CCT": true, "WADT": true, "WST": true,
	"JT": true, "ALMST": true, "WAST": true, "CXT": true,
	"ALMT": true, "MAWT": true, "IOT": true, "MVT": true,
	"TFT": true, "AFT": true, "MUT": true,
	"RET": true, "SCT": true, "IT": true, "EAT": true,
	"BT": true, "EETDST": true, "HMT": true, "BDST": true,
	"CEST": true, "CETDST": true, "EET": true, "FWT": true,
	"IST": true, "MEST": true, "METDST": true, "SST": true,
	"BST": true, "CET": true, "DNT": true, "FST": true,
	"MET": true, "MEWT": true, "MEZ": true, "NOR": true,
	"SET": true, "SWT": true, "WETDST": true, "GMT": true,
	"UT": true, "UTC": true, "Z": true, "ZULU": true,
	"WET": true, "WAT": true, "NDT": true, "ADT": true,
	"AWT": true, "NFT": true, "NST": true, "AST": true,
	"ACST": true, "ACT": true, "EDT": true, "CDT": true,
	"EST": true, "CST": true, "MDT": true, "MST": true,
	"PDT": true, "AKDT": true, "PST": true, "YDT": true,
	"AKST": true, "HDT": true, "YST": true, "AHST": true,
	"HST": true, "CAT": true, "NT": true, "IDLW": true,
}

type OrganizationConfig struct {
	ID                          int64  `json:"id"`
	RewardMultiplier            int    `json:"reward_multiplier"`
	RewardQuotaRenewalFrequency int    `json:"reward_quota_renewal_frequency"`
	Timezone                    string `json:"timezone"`
	CreatedAt                   int64  `json:"created_at"`
	CreatedBy                   int64  `json:"created_by"`
	UpdatedAt                   int64  `json:"updated_at"`
	UpdatedBy                   int64  `json:"updated_by"`
}

func (orgConfig OrganizationConfig) OrgValidate() (err error) {

	if orgConfig.RewardMultiplier <= 0 {
		return apperrors.InvalidRewardMultiplier
	}

	if orgConfig.RewardQuotaRenewalFrequency <= 0 {
		return apperrors.InvalidRewardQuotaRenewalFrequency
	}

	if !isTimeZoneValid(orgConfig.Timezone) {
		return apperrors.InvalidTimezone
	}

	return
}

func (orgConfig OrganizationConfig) OrgUpdateValidate() (err error) {

	if orgConfig.Timezone != "" {
		if !isTimeZoneValid(orgConfig.Timezone) {
			return apperrors.InvalidTimezone
		}
	}

	return
}

func isTimeZoneValid(tz string) bool {
	_, exists := timeZones[tz]
	return exists
}
