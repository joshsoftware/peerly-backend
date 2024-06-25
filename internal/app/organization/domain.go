package organization

import (
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

func OrganizationDBToOrganization(orgDB repository.Organization) dto.Organization {
	return dto.Organization{
		ID:                       orgDB.ID,
		Name:                     orgDB.Name,
		ContactEmail:             orgDB.ContactEmail,
		DomainName:               orgDB.DomainName,
		SubscriptionStatus:       orgDB.SubscriptionStatus,
		SubscriptionValidUpto:    orgDB.SubscriptionValidUpto,
		Hi5Limit:                 orgDB.Hi5Limit,
		Hi5QuotaRenewalFrequency: orgDB.Hi5QuotaRenewalFrequency,
		Timezone:                 orgDB.Timezone,
		CreatedAt:                orgDB.CreatedAt,
		CreatedBy:                orgDB.CreatedBy,
		UpdatedAt:                orgDB.UpdatedAt,
	}
}
