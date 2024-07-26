package appreciation

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	appreciationRepo repository.AppreciationStorer
	corevaluesRespo  repository.CoreValueStorer
}

// Service contains all
type Service interface {
	CreateAppreciation(ctx context.Context, appreciation dto.Appreciation) (dto.Appreciation, error)
	GetAppreciationById(ctx context.Context, appreciationId int32) (dto.AppreciationResponse, error)
	ListAppreciations(ctx context.Context, filter dto.AppreciationFilter) (dto.ListAppreciationsResponse, error)
	DeleteAppreciation(ctx context.Context, apprId int32) error
}

func NewService(appreciationRepo repository.AppreciationStorer, coreValuesRepo repository.CoreValueStorer) Service {
	return &service{
		appreciationRepo: appreciationRepo,
		corevaluesRespo:  coreValuesRepo,
	}
}

func (apprSvc *service) CreateAppreciation(ctx context.Context, appreciation dto.Appreciation) (dto.Appreciation, error) {

	//add quarter
	appreciation.Quarter = utils.GetQuarter()

	//add sender
	data := ctx.Value(constants.UserId)
	sender, ok := data.(int64)
	if !ok {
		logger.Error("err in parsing userid from token")
		return dto.Appreciation{}, apperrors.InternalServer
	}

	//check is receiver present in database
	chk, err := apprSvc.appreciationRepo.IsUserPresent(ctx, nil, appreciation.Receiver)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.Appreciation{}, err
	}
	if !chk {
		return dto.Appreciation{}, apperrors.UserNotFound
	}
	appreciation.Sender = sender

	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)

	if err != nil {
		return dto.Appreciation{}, err
	}

	defer func() {
		rvr := recover()
		defer func() {
			if rvr != nil {
				logger.Info(ctx, "Transaction aborted because of panic: %v, Propagating panic further", rvr)
				panic(rvr)
			}
		}()

		txErr := apprSvc.appreciationRepo.HandleTransaction(ctx, tx, err == nil && rvr == nil)
		if txErr != nil {
			err = txErr
			logger.Info(ctx, "error in creating transaction, err: %s", txErr.Error())
			return
		}
	}()

	//check is corevalue present in database
	_, err = apprSvc.corevaluesRespo.GetCoreValue(ctx, int64(appreciation.CoreValueID))
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.Appreciation{}, err
	}

	// check self appreciation
	if appreciation.Receiver == sender {
		return dto.Appreciation{}, apperrors.SelfAppreciationError
	}

	appr, err := apprSvc.appreciationRepo.CreateAppreciation(ctx, tx, appreciation)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.Appreciation{}, err
	}

	return mapAppreciationDBToDTO(appr), nil
}

func (apprSvc *service) GetAppreciationById(ctx context.Context, appreciationId int32) (dto.AppreciationResponse, error) {

	resAppr, err := apprSvc.appreciationRepo.GetAppreciationById(ctx, nil, appreciationId)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.AppreciationResponse{}, err
	}

	return mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(resAppr), nil
}

func (apprSvc *service) ListAppreciations(ctx context.Context, filter dto.AppreciationFilter) (dto.ListAppreciationsResponse, error) {

	infos, pagination, err := apprSvc.appreciationRepo.ListAppreciations(ctx, nil, filter)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.ListAppreciationsResponse{}, err
	}

	responses := make([]dto.AppreciationResponse, 0)
	for _, info := range infos {
		responses = append(responses, mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(info))
	}
	paginationResp := dtoPagination(pagination)
	return dto.ListAppreciationsResponse{Appreciations: responses, MetaData: paginationResp}, nil
}

func (apprSvc *service) DeleteAppreciation(ctx context.Context, apprId int32) error {
	return apprSvc.appreciationRepo.DeleteAppreciation(ctx, nil, apprId)
}
