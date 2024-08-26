package corevalues

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	coreValuesRepo repository.CoreValueStorer
}

type Service interface {
	ListCoreValues(ctx context.Context) (resp []dto.CoreValue, err error)
	GetCoreValue(ctx context.Context, coreValueID string) (coreValue dto.CoreValue, err error)
	CreateCoreValue(ctx context.Context, coreValue dto.CreateCoreValueReq) (resp dto.CoreValue, err error)
	UpdateCoreValue(ctx context.Context, coreValueID string, coreValue dto.UpdateQueryRequest) (resp dto.CoreValue, err error)
}

func NewService(coreValuesRepo repository.CoreValueStorer) Service {
	return &service{
		coreValuesRepo: coreValuesRepo,
	}
}

func (cs *service) ListCoreValues(ctx context.Context) (resp []dto.CoreValue, err error) {

	dbResp, err := cs.coreValuesRepo.ListCoreValues(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}

	for _, value := range dbResp {
		coreValue := mapCoreValueDbToService(value)
		resp = append(resp, coreValue)
	}

	return

}

func (cs *service) GetCoreValue(ctx context.Context, coreValueID string) (coreValue dto.CoreValue, err error) {

	coreValId, err := utils.VarsStringToInt(coreValueID, "coreValueId")
	if err != nil {
		return
	}

	dbResp, err := cs.coreValuesRepo.GetCoreValue(ctx, coreValId)
	if err != nil {
		return
	}

	coreValue = mapCoreValueDbToService(dbResp)

	return
}

func (cs *service) CreateCoreValue(ctx context.Context, coreValue dto.CreateCoreValueReq) (resp dto.CoreValue, err error) {

	isUnique, err := cs.coreValuesRepo.CheckUniqueCoreVal(ctx, coreValue.Name)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}
	if !isUnique {
		err = apperrors.UniqueCoreValue
		return
	}

	err = cs.validate(ctx, coreValue)
	if err != nil {
		return
	}

	dbResp, err := cs.coreValuesRepo.CreateCoreValue(ctx, coreValue)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}

	resp = mapCoreValueDbToService(dbResp)

	return
}

func (cs *service) UpdateCoreValue(ctx context.Context, coreValueID string, reqData dto.UpdateQueryRequest) (resp dto.CoreValue, err error) {

	coreValId, err := utils.VarsStringToInt(coreValueID, "coreValueId")
	if err != nil {
		return
	}

	//validate corevalue
	//get data
	coreValue, err := cs.coreValuesRepo.GetCoreValue(ctx, coreValId)
	if err != nil {
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
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}
	if !isUnique && reqData.Name != coreValue.Name {
		err = apperrors.UniqueCoreValue
		return
	}

	reqData.Id = coreValId

	dbResp, err := cs.coreValuesRepo.UpdateCoreValue(ctx, reqData)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError

		return
	}

	resp = mapCoreValueDbToService(dbResp)

	return
}

func (cs *service) validateParentCoreValue(ctx context.Context, coreValueID int64) (ok bool) {
	coreValue, err := cs.coreValuesRepo.GetCoreValue(ctx, coreValueID)
	if err != nil {
		logger.Errorf(ctx, "parent core value id not present, err: %s", err.Error())
		return
	}

	if coreValue.ParentCoreValueID.Valid {
		logger.Error(ctx, "Invalid parent core value id")
		return
	}

	return true
}

func (cs *service) validate(ctx context.Context, coreValue dto.CreateCoreValueReq) (err error) {

	if coreValue.Name == "" {
		err = apperrors.TextFieldBlank
	}
	if coreValue.Description == "" {
		err = apperrors.DescFieldBlank
	}
	if coreValue.ParentCoreValueID != nil {
		if !cs.validateParentCoreValue(ctx, *coreValue.ParentCoreValueID) {
			err = apperrors.InvalidParentValue
		}
	}

	return
}

func mapCoreValueDbToService(dbStruct repository.CoreValue) (svcStruct dto.CoreValue) {
	svcStruct.ID = dbStruct.ID
	svcStruct.Name = dbStruct.Name
	svcStruct.Description = dbStruct.Description
	svcStruct.ParentCoreValueID = dbStruct.ParentCoreValueID.Int64
	return
}
