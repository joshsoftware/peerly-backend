package organizationConfig

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	// logger "github.com/sirupsen/logrus"
)

type service struct {
	OranizationRepo     repository.OrganizationStorer
}

type Service interface {
	GetOrganizationConfig(ctx context.Context) (dto.OrganizationConfig, error)
	CreateOrganizationConfig(ctx context.Context, organization dto.OrganizationConfig) (dto.OrganizationConfig, error)
	UpdateOrganizationConfig(ctx context.Context, organization dto.OrganizationConfig) (dto.OrganizationConfig, error)
}

func NewService(oranizationRepo repository.OrganizationStorer) Service {
	return &service{
		OranizationRepo:     oranizationRepo,
	}
}


func (orgSvc *service) GetOrganizationConfig(ctx context.Context) (dto.OrganizationConfig, error) {

	organization, err := orgSvc.OranizationRepo.GetOrganizationConfig(ctx)
	if err != nil {
		return dto.OrganizationConfig{}, err
	}
	org := OrganizationConfigToDTO(organization)
	return org, nil

}


func (orgSvc *service) CreateOrganizationConfig(ctx context.Context, organizationConfig dto.OrganizationConfig) (dto.OrganizationConfig, error) {

	userID := ctx.Value("userId")

	var userIDInt64 int64
	switch v := userID.(type) {
	case int:
    	userIDInt64 = int64(v)
	default:
    return dto.OrganizationConfig{}, apperrors.UserNotFound
	}
	organizationConfig.CreatedBy = userIDInt64
	organizationConfig.UpdatedBy = userIDInt64

	_ ,err := orgSvc.OranizationRepo.GetOrganizationConfig(ctx);
	if err != apperrors.OrganizationNotFound {
		return dto.OrganizationConfig{},apperrors.OrganizationConfigAlreadyPresent
	}
	createdOrganization, err := orgSvc.OranizationRepo.CreateOrganizationConfig(ctx, organizationConfig)
	if err != nil {
		return dto.OrganizationConfig{}, err
	}
	org := OrganizationConfigToDTO(createdOrganization)
	return org, nil
}

func (orgSvc *service) UpdateOrganizationConfig(ctx context.Context, organizationConfig dto.OrganizationConfig) (dto.OrganizationConfig, error) {
	
	userID := ctx.Value("userId")

	var userIDInt64 int64
	switch v := userID.(type) {
	case int:
    	userIDInt64 = int64(v)
	default:
    return dto.OrganizationConfig{}, apperrors.UserNotFound
	}
	organizationConfig.UpdatedBy = userIDInt64

	_ ,err := orgSvc.OranizationRepo.GetOrganizationConfig(ctx);
	if err != nil {
		return dto.OrganizationConfig{},err
	}

	updatedOrganization, err := orgSvc.OranizationRepo.UpdateOrganizationCofig(ctx, organizationConfig)
	if err != nil {
		return dto.OrganizationConfig{}, err
	}
	org := OrganizationConfigToDTO(updatedOrganization)
	return org, nil
}
