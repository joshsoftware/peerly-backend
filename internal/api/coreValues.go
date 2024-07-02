package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	corevalues "github.com/joshsoftware/peerly-backend/internal/app/coreValues"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func listCoreValuesHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		coreValues, err := coreValueSvc.ListCoreValues(req.Context())
		if err != nil {

			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: coreValues, Message: "List of all the core values", Success: true})
	})
}

func getCoreValueHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		coreValue, err := coreValueSvc.GetCoreValue(req.Context(), vars["id"])
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: coreValue, Message: fmt.Sprintf("Data of core value with id: %s", vars["id"]), Success: true})
	})
}

func createCoreValueHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		const userId int64 = 1
		var coreValue dto.CreateCoreValueReq
		err := json.NewDecoder(req.Body).Decode(&coreValue)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding request data")
			err = apperrors.JSONParsingErrorReq
			apperrors.ErrorResp(rw, err)
			return
		}

		resp, err := coreValueSvc.CreateCoreValue(req.Context(), userId, coreValue)
		if err != nil {

			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusCreated, dto.SuccessResponse{Data: resp, Message: "Core value successfully created", Success: true})
	})
}

func updateCoreValueHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		var updateReq dto.UpdateQueryRequest
		err := json.NewDecoder(req.Body).Decode(&updateReq)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding request data")
			err = apperrors.JSONParsingErrorReq
			apperrors.ErrorResp(rw, err)
			return
		}

		resp, err := coreValueSvc.UpdateCoreValue(req.Context(), vars["id"], updateReq)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: resp, Message: "Core value successfully updated", Success: true})
	})
}
