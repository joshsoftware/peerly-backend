package repository

import (
	"context"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type OrganizationConfigStorer interface {
	GetOrganizationConfig(ctx context.Context, tx Transaction) (organization OrganizationConfig, err error)
	UpdateOrganizationConfig(ctx context.Context, tx Transaction, reqOrganization dto.OrganizationConfig) (updatedOrganization OrganizationConfig, err error)
	CreateOrganizationConfig(ctx context.Context, tx Transaction, org dto.OrganizationConfig) (createdOrganization OrganizationConfig, err error)
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
