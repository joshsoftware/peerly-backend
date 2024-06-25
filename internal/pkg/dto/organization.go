package dto

import (
	"time"
)

type Organization struct {
	ID                       int64     `json:"id" `
	Name                     string    `json:"name"`
	ContactEmail             string    `json:"email"`
	DomainName               string    `json:"domain_name"`
	SubscriptionStatus       int       `json:"subscription_status"`
	SubscriptionValidUpto    time.Time `json:"subscription_valid_upto"`
	Hi5Limit                 int       `json:"hi5_limit"`
	Hi5QuotaRenewalFrequency string    `json:"hi5_quota_renewal_frequency"`
	Timezone                 string    `json:"timezone"`
	CreatedAt                time.Time `json:"created_at"`
	CreatedBy                int64     `json:"created_by"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type OTP struct {
	CreatedAt   time.Time
	OrgId      int64 `json:"org_id"`
	OTPCode     string `json:"otpcode"`
}