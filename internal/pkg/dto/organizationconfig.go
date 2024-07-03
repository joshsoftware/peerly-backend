package dto

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
