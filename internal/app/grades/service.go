package grades

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	gradesRepo repository.GradesStorer
}

type Service interface {
	ListGrades(ctx context.Context) (resp []dto.Grade, err error)
}

func NewService(gradesRepo repository.GradesStorer) Service {
	return &service{
		gradesRepo: gradesRepo,
	}
}

func (gs *service) ListGrades(ctx context.Context) (resp []dto.Grade, err error) {

	dbResp, err := gs.gradesRepo.ListGrades(ctx)
	if err != nil {
		logger.Error(err.Error())
		err = apperrors.InternalServerError
	}

	for _, item := range dbResp {
		svcItem := mapDbToSvc(item)
		resp = append(resp, svcItem)
	}

	return

}

func mapDbToSvc(dbResp repository.Grade) (svcResp dto.Grade) {
	svcResp.Id = dbResp.Id
	svcResp.Name = dbResp.Name
	svcResp.Points = dbResp.Points
	return
}
