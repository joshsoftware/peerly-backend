package repository

import (
	"context"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type OrganizationStorer interface {
	ListOrganizations(ctx context.Context) (organizations []Organization, err error)
	GetOrganization(ctx context.Context, organizationID int) (organization Organization, err error)
	GetOrganizationByDomainName(ctx context.Context, domainName string) (organization Organization, err error)
	DeleteOrganization(ctx context.Context, organizationID int, userId int64) (err error)
	UpdateOrganization(ctx context.Context, reqOrganization dto.Organization) (updatedOrganization Organization, err error)
	CreateOrganization(ctx context.Context, org dto.Organization) (createdOrganization Organization, err error)

	IsEmailPresent(ctx context.Context, email string) bool
	IsDomainPresent(ctx context.Context, domainName string) bool
	IsOrganizationIdPresent(ctx context.Context, organizationId int64) bool
}

type Organization struct {
	ID                       int64     `db:"id"`
	Name                     string    `db:"name"`
	ContactEmail             string    `db:"contact_email"`
	DomainName               string    `db:"domain_name"`
	SubscriptionStatus       int       `db:"subscription_status"`
	SubscriptionValidUpto    time.Time `db:"subscription_valid_upto"`
	Hi5Limit                 int       `db:"hi5_limit"`
	Hi5QuotaRenewalFrequency string    `db:"hi5_quota_renewal_frequency"`
	Timezone                 string    `db:"timezone"`
	CreatedAt                time.Time `db:"created_at"`
	CreatedBy                int64     `db:"created_by"`
	UpdatedAt                time.Time `db:"updated_at"`
}