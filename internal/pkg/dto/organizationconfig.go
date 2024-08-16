package dto
import(
	"time"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
)


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

func(orgConfig OrganizationConfig) OrgValidate() (err error)  {

	if orgConfig.RewardMultiplier <= 0 {
		return apperrors.InvalidRewardMultiplier
	}

	if orgConfig.RewardQuotaRenewalFrequency <= 0  {
		return apperrors.InvalidRewardQuotaRenewalFrequency
	}

	if !isValidTimeZone(orgConfig.Timezone) {
		return apperrors.InvalidTimezone
	}

	return
}

func (orgConfig OrganizationConfig) OrgUpdateValidate() (err error)  {

	if orgConfig.Timezone != "" {
		if !isValidTimeZone(orgConfig.Timezone) {
			return apperrors.InvalidTimezone
		}
	}

	return
}

// Check if a given time zone is valid
func isValidTimeZone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}