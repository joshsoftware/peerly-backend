package grades

import (
	"context"
	"fmt"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	gradesRepo repository.GradesStorer
	userRepo   repository.UserStorer
}

type Service interface {
	ListGrades(ctx context.Context) (resp []dto.Grade, err error)
	EditGrade(ctx context.Context, id string, points int64) (err error)
}

func NewService(gradesRepo repository.GradesStorer, userRepo repository.UserStorer) Service {
	return &service{
		gradesRepo: gradesRepo,
		userRepo:   userRepo,
	}
}

func (gs *service) ListGrades(ctx context.Context) (resp []dto.Grade, err error) {

	dbResp, err := gs.gradesRepo.ListGrades(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
	}

	for _, item := range dbResp {
		fmt.Println("updated by: ", item.UpdatedBy)
		if item.UpdatedBy.Valid {
			reqData := dto.GetUserByIdReq{
				UserId:          item.UpdatedBy.Int64,
				QuaterTimeStamp: GetQuarterStartUnixTime(),
			}
			user, err := gs.userRepo.GetUserById(ctx, reqData)
			if err != nil {
				return nil, err
			}
			svcItem := mapDbToSvc(item, user)
			resp = append(resp, svcItem)
		} else {
			svcItem := mapDbToSvc(item, dto.GetUserByIdResp{})
			resp = append(resp, svcItem)
		}
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
	userId := ctx.Value(constants.UserId)
	data, ok := userId.(int64)
	if !ok {
		logger.Error(ctx,"Error in typecasting user id")
		err = apperrors.InternalServerError
		return
	}
	reqData.UpdatedBy = data
	err = gs.gradesRepo.EditGrade(ctx, reqData)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}
	return
}

func mapDbToSvc(dbResp repository.Grade, user dto.GetUserByIdResp) (svcResp dto.Grade) {
	svcResp.Id = dbResp.Id
	svcResp.Name = dbResp.Name
	svcResp.Points = dbResp.Points
	svcResp.UpdatedBy = user.FirstName + " " + user.LastName
	return
}

func GetQuarterStartUnixTime() int64 {
	// Example function to get the Unix timestamp of the start of the quarter
	now := time.Now()
	quarterStart := time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.UTC)
	return quarterStart.Unix() * 1000 // convert to milliseconds
}
