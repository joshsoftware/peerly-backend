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
	OrganizationConfigRepo     repository.OrganizationConfigStorer
}

type Service interface {
	GetOrganizationConfig(ctx context.Context) (dto.OrganizationConfig, error)
	CreateOrganizationConfig(ctx context.Context, organizationConfigInfo dto.OrganizationConfig) (dto.OrganizationConfig, error)
	UpdateOrganizationConfig(ctx context.Context, organizationConfigInfo dto.OrganizationConfig) (dto.OrganizationConfig, error)
}

func NewService(organizationConfigRepo repository.OrganizationConfigStorer) Service {
	return &service{
		OrganizationConfigRepo:     organizationConfigRepo,
	}
}


func (orgSvc *service) GetOrganizationConfig(ctx context.Context) (dto.OrganizationConfig, error) {

	organization, err := orgSvc.OrganizationConfigRepo.GetOrganizationConfig(ctx,nil)
	if err != nil {
		logger.Errorf("err: %v",err)
		return dto.OrganizationConfig{}, err
	}
	org := organizationConfigToDTO(organization)
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

	_ ,err := orgSvc.OrganizationConfigRepo.GetOrganizationConfig(ctx,nil);
	if err != apperrors.OrganizationConfigNotFound {
		return dto.OrganizationConfig{},apperrors.OrganizationConfigAlreadyPresent
	}

	createdOrganizationConfig, err := orgSvc.OrganizationConfigRepo.CreateOrganizationConfig(ctx,nil, organizationConfig)
	if err != nil {
		logger.Errorf("err: %v",err)
		return dto.OrganizationConfig{}, err
	}
	
	org := organizationConfigToDTO(createdOrganizationConfig)
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

	_ ,err := orgSvc.OrganizationConfigRepo.GetOrganizationConfig(ctx,nil);
	if err != nil {
		logger.Errorf("err: %v",err)
		return dto.OrganizationConfig{},err
	}

	updatedOrganization, err := orgSvc.OrganizationConfigRepo.UpdateOrganizationConfig(ctx,nil, organizationConfig)
	if err != nil {
		logger.Errorf("err: %v",err)
		return dto.OrganizationConfig{}, err
	}
	org := organizationConfigToDTO(updatedOrganization)
	return org, nil
}
