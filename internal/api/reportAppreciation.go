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
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
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
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}
		reqData.AppreciationId = appreciationId

		resp, err := reportAppreciationSvc.ReportAppreciation(req.Context(), reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusCreated, "Appreciation reported successfully", resp)
	})
}
