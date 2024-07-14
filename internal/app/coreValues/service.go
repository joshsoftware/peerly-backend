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
	ListCoreValues(ctx context.Context) (resp []dto.ListCoreValuesResp, err error)
	GetCoreValue(ctx context.Context, coreValueID string) (coreValue dto.GetCoreValueResp, err error)
	CreateCoreValue(ctx context.Context, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error)
	UpdateCoreValue(ctx context.Context, coreValueID string, coreValue dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error)
}

func NewService(coreValuesRepo repository.CoreValueStorer) Service {
	return &service{
		coreValuesRepo: coreValuesRepo,
	}
}

func (cs *service) ListCoreValues(ctx context.Context) (resp []dto.ListCoreValuesResp, err error) {

	resp, err = cs.coreValuesRepo.ListCoreValues(ctx)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching data")
		err = apperrors.InternalServerError
	}

	return

}

func (cs *service) GetCoreValue(ctx context.Context, coreValueID string) (coreValue dto.GetCoreValueResp, err error) {

	coreValId, err := VarsStringToInt(coreValueID, "coreValueId")
	if err != err {
		return
	}

	coreValue, err = cs.coreValuesRepo.GetCoreValue(ctx, coreValId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching data")
		return
	}

	return
}

func (cs *service) CreateCoreValue(ctx context.Context, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error) {

	isUnique, err := cs.coreValuesRepo.CheckUniqueCoreVal(ctx, coreValue.Name)
	if err != nil {
		return
	}
	if !isUnique {
		err = apperrors.UniqueCoreValue
		return
	}

	err = Validate(ctx, coreValue, cs.coreValuesRepo)
	if err != nil {
		return
	}

	resp, err = cs.coreValuesRepo.CreateCoreValue(ctx, userId, coreValue)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while creating core value")
		err = apperrors.InternalServerError
		return
	}

	return
}

func (cs *service) UpdateCoreValue(ctx context.Context, coreValueID string, reqData dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error) {

	coreValId, err := VarsStringToInt(coreValueID, "coreValueId")
	if err != nil {
		return
	}

	//validate corevalue
	//get data
	coreValue, err := cs.coreValuesRepo.GetCoreValue(ctx, coreValId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching data")
		return
	}

	//set empty fields
	if reqData.Name == "" {
		reqData.Name = coreValue.Name
	}
	if reqData.Description == "" {
		reqData.Description = coreValue.Description
	}

	isUnique, err := cs.coreValuesRepo.CheckUniqueCoreVal(ctx, reqData.Name)
	if err != nil {
		return
	}
	if !isUnique && reqData.Name != coreValue.Name {
		err = apperrors.UniqueCoreValue
		return
	}

	resp, err = cs.coreValuesRepo.UpdateCoreValue(ctx, coreValId, reqData)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while updating core value")
		err = apperrors.InternalServerError

		return
	}

	return
}
