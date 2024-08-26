package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/grades"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

func listGradesHandler(gradeSvc grades.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		resp, err := gradeSvc.ListGrades(ctx)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, 200, "grades list fetched successfully", resp)
	})
}

func editGradesHandler(gradeSvc grades.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		var reqData dto.UpdateGradeReq
		err := json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.Errorf(ctx, "error while decoding request data, err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}
		err = gradeSvc.EditGrade(ctx, vars["id"], reqData.Points)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, 200, "grade points updated successfully", nil)
	})
}
