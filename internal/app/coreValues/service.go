package corevalues

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	coreValuesRepo repository.CoreValueStorer
}

type Service interface {
	ListCoreValues(ctx context.Context, organisationID string) (resp []dto.ListCoreValuesResp, err error)
	GetCoreValue(ctx context.Context, organisationID string, coreValueID string) (coreValue dto.GetCoreValueResp, err error)
	CreateCoreValue(ctx context.Context, organisationID string, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error)
	DeleteCoreValue(ctx context.Context, organisationID string, coreValueID string, userId int64) (err error)
	UpdateCoreValue(ctx context.Context, organisationID string, coreValueID string, coreValue dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error)
}

func NewService(coreValuesRepo repository.CoreValueStorer) Service {
	return &service{
		coreValuesRepo: coreValuesRepo,
	}
}

func (cs *service) ListCoreValues(ctx context.Context, organisationID string) (resp []dto.ListCoreValuesResp, err error) {

	if organisationID == "" {
		err = apperrors.InvalidOrgId
		return
	}

	orgId, err := VarsStringToInt(organisationID, "organisationId")
	if err != nil {
		return
	}

	err = cs.coreValuesRepo.CheckOrganisation(ctx, orgId)
	if err != nil {
		return
	}

	resp, err = cs.coreValuesRepo.ListCoreValues(ctx, orgId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching data")
		err = apperrors.InternalServerError
	}

	return

}

func (cs *service) GetCoreValue(ctx context.Context, organisationID string, coreValueID string) (coreValue dto.GetCoreValueResp, err error) {

	orgId, err := VarsStringToInt(organisationID, "organisationId")
	if err != nil {
		return
	}

	err = cs.coreValuesRepo.CheckOrganisation(ctx, orgId)
	if err != nil {
		return
	}

	coreValId, err := VarsStringToInt(coreValueID, "coreValueId")
	if err != err {
		return
	}

	coreValue, err = cs.coreValuesRepo.GetCoreValue(ctx, orgId, coreValId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching data")
		return
	}

	return
}

func (cs *service) CreateCoreValue(ctx context.Context, organisationID string, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error) {

	orgId, err := VarsStringToInt(organisationID, "organisationId")
	if err != nil {
		return
	}

	err = cs.coreValuesRepo.CheckOrganisation(ctx, orgId)
	if err != nil {
		return
	}

	isUnique, err := cs.coreValuesRepo.CheckUniqueCoreVal(ctx, orgId, coreValue.Text)
	if err != nil {
		return
	}
	if !isUnique {
		err = apperrors.UniqueCoreValue
		return
	}

	err = Validate(ctx, coreValue, cs.coreValuesRepo, orgId)
	if err != nil {
		return
	}

	resp, err = cs.coreValuesRepo.CreateCoreValue(ctx, orgId, userId, coreValue)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while creating core value")
		err = apperrors.InternalServerError
		return
	}

	return
}

func (cs *service) DeleteCoreValue(ctx context.Context, organisationID string, coreValueID string, userId int64) (err error) {

	orgId, err := VarsStringToInt(organisationID, "organisationId")
	if err != nil {
		return
	}

	err = cs.coreValuesRepo.CheckOrganisation(ctx, orgId)
	if err != nil {
		return
	}

	coreValId, err := VarsStringToInt(coreValueID, "coreValueId")
	if err != nil {
		return
	}

	coreValue, err := cs.coreValuesRepo.GetCoreValue(ctx, orgId, coreValId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching data")
		return
	}

	if coreValue.SoftDelete {
		err = apperrors.InvalidCoreValueData
		return
	}

	err = cs.coreValuesRepo.DeleteCoreValue(ctx, orgId, coreValId, userId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while deleting core value")
		err = apperrors.InternalServerError

		return
	}

	return
}

func (cs *service) UpdateCoreValue(ctx context.Context, organisationID string, coreValueID string, reqData dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error) {

	orgId, err := VarsStringToInt(organisationID, "organisationId")
	if err != nil {
		return
	}

	coreValId, err := VarsStringToInt(coreValueID, "coreValueId")
	if err != nil {
		return
	}

	//validate organisation
	err = cs.coreValuesRepo.CheckOrganisation(ctx, orgId)
	if err != nil {
		return
	}

	//validate corevalue
	//get data
	coreValue, err := cs.coreValuesRepo.GetCoreValue(ctx, orgId, coreValId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching data")
		return
	}

	if coreValue.SoftDelete {
		err = apperrors.InvalidCoreValueData
		return
	}

	//set empty fields
	if reqData.Text == "" {
		reqData.Text = coreValue.Text
	}
	if reqData.Description == "" {
		reqData.Description = coreValue.Description
	}
	if reqData.ThumbnailUrl == "" {
		reqData.ThumbnailUrl = *coreValue.ThumbnailURL
	}

	isUnique, err := cs.coreValuesRepo.CheckUniqueCoreVal(ctx, orgId, reqData.Text)
	if err != nil {
		return
	}
	if !isUnique && reqData.Text != coreValue.Text {
		err = apperrors.UniqueCoreValue
		return
	}

	resp, err = cs.coreValuesRepo.UpdateCoreValue(ctx, orgId, coreValId, reqData)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while updating core value")
		err = apperrors.InternalServerError

		return
	}

	return
}
