package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	reportappreciations "github.com/joshsoftware/peerly-backend/internal/app/reportAppreciations"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

func reportAppreciationHandler(reportAppreciationSvc reportappreciations.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		if vars["id"] == "" {
			err := apperrors.InvalidId
			dto.ErrorRepsonse(rw, err)
			return
		}

		appreciationId, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			logger.Errorf(ctx, "error while parsing appreciation id from url, err: %v", err)
			err = apperrors.InternalServerError
			return

		}
		var reqData dto.ReportAppreciationReq
		err = json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.Errorf(ctx, "err while decoding request data, err: %v", err)
			dto.ErrorRepsonse(rw, err)
			return
		}
		reqData.AppreciationId = appreciationId

		resp, err := reportAppreciationSvc.ReportAppreciation(ctx, reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusCreated, "Appreciation reported successfully", resp)
	})
}

func listReportedAppreciations(reportAppreciationSvc reportappreciations.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		quarterStr := req.URL.Query().Get("quarter")
		yearStr := req.URL.Query().Get("year")
		var quarter, year int
		var err error
		if quarterStr != "" {
			quarter, err = strconv.Atoi(quarterStr)
			if err != nil {
				http.Error(rw, "Invalid quarter", http.StatusBadRequest)
				return
			}
		}
		if yearStr != "" {
			year, err = strconv.Atoi(yearStr)
			if err != nil {
				http.Error(rw, "Invalid year", http.StatusBadRequest)
				return
			}
		}

		resp, err := reportAppreciationSvc.ListReportedAppreciations(req.Context(), quarter, year)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Reported appreciations listed successfully", resp)
	})
}

func moderateAppreciation(reportAppreciationSvc reportappreciations.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		if vars["id"] == "" {
			err := apperrors.InvalidId
			dto.ErrorRepsonse(rw, err)
			return
		}
		resolutionId, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			logger.Errorf(ctx, "error while parsing id, err: %s", err.Error())
			err = apperrors.InternalServerError
			return
		}
		var reqData dto.ModerationReq
		err = json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.Errorf(ctx, "error while decoding request data, err:%s", err.Error())
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

func resolveAppreciation(reportAppreciationSvc reportappreciations.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		if vars["id"] == "" {
			err := apperrors.InvalidId
			dto.ErrorRepsonse(rw, err)
			return
		}
		resolutionId, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			logger.Errorf(ctx, "error while parsing id, err: %s", err.Error())
			err = apperrors.InternalServerError
			return
		}
		var reqData dto.ModerationReq
		err = json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.Errorf(ctx, "error while decoding request data, err:%s", err.Error())
			dto.ErrorRepsonse(rw, err)
			return
		}
		reqData.ResolutionId = resolutionId
		err = reportAppreciationSvc.ResolveAppreciation(req.Context(), reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Appreciation resolved successfully", nil)
	})
}
