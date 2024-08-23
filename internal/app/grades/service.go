package grades

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	gradesRepo repository.GradesStorer
}

type Service interface {
	ListGrades(ctx context.Context) (resp []dto.Grade, err error)
	EditGrade(ctx context.Context, id string, points int64) (err error)
}

func NewService(gradesRepo repository.GradesStorer) Service {
	return &service{
		gradesRepo: gradesRepo,
	}
}

func (gs *service) ListGrades(ctx context.Context) (resp []dto.Grade, err error) {

	dbResp, err := gs.gradesRepo.ListGrades(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
	}

	for _, item := range dbResp {
		svcItem := mapDbToSvc(item)
		resp = append(resp, svcItem)
	}

	return

}

func (gs *service) EditGrade(ctx context.Context, id string, points int64) (err error) {
	gradeId, err := utils.VarsStringToInt(id, "gradeId")
	if err != nil {
		return
	}
	var reqData dto.UpdateGradeReq
	reqData.Id = gradeId
	if points < 0 {
		logger.Errorf(ctx, "grade points cannot be negative, grade points: %d", points)
		err = apperrors.NegativeGradePoints
		return
	}
	reqData.Points = points
	err = gs.gradesRepo.EditGrade(ctx, reqData)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}
	return
}

func mapDbToSvc(dbResp repository.Grade) (svcResp dto.Grade) {
	svcResp.Id = dbResp.Id
	svcResp.Name = dbResp.Name
	svcResp.Points = dbResp.Points
	return
}
