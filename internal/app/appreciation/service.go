package appreciation

import (
	"context"
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	appreciationRepo repository.AppreciationStorer
	corevaluesRespo  repository.CoreValueStorer
}

type Service interface {
	CreateAppreciation(ctx context.Context, apprecication dto.Appreciation) (dto.Appreciation, error)
	GetAppreciationById(ctx context.Context, appreciationId int) (dto.ResponseAppreciation, error)
	GetAppreciation(ctx context.Context, filter dto.AppreciationFilter) ([]dto.ResponseAppreciation, error)
	ValidateAppreciation(ctx context.Context, isValid bool, apprId int) (bool, error)
}

func NewService(appreciationRepo repository.AppreciationStorer, coreValuesRepo repository.CoreValueStorer) Service {
	return &service{
		appreciationRepo: appreciationRepo,
		corevaluesRespo:  coreValuesRepo,
	}
}

func (apprSvc *service) CreateAppreciation(ctx context.Context, apprecication dto.Appreciation) (dto.Appreciation, error) {

	//add quarter
	apprecication.Quarter = GetQuarter()

	//add sender
	apprecication.Sender = ctx.Value("userId").(int64)

	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)
	fmt.Println("ERr", err, tx)
	
	if err != nil {
		return dto.Appreciation{}, err
	}

	defer func() {
		txErr := apprSvc.appreciationRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()


	//check is corevalue present in database
	_, err = apprSvc.corevaluesRespo.GetCoreValue(ctx, int64(apprecication.CoreValueID))
	if err != nil {
		return dto.Appreciation{}, err
	}

	//check is receiver present in database
	chk, err := apprSvc.appreciationRepo.IsUserPresent(ctx, tx, apprecication.Receiver)
	if err != nil {
		return dto.Appreciation{}, err
	}
	if !chk {
		return dto.Appreciation{}, apperrors.UserNotFound
	}

	appr, err := apprSvc.appreciationRepo.CreateAppreciation(ctx, tx, apprecication)
	if err != nil {
		return dto.Appreciation{}, err
	}

	return MapAppreciationDBToDTO(appr), nil
}

func (apprSvc *service) GetAppreciationById(ctx context.Context, appreciationId int) (dto.ResponseAppreciation, error) {

	resAppr, err := apprSvc.appreciationRepo.GetAppreciationById(ctx, nil, appreciationId)
	if err != nil {
		return dto.ResponseAppreciation{}, err
	}

	return mapAppreciationInfoToResponse(resAppr), nil
}

func (apprSvc *service) GetAppreciation(ctx context.Context, filter dto.AppreciationFilter) ([]dto.ResponseAppreciation, error) {
	infos, err := apprSvc.appreciationRepo.GetAppreciation(ctx, nil, filter)
	if err != nil {
		return []dto.ResponseAppreciation{}, err
	}

	var responses []dto.ResponseAppreciation
	for _, info := range infos {
		response := mapAppreciationInfoToResponse(info)
		responses = append(responses, response)
	}

	return responses, nil
}

func (apprSvc *service) ValidateAppreciation(ctx context.Context, isValid bool, apprId int) (bool, error) {
	return apprSvc.appreciationRepo.ValidateAppreciation(ctx, nil, isValid, apprId)
}
