package organizationConfig

import (
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

func organizationConfigToDTO(org repository.OrganizationConfig) dto.OrganizationConfig {
	return dto.OrganizationConfig{
		ID:                          org.ID,
		RewardMultiplier:            org.RewardMultiplier,
		RewardQuotaRenewalFrequency: org.RewardQuotaRenewalFrequency,
		Timezone:                    org.Timezone,
		CreatedAt:                   org.CreatedAt,
		CreatedBy:                   org.CreatedBy,
		UpdatedAt:                   org.UpdatedAt,
		UpdatedBy:                   org.UpdatedBy,
	}
}
