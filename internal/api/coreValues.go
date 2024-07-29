package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	corevalues "github.com/joshsoftware/peerly-backend/internal/app/coreValues"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func listCoreValuesHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		coreValues, err := coreValueSvc.ListCoreValues(ctx)
		if err != nil {

			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Core values listed", coreValues)
	})
}

func getCoreValueHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		ctx := req.Context()

		coreValue, err := coreValueSvc.GetCoreValue(ctx, vars["id"])
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Core value listed", coreValue)
	})
}

func createCoreValueHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var coreValue dto.CreateCoreValueReq
		err := json.NewDecoder(req.Body).Decode(&coreValue)
		if err != nil {
			logger.Errorf("error while decoding request data, err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		ctx := req.Context()

		resp, err := coreValueSvc.CreateCoreValue(ctx, coreValue)
		if err != nil {

			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Core value created", resp)
	})
}

func updateCoreValueHandler(coreValueSvc corevalues.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		var updateReq dto.UpdateQueryRequest
		err := json.NewDecoder(req.Body).Decode(&updateReq)
		if err != nil {
			logger.Errorf("error while decoding request data, err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		ctx := req.Context()
		resp, err := coreValueSvc.UpdateCoreValue(ctx, vars["id"], updateReq)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Core value updated", resp)
	})
}
