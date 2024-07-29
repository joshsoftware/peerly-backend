package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	reportappreciations "github.com/joshsoftware/peerly-backend/internal/app/reportAppreciations"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func reportAppreciationHandler(reportAppreciationSvc reportappreciations.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		if vars["id"] == "" {
			err := apperrors.InvalidId
			dto.ErrorRepsonse(rw, err)
			return
		}

		appreciationId, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while parsing appreciation id from url")
			err = apperrors.InternalServerError
			return

		}
		var reqData dto.ReportAppreciationReq
		err = json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding request data")
			dto.ErrorRepsonse(rw, err)
			return
		}
		reqData.AppreciationId = appreciationId

		resp, err := reportAppreciationSvc.ReportAppreciation(req.Context(), reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusCreated, "Appreciation reported successfully", resp)
	})
}

func listReportedAppreciations(reportAppreciationSvc reportappreciations.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		appreciations, err := reportAppreciationSvc.ListReportedAppreciations(req.Context())
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation reported successfully", appreciations)
	})
}

func moderateAppriciation(reportAppreciationSvc reportappreciations.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		if vars["id"] == "" {
			err := apperrors.InvalidId
			dto.ErrorRepsonse(rw, err)
			return
		}
		resolutionId, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			logger.Errorf("error while parsing id, err: %s", err.Error())
			err = apperrors.InternalServerError
			return
		}
		var reqData dto.ModerationReq
		err = json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.Errorf("error while decoding request data, err:%s", err.Error())
			dto.ErrorRepsonse(rw, err)
			return
		}
		reqData.ResolutionId = resolutionId
		err = reportAppreciationSvc.DeleteAppreciation(req.Context(), reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation deleted successfully", nil)
	})
}
