package organizationConfig

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	OrganizationRepo     repository.OrganizationStorer
}

type Service interface {
	GetOrganizationConfig(ctx context.Context) (dto.OrganizationConfig, error)
	CreateOrganizationConfig(ctx context.Context, organization dto.OrganizationConfig) (dto.OrganizationConfig, error)
	UpdateOrganizationConfig(ctx context.Context, organization dto.OrganizationConfig) (dto.OrganizationConfig, error)
}

func NewService(organizationRepo repository.OrganizationStorer) Service {
	return &service{
		OrganizationRepo:     organizationRepo,
	}
}


func (orgSvc *service) GetOrganizationConfig(ctx context.Context) (dto.OrganizationConfig, error) {

	organization, err := orgSvc.OrganizationRepo.GetOrganizationConfig(ctx,nil)
	if err != nil {
		return dto.OrganizationConfig{}, err
	}
	org := OrganizationConfigToDTO(organization)
	return org, nil

}


func (orgSvc *service) CreateOrganizationConfig(ctx context.Context, organizationConfig dto.OrganizationConfig) (dto.OrganizationConfig, error) {

	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		logger.Error("err in parsing userid from token")
		return dto.OrganizationConfig{},apperrors.InternalServer
	}
	organizationConfig.CreatedBy = userID
	organizationConfig.UpdatedBy = userID

	_ ,err := orgSvc.OrganizationRepo.GetOrganizationConfig(ctx,nil);
	if err != apperrors.OrganizationConfigNotFound {
		return dto.OrganizationConfig{},apperrors.OrganizationConfigAlreadyPresent
	}
	createdOrganization, err := orgSvc.OrganizationRepo.CreateOrganizationConfig(ctx,nil, organizationConfig)
	if err != nil {
		return dto.OrganizationConfig{}, err
	}
	org := OrganizationConfigToDTO(createdOrganization)
	return org, nil
}

func (orgSvc *service) UpdateOrganizationConfig(ctx context.Context, organizationConfig dto.OrganizationConfig) (dto.OrganizationConfig, error) {
	
	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		logger.Error("err in parsing userid from token")
		return dto.OrganizationConfig{},apperrors.InternalServer
	}
	organizationConfig.UpdatedBy = userID

	_ ,err := orgSvc.OrganizationRepo.GetOrganizationConfig(ctx,nil);
	if err != nil {
		return dto.OrganizationConfig{},err
	}

	updatedOrganization, err := orgSvc.OrganizationRepo.UpdateOrganizationConfig(ctx,nil, organizationConfig)
	if err != nil {
		return dto.OrganizationConfig{}, err
	}
	org := OrganizationConfigToDTO(updatedOrganization)
	return org, nil
}
