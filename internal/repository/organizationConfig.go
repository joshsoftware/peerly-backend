package repository

import (
	"context"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type OrganizationStorer interface {
	GetOrganizationConfig(ctx context.Context) (organization OrganizationConfig, err error)
	UpdateOrganizationCofig(ctx context.Context, reqOrganization dto.OrganizationConfig) (updatedOrganization OrganizationConfig, err error)
	CreateOrganizationConfig(ctx context.Context, org dto.OrganizationConfig) (createdOrganization OrganizationConfig, err error)
}

type OrganizationConfig struct {
	ID                          int64  `db:"id"`
	RewardMultiplier            int    `db:"reward_multiplier"`
	RewardQuotaRenewalFrequency int    `db:"reward_quota_renewal_frequency"`
	Timezone                    string `db:"timezone"`
	CreatedAt                   int64  `db:"created_at"`
	CreatedBy                   int64  `db:"created_by"`
	UpdatedAt                   int64  `db:"updated_at"`
	UpdatedBy                   int64  `db:"updated_by"`
}
