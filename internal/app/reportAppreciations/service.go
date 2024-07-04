package reportappreciations

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	reportAppreciationRepo repository.ReportAppreciationStorer
}

type Service interface {
	ReportAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (resp dto.ReportAppricaitionResp, err error)
}

func NewService(reportAppreciationRepo repository.ReportAppreciationStorer) Service {
	return &service{
		reportAppreciationRepo: reportAppreciationRepo,
	}
}

func (rs *service) ReportAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (resp dto.ReportAppricaitionResp, err error) {

	reporterId := ctx.Value(constants.UserId)
	data, ok := reporterId.(int64)
	if !ok {
		logger.Error("Error in typecasting reporter id")
		err = apperrors.InternalServerError
		return
	}
	reqData.ReportedBy = data

	doesAppreciationExist, err := rs.reportAppreciationRepo.CheckAppreciation(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	if !doesAppreciationExist {
		err = apperrors.InvalidId
		return
	}

	isDupliate, err := rs.reportAppreciationRepo.CheckDuplicateReport(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	if isDupliate {
		err = apperrors.RepeatedReport
		return
	}

	usersData, err := rs.reportAppreciationRepo.GetSenderAndReceiver(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	if usersData.Sender == reqData.ReportedBy || usersData.Receiver == reqData.ReportedBy {
		err = apperrors.CannotReportOwnAppreciation
		return
	}

	resp, err = rs.reportAppreciationRepo.ReportAppreciation(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}

	return
}
